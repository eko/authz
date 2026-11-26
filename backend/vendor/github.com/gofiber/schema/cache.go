// Copyright 2012 The Gorilla Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package schema

import (
	"errors"
	"fmt"
	"maps"
	"reflect"
	"slices"
	"strings"
	"sync"
	"sync/atomic"
	"unsafe"

	utils "github.com/gofiber/utils/v2"
	utilstrings "github.com/gofiber/utils/v2/strings"
	"github.com/gofiber/utils/v2/swar"
)

const maxParserIndex = 1000

// maxDirectKeyLen bounds the stack buffers keys are case-folded or assembled
// in for a map probe; longer keys take a path that allocates.
const maxDirectKeyLen = 64

// maxDirectPaths caps the nested entries precomputed per struct type: deep
// fan-out nesting multiplies dotted chains, and without a cap build time and
// retained memory grow exponentially. Excess keys use the generic parser.
const maxDirectPaths = 512

// maxPathReserve caps the hop and part slices the parser reserves from a
// path's separator count. That count comes from a key nothing has validated
// yet, so an unbounded reservation lets a key of mostly separators claim
// megabytes before its first segment is rejected. Appending past the
// reservation still works, and real paths are a few segments deep.
const maxPathReserve = 16

// maxPathDepth caps the fields one path names: decoding recurses per level,
// and a type that contains itself puts no other bound on the depth.
const maxPathDepth = 64

var (
	errInvalidPath   = errors.New("schema: invalid path")
	errIndexTooLarge = errors.New("schema: index exceeds parser limit")
	errPathTooDeep   = errors.New("schema: path exceeds depth limit")
)

// newCache returns a new cache.
func newCache() *cache {
	c := cache{
		tag: "schema",
	}
	return &c
}

// cache caches meta-data about a struct.
type cache struct {
	l sync.RWMutex // serializes configuration writes (tag, regconv)
	m sync.Map     // map[reflect.Type]cacheEntry
	// regconv holds the registered converters as an immutable map published
	// atomically: registerConverter replaces the whole map (copy-on-write
	// under l), so readers never touch a map that is being written.
	regconv atomic.Pointer[converterRegistry]
	tag     string
	// gen is bumped (under l) before m is cleared on configuration changes;
	// cached entries are tagged with the generation they were built under
	// and ignored on mismatch, so any call starting after a reconfiguration
	// returns observes the new configuration even if a racing build stored
	// a stale entry after the clear.
	gen atomic.Uint64
	// slots is a direct-mapped cache in front of m, indexed by the address
	// of the type's descriptor: a hit costs a few loads and a compare, where
	// m's interface-keyed hash was the largest part of what Decode paid to
	// find a type's metadata. Entries are m's own, generation included, so
	// a collision or a stale entry only falls through to m.
	slots [typeCacheSlots]atomic.Pointer[cacheEntry]
}

// typeCacheSlots is the size of cache.slots, a power of two so a slot is a
// mask rather than a division.
const typeCacheSlots = 32

// typeSlot returns the cache.slots index of t: bits of the address of its
// type descriptor, the data word of the reflect.Type interface, read
// directly since reflect.ValueOf(t).Pointer() costs several times as much.
func typeSlot(t reflect.Type) uintptr {
	return (*[2]uintptr)(unsafe.Pointer(&t))[1] >> 4 & (typeCacheSlots - 1)
}

// cacheEntry tags a structInfo with its type and the configuration
// generation it was built under; entries are stored by pointer in c.m
// (matching the encoder's encPlan pattern) and in c.slots.
type cacheEntry struct {
	info *structInfo
	typ  reflect.Type
	gen  uint64
}

type converterRegistry struct {
	m map[reflect.Type]*valueConverter
	// slots has bit typeSlot(t) set for every registered t: a clear bit skips the lookup in m.
	slots uint64
}

// slots must hold a bit for each of typeCacheSlots.
const _ = uint8(64 - typeCacheSlots)

// registerConverter registers a converter function for a custom type.
func (c *cache) registerConverter(value any, converterFunc Converter) {
	var vc *valueConverter
	if converterFunc != nil {
		vc = &valueConverter{plain: converterFunc}
	}
	c.storeConverter(value, vc)
}

// registerConverterE registers an error-returning converter function for a
// custom type.
func (c *cache) registerConverterE(value any, converterFunc ConverterE) {
	var vc *valueConverter
	if converterFunc != nil {
		vc = &valueConverter{withErr: converterFunc}
	}
	c.storeConverter(value, vc)
}

// storeConverter publishes vc as the converter of value's type, replacing any
// converter registered for it before, of either kind. A nil vc is stored as
// is, so a nil function reads back as no converter, as it always has.
func (c *cache) storeConverter(value any, vc *valueConverter) {
	c.l.Lock()
	t := reflect.TypeOf(value)
	next := &converterRegistry{m: make(map[reflect.Type]*valueConverter)}
	if prev := c.regconv.Load(); prev != nil {
		maps.Copy(next.m, prev.m)
		next.slots = prev.slots
	}
	next.m[t] = vc
	next.slots |= 1 << typeSlot(t)
	c.regconv.Store(next)
	c.reset()
	c.l.Unlock()
}

// parsePath parses a path in dotted notation verifying that it is a valid
// path to a struct field.
//
// It returns "path parts" which contain indices to fields to be used by
// reflect.Value.FieldByString(). Multiple parts are required for slices of
// structs.
func (c *cache) parsePath(p string, t reflect.Type) ([]pathPart, error) {
	if t.Kind() != reflect.Struct {
		return nil, errInvalidPath
	}
	return c.parsePathInfo(p, c.get(t))
}

// parsePathInfo is parsePath with the root struct's info already resolved,
// letting Decode look it up once per call instead of once per key. The
// parsed-path cache lives on that structInfo, keyed by the plain path
// string, which hashes much cheaper than a composite key.
func (c *cache) parsePathInfo(p string, rootInfo *structInfo) ([]pathPart, error) {
	// Fast path: probe the precomputed direct paths, whose index matches the
	// key case-insensitively in one probe, as the fields do.
	if rootInfo.directIndex.len > 0 {
		if parts, ok := rootInfo.directIndex.lookup(p); ok {
			return parts, nil
		}
		// Every field a bare alias can reach is in the direct paths
		// (buildDirectPaths adds the flat aliases before the cap can
		// truncate anything), and the probe matched the key folded, so a
		// dotless miss names no field.
		if strings.IndexByte(p, '.') < 0 {
			return nil, errInvalidPath
		}
	}

	if cached, ok := rootInfo.paths.load(p); ok {
		return cached, nil
	}

	if rootInfo.hasMapFields {
		if i := strings.IndexByte(p, '.'); i > 0 && i+1 < len(p) {
			if f := rootInfo.get(p[:i]); f != nil && f.mapRoot != nil {
				return f.mapRoot, nil
			}
		}
	}

	struc := rootInfo
	var t reflect.Type
	var field *fieldInfo
	var index64 int64
	// A path yields at most one hop and one part per segment. The hops of
	// every part share one backing array, cut into capped slices as the
	// parts are emitted.
	var parts []pathPart
	var hopBuf []pathHop
	hopStart := 0
	// Only paths spelled as declared (alias or its lower case, plain indexes) are cached.
	canonical := true
	for keyStart, depth := 0, 0; ; depth++ {
		if depth == maxPathDepth {
			return nil, errPathTooDeep
		}
		keyEnd, segment, err := nextPathSegment(p, keyStart)
		if err != nil {
			return nil, errInvalidPath
		}
		if field = struc.get(segment); field == nil {
			return nil, errInvalidPath
		}
		if segment != field.alias && segment != field.aliasLower {
			canonical = false
		}
		if parts == nil {
			segments := min(strings.Count(p, ".")+1, maxPathReserve)
			parts = make([]pathPart, 0, segments)
			hopBuf = make([]pathHop, 0, segments)
		}
		// Valid field. Append the hop; the field's index chain was resolved
		// when the structInfo was built, so the decoder walks plain indices
		// instead of repeating FieldByName lookups on every Decode call.
		hopBuf = append(hopBuf, pathHop{index: field.index, ensure: struc.anonymousPtrFields})
		if field.mapValue != nil {
			// The rest of the path is a key of this map field, kept exactly
			// as sent. Keys are chosen by the client, so this path is not
			// cached: caching it would keep every key ever sent for the life
			// of the process.
			keyStart = keyEnd + 1
			if keyStart >= len(p) {
				return nil, errInvalidPath
			}
			hops := hopBuf[hopStart:len(hopBuf):len(hopBuf)]
			return append(parts, pathPart{
				hops:      hops,
				field:     field,
				index:     -1,
				soleIndex: soleHopIndex(hops),
				keyStart:  int32(keyStart),
			}), nil
		}
		if field.takesIndex() {
			// Parse a special case: slices of structs.
			// i+1 must be the slice index.
			//
			// Now that struct can implements TextUnmarshaler interface,
			// we don't need to force the struct's fields to appear in the path.
			// So checking i+2 is not necessary anymore.
			// We can skip this part if the type is multipart.FileHeader. It is another special case too.
			keyStart = keyEnd + 1
			if keyStart >= len(p) {
				return nil, errInvalidPath
			}
			keyEnd, segment, err = nextPathSegment(p, keyStart)
			if err != nil {
				return nil, errInvalidPath
			}
			// A negative index would wrap to a large one where int is 32 bits.
			if index64, err = utils.ParseInt(segment); err != nil || index64 < 0 {
				return nil, errInvalidPath
			}
			if index64 > maxParserIndex {
				return nil, errIndexTooLarge
			}
			if c := segment[0]; c == '+' || c == '-' || (c == '0' && len(segment) > 1) {
				canonical = false
			}
			hops := hopBuf[hopStart:len(hopBuf):len(hopBuf)]
			hopStart = len(hopBuf)
			parts = append(parts, pathPart{
				hops:      hops,
				field:     field,
				index:     int(index64),
				soleIndex: soleHopIndex(hops),
			})

			// Get the next struct type, dropping ptrs.
			if field.typ.Kind() == reflect.Pointer {
				t = field.typ.Elem()
			} else {
				t = field.typ
			}
			if t.Kind() == reflect.Slice {
				t = t.Elem()
				if t.Kind() == reflect.Pointer {
					t = t.Elem()
				}
			}
		} else if field.typ.Kind() == reflect.Pointer {
			t = field.typ.Elem()
		} else {
			t = field.typ
		}

		if keyEnd == len(p) {
			break
		}
		keyStart = keyEnd + 1
		if keyStart >= len(p) {
			return nil, errInvalidPath
		}
		if t.Kind() != reflect.Struct || field.whole {
			return nil, errInvalidPath
		}
		struc = c.get(t)
	}
	// Add the remaining. A part without hops means the path terminated at a
	// slice index ("a.0"), so the decoder receives a slice element there.
	hops := hopBuf[hopStart:len(hopBuf):len(hopBuf)]
	parts = append(parts, pathPart{
		hops:      hops,
		field:     field,
		index:     -1,
		soleIndex: soleHopIndex(hops),
		elem:      len(hops) == 0,
	})

	if !canonical || len(p) > maxCachedPathLen || rootInfo.paths.full() {
		return parts, nil
	}
	// Detach the key: callers may pass strings aliasing reused request buffers.
	cost := len(p) + cap(parts)*int(unsafe.Sizeof(pathPart{})) + cap(hopBuf)*int(unsafe.Sizeof(pathHop{}))
	cached, _ := rootInfo.paths.loadOrStore(strings.Clone(p), parts, cost)
	return cached, nil
}

// maxFastPaths bounds pathCache's copy-on-write map: each publish copies the
// whole map, so past this many paths further entries spill instead. Real
// structs stay well inside it.
const maxFastPaths = 512

// Bounds on the path cache; past any, a path is parsed each time instead of kept.
const (
	maxSpillPaths     = 16384
	maxCachedPathLen  = 256
	maxPathCacheBytes = 8 << 20
)

// pathCache caches parsed paths for one struct type, keyed by the raw source
// key. Lookups run against a plain string-keyed map published atomically —
// sync.Map's interface keys cost a type hash and a trie walk, and
// parsePathInfo probes once per source key. Writes copy the map under the
// lock, so a reader only ever sees a complete one.
type pathCache struct {
	fast atomic.Pointer[map[string][]pathPart]
	// spill holds paths seen after fast was sealed, so a key space larger
	// than maxFastPaths keeps its cache.
	spill    sync.Map // map[string][]pathPart
	spillLen atomic.Int64
	// bytes is what the kept paths weigh.
	bytes  atomic.Int64
	mu     sync.Mutex
	sealed atomic.Bool
}

// full reports whether a bound is reached; racing writers may overshoot it.
func (c *pathCache) full() bool {
	return c.bytes.Load() >= maxPathCacheBytes || (c.sealed.Load() && c.spillLen.Load() >= maxSpillPaths)
}

func (c *pathCache) load(key string) ([]pathPart, bool) {
	if m := c.fast.Load(); m != nil {
		if parts, ok := (*m)[key]; ok {
			return parts, true
		}
	}
	// spill only ever receives entries once fast is sealed.
	if c.sealed.Load() {
		if v, ok := c.spill.Load(key); ok {
			return v.([]pathPart), true
		}
	}
	return nil, false
}

// loadOrStore caches parts under key, or returns what a concurrent call
// cached there first.
func (c *pathCache) loadOrStore(key string, parts []pathPart, cost int) ([]pathPart, bool) {
	if !c.sealed.Load() {
		c.mu.Lock()
		old := c.fast.Load()
		if old != nil {
			if existing, ok := (*old)[key]; ok {
				c.mu.Unlock()
				return existing, true
			}
		}
		if old == nil || len(*old) < maxFastPaths {
			var next map[string][]pathPart
			if old == nil {
				next = make(map[string][]pathPart, 1)
			} else {
				next = make(map[string][]pathPart, len(*old)+1)
				maps.Copy(next, *old)
			}
			next[key] = parts
			c.fast.Store(&next)
			c.bytes.Add(int64(cost))
			c.mu.Unlock()
			return parts, false
		}
		c.sealed.Store(true)
		c.mu.Unlock()
	}
	v, loaded := c.spill.LoadOrStore(key, parts)
	if !loaded {
		c.spillLen.Add(1)
		c.bytes.Add(int64(cost))
	}
	return v.([]pathPart), loaded
}

func (c *pathCache) clear() {
	c.mu.Lock()
	c.fast.Store(nil)
	// Nothing reaches spill before the seal, so an unsealed cache has none.
	if c.sealed.Load() {
		c.spill.Clear()
		c.spillLen.Store(0)
		c.sealed.Store(false)
	}
	c.bytes.Store(0)
	c.mu.Unlock()
}

// foldASCIILower writes the ASCII-lowercased form of s into buf and reports
// whether any byte changed. Folding into a caller-owned stack buffer lets the
// case-insensitive map probes below skip the allocation utilstrings.ToLower
// makes for a mixed-case key. The caller guarantees len(buf) >= len(s).
func foldASCIILower(buf []byte, s string) bool {
	n := len(s)
	changed := false
	i := 0
	for ; i+swar.WordLen <= n; i += swar.WordLen {
		w := swar.Load8(s, i)
		lw := swar.ToLowerWord(w)
		changed = changed || lw != w
		swar.Store8(buf, i, lw)
	}
	for ; i < n; i++ {
		ch := s[i]
		if ch >= 'A' && ch <= 'Z' {
			ch += 'a' - 'A'
			changed = true
		}
		buf[i] = ch
	}
	return changed
}

// dotBroadcast is the SWAR needle for '.'; hoisted so the word loop in
// nextPathSegment pays no per-call broadcast cost.
var dotBroadcast = swar.Broadcast('.')

func nextPathSegment(path string, start int) (int, string, error) {
	end := start
	for end+swar.WordLen <= len(path) {
		if m := swar.ZeroLanes(swar.Load8(path, end) ^ dotBroadcast); m != 0 {
			end += swar.FirstLane(m)
			if start == end {
				return 0, "", errInvalidPath
			}
			return end, path[start:end], nil
		}
		end += swar.WordLen
	}
	for end < len(path) && path[end] != '.' {
		end++
	}
	if start == end {
		return 0, "", errInvalidPath
	}
	return end, path[start:end], nil
}

// get returns a cached structInfo, creating it if necessary.
func (c *cache) get(t reflect.Type) *structInfo {
	gen := c.gen.Load()
	slot := &c.slots[typeSlot(t)]
	// Ignore entries built under an older configuration: a build racing a
	// reconfiguration may store one after the clear. Hit-validation
	// guarantees that any call starting after the reconfiguration returned
	// observes the new configuration (a call already in flight during the
	// reconfiguration may still briefly use old metadata, which is inherent
	// to concurrent reconfiguration).
	cur := slot.Load()
	if cur != nil && cur.typ == t && cur.gen == gen {
		return cur.info
	}
	if v, ok := c.m.Load(t); ok {
		if e := v.(*cacheEntry); e.gen == gen {
			claimSlot(slot, cur, e)
			return e.info
		}
	}
	info := c.create(t, "")
	if c.gen.Load() == gen {
		e := &cacheEntry{info: info, typ: t, gen: gen}
		c.m.Store(t, e)
		claimSlot(slot, cur, e)
	}
	// If the configuration changed while building, serve the result once
	// without caching it (or with a stale tag that hit-validation ignores);
	// the next call rebuilds fresh.
	return info
}

// claimSlot stores e in slot, which held cur when get read it, unless cur is
// the live entry of another type: two types sharing a slot would otherwise
// take it from each other on every call, and each store would invalidate the
// slots' cache line on every core reading it. A later type that loses the
// slot finds its metadata in c.m, as every type did before the slots.
func claimSlot(slot *atomic.Pointer[cacheEntry], cur, e *cacheEntry) {
	if cur == nil || cur.gen < e.gen {
		slot.CompareAndSwap(cur, e)
	}
}

// reset clears cached metadata and must be called with c.l held. Parsed
// path caches live on the structInfos, so dropping them drops those too.
func (c *cache) reset() {
	c.gen.Add(1)
	c.m.Clear()
	// The generation already turns the slots' entries away; dropping them
	// lets their metadata go.
	for i := range c.slots {
		c.slots[i].Store(nil)
	}
}

// aliasTag returns the configured tag name under the configuration lock, so
// metadata builds racing SetAliasTag read a consistent value.
func (c *cache) aliasTag() string {
	c.l.RLock()
	tag := c.tag
	c.l.RUnlock()
	return tag
}

// create creates a structInfo with meta-data about a struct.
func (c *cache) create(t reflect.Type, parentAlias string) *structInfo {
	// Snapshot the alias tag once so every field of this type is analyzed
	// under a consistent configuration and we don't churn the config lock
	// per field.
	tag := c.aliasTag()
	info := &structInfo{}
	var anonymousInfos []*structInfo
	var anonymousIdx [][]int
	for i := 0; i < t.NumField(); i++ {
		structField := t.Field(i)
		// Only exported anonymous pointers can be allocated; unexported ones
		// are not settable and Set would panic.
		if structField.Anonymous && structField.Type.Kind() == reflect.Pointer && structField.IsExported() {
			info.anonymousPtrFields = append(info.anonymousPtrFields, i)
		}
		if f := c.createField(structField, parentAlias, tag); f != nil {
			f.index = structField.Index
			info.fields = append(info.fields, f)
			if ft := indirectType(f.typ); ft.Kind() == reflect.Struct && f.isAnonymous && !f.whole {
				anonymousInfos = append(anonymousInfos, c.create(ft, f.canonicalAlias))
				anonymousIdx = append(anonymousIdx, structField.Index)
			}
		}
	}
	// The shadow check below reads the struct's own fields from fieldsByName.
	info.fieldsByName = make(map[string]*fieldInfo, len(info.fields))
	info.addByName(info.fields)
	own := len(info.fields)
	for i, a := range anonymousInfos {
		others := []*structInfo{info}
		others = append(others, anonymousInfos[:i]...)
		others = append(others, anonymousInfos[i+1:]...)
		for _, f := range a.fields {
			if !containsAlias(others, f.alias) {
				// Copy the promoted field so its index chain can be prefixed
				// with the embedded field's index; the original stays valid
				// for the embedded type's own structInfo.
				pf := *f
				pf.index = append(append(make([]int, 0, len(anonymousIdx[i])+len(f.index)), anonymousIdx[i]...), f.index...)
				info.fields = append(info.fields, &pf)
			}
		}
	}
	info.addByName(info.fields[own:])
	for _, f := range info.fields {
		if f.mapValue != nil {
			// Per field, not per type: a promoted copy's path is its own.
			hops := []pathHop{{index: f.index, ensure: info.anonymousPtrFields}}
			f.mapRoot = []pathPart{{hops: hops, field: f, index: -1, soleIndex: soleHopIndex(hops)}}
			info.hasMapFields = true
		}
	}
	info.requiredGroups = c.buildRequiredFields(info)
	buildRequiredIndex(info)
	info.direct = c.buildDirectPaths(info)
	info.directIndex = newFoldIndex(info.direct)
	// The setDefaults walk also allocates nil anonymous embedded pointers,
	// so it can only be skipped when neither defaults nor such pointers
	// exist anywhere in the tree.
	info.needsDefaultsWalk = c.needsDefaultsWalk(t, tag, map[reflect.Type]bool{})
	if info.needsDefaultsWalk && parentAlias == "" { // promotion-only builds never walk
		info.defaults = c.buildDefaults(info)
	}
	return info
}

// maxDefaultKeys bounds the spellings kept per field: each embedding level can
// be spelled or skipped, so they double per level.
const maxDefaultKeys = 8

// buildDefaults plans every defaulted field with its spellings through embedded
// structs, so the walk joins prefixes only at named pointers.
func (c *cache) buildDefaults(info *structInfo) defaultsPlan {
	var p defaultsPlan
	at := map[string]int{} // index chain -> step, merging the paths to one field
	add := func(f *fieldInfo, index []int, keys ...string) {
		chain := fmt.Sprint(index)
		i, ok := at[chain]
		if !ok {
			i = len(p.steps)
			at[chain] = i
			p.steps = append(p.steps, defaultStep{field: f, index: index})
		}
		for _, key := range keys {
			if !slices.Contains(p.steps[i].keys, key) {
				p.steps[i].keys = append(p.steps[i].keys, key)
			}
		}
	}
	for _, idx := range info.anonymousPtrFields {
		p.allocs = append(p.allocs, []int{idx})
	}
	for _, f := range info.fields {
		if len(f.index) > 1 && info.fieldsByName[f.aliasLower] != f {
			continue // shadowed promoted copy, reached through its embedded struct
		}
		if f.def != nil {
			add(f, f.index, f.alias)
		}
		k := f.typ.Kind()
		structPtr := k == reflect.Pointer && f.typ.Elem().Kind() == reflect.Struct
		switch {
		case f.def != nil && !f.isAnonymous: // a defaulted named struct is reported, never descended
		case k == reflect.Struct || structPtr && f.isAnonymous:
			child := c.get(indirectType(f.typ)).defaults
			if len(f.index) == 1 && f.def == nil { // a defaulted embedded struct stays zero, so its error is reported
				for _, a := range child.allocs {
					p.allocs = append(p.allocs, slices.Concat(f.index, a))
				}
			}
			for _, s := range child.steps {
				keys := make([]string, len(s.keys))
				for i, sk := range s.keys {
					keys[i] = f.alias + "." + sk
				}
				add(s.field, slices.Concat(f.index, s.index), keys...)
			}
		case structPtr:
			add(f, f.index, f.alias+".") // walked when set, under its keys
		}
	}
	for i := range p.steps {
		// Shortest first: a pointer step's first and last keys are its bare and full spellings.
		keys := p.steps[i].keys
		slices.SortStableFunc(keys, func(a, b string) int { return len(a) - len(b) })
		if len(keys) > maxDefaultKeys {
			p.steps[i].keys = append(keys[:maxDefaultKeys-1], keys[len(keys)-1])
		}
	}
	return p
}

// buildDirectPaths precomputes, keyed by lowercase path, the parsed paths for
// every key resolvable without runtime state: flat aliases plus dotted chains
// through non-pointer nested struct fields.
func (c *cache) buildDirectPaths(info *structInfo) map[string][]pathPart {
	direct := make(map[string][]pathPart, len(info.fieldsByName))
	// Flat aliases first (linear in field count) so the cap below can never
	// crowd them out; iterate fields in declaration order, honoring
	// fieldsByName's first-wins rule.
	for _, f := range info.fields {
		if !directEligible(info, f) {
			continue
		}
		hops := []pathHop{{index: f.index, ensure: info.anonymousPtrFields}}
		direct[f.aliasLower] = []pathPart{{
			hops:      hops,
			field:     f,
			index:     -1,
			soleIndex: soleHopIndex(hops),
		}}
	}
	for _, f := range info.fields {
		if len(direct) >= maxDirectPaths {
			break
		}
		if !directEligible(info, f) || f.typ.Kind() != reflect.Struct || f.whole {
			continue
		}
		// Non-pointer struct nesting cannot recurse (the type would be
		// illegal), so the child's info is always buildable here.
		hop := pathHop{index: f.index, ensure: info.anonymousPtrFields}
		for childKey, childParts := range c.get(f.typ).direct {
			if len(direct) >= maxDirectPaths {
				break
			}
			cp := childParts[0]
			hops := make([]pathHop, 0, len(cp.hops)+1)
			hops = append(hops, hop)
			hops = append(hops, cp.hops...)
			direct[f.aliasLower+"."+childKey] = []pathPart{{
				hops:      hops,
				field:     cp.field,
				index:     -1,
				soleIndex: soleHopIndex(hops),
			}}
		}
	}
	return direct
}

// directEligible reports whether f can serve as a direct-path terminal: it is
// its alias's first-wins winner, the alias has no dot, and a bare alias is a
// valid path.
func directEligible(info *structInfo, f *fieldInfo) bool {
	if info.fieldsByName[f.aliasLower] != f || strings.IndexByte(f.aliasLower, '.') >= 0 {
		return false
	}
	return f.endsPath()
}

// endsPath reports whether a path may end at f: a slice of structs needs an index after it, a map field a key.
func (f *fieldInfo) endsPath() bool {
	return !f.takesIndex() && f.mapValue == nil
}

// takesIndex reports whether a path through f must carry a slice element
// index: f is a slice of structs walked element by element, not one handed to
// a multipart binder or a TextUnmarshaler whole.
func (f *fieldInfo) takesIndex() bool {
	return f.isSliceOfStructs && !f.isMultipart &&
		(!f.unmarshalerInfo.IsValid || f.unmarshalerInfo.IsSliceElement)
}

// needsDefaultsWalk reports whether the setDefaults walk can have any effect
// on the struct tree rooted at t: it declares a default tag option, or has an
// (exported) anonymous pointer field the walk allocates, anywhere in the
// tree. visited guards against recursive types. tag is the alias tag
// snapshot for this build, so no per-field config lock is taken.
func (c *cache) needsDefaultsWalk(t reflect.Type, tag string, visited map[reflect.Type]bool) bool {
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct || visited[t] {
		return false
	}
	visited[t] = true
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		// Mirror the anonymousPtrFields filter: only exported anonymous
		// pointers are ever allocated, so only they justify the walk.
		if field.Anonymous && field.Type.Kind() == reflect.Pointer && field.IsExported() {
			return true
		}
		alias, options := fieldAlias(field, tag)
		if alias == "-" {
			continue
		}
		if options.getDefaultOptionValue() != "" {
			return true
		}
		if c.needsDefaultsWalk(field.Type, tag, visited) {
			return true
		}
	}
	return false
}

// createField creates a fieldInfo for the given field. tag is the alias tag
// snapshot for the enclosing type's build.
func (c *cache) createField(field reflect.StructField, parentAlias, tag string) *fieldInfo {
	alias, options := fieldAlias(field, tag)
	if alias == "-" {
		// Ignore this field.
		return nil
	}
	canonicalAlias := alias
	if parentAlias != "" {
		canonicalAlias = parentAlias + "." + alias
	}
	if field.Type.Kind() == reflect.Map && c.converter(field.Type) == nil &&
		!isTextUnmarshaler(reflect.Zero(field.Type)).IsValid {
		return c.createMapField(field, alias, canonicalAlias, options, tag)
	}
	// Check if the type is supported and don't cache it if not.
	// First let's get the basic type.
	isSlice, isStruct := false, false
	ft := field.Type
	m := isTextUnmarshaler(reflect.Zero(ft))
	if ft.Kind() == reflect.Pointer {
		ft = ft.Elem()
	}
	if isSlice = ft.Kind() == reflect.Slice; isSlice {
		ft = ft.Elem()
		if ft.Kind() == reflect.Pointer {
			ft = ft.Elem()
		}
	}
	if ft.Kind() == reflect.Array {
		ft = ft.Elem()
		if ft.Kind() == reflect.Pointer {
			ft = ft.Elem()
		}
	}
	if isStruct = ft.Kind() == reflect.Struct; !isStruct {
		// A map type that unmarshals itself is decoded as a whole; a slice of
		// them has no decode path, so it stays unsupported.
		if c.converter(ft) == nil && getBuiltinConverter(ft.Kind()) == nil &&
			(isSlice || ft.Kind() != reflect.Map || !m.IsValid) {
			// Type is not supported.
			return nil
		}
	}

	// Reuse the unmarshaler facts when the successive type unwrappings land
	// on the same type (the common non-pointer, non-slice case).
	derefT := indirectType(field.Type)
	derefU := m
	if derefT != field.Type {
		derefU = isTextUnmarshaler(reflect.Zero(derefT))
	}
	elemU := derefU
	if ft != derefT {
		elemU = isTextUnmarshaler(reflect.Zero(ft))
	}

	// Non-pointer builtin scalars without unmarshalers or custom converters
	// can skip decode's dispatch entirely; converter registration resets the
	// cache, so this build-time decision stays valid.
	fastKind := reflect.Invalid
	if k := field.Type.Kind(); k != reflect.Pointer && !m.IsValid &&
		getBuiltinConverter(k) != nil && c.converter(field.Type) == nil {
		fastKind = k
	}

	f := &fieldInfo{
		typ:              field.Type,
		fastKind:         fastKind,
		name:             field.Name,
		alias:            alias,
		aliasLower:       utilstrings.ToLower(alias),
		canonicalAlias:   canonicalAlias,
		unmarshalerInfo:  m,
		derefUnmarshaler: derefU,
		elemUnmarshaler:  elemU,
		isMultipart:      isMultipartField(field.Type),
		isSliceOfStructs: isSlice && isStruct,
		isAnonymous:      field.Anonymous,
		isRequired:       options.Contains("required"),
		defaultValue:     options.getDefaultOptionValue(),
	}
	f.whole = f.isMultipart || isStruct && (elemU.IsValid || c.converter(ft) != nil)
	f.resolveDefault()
	return f
}

// createMapField creates the fieldInfo for a map field, which is decoded from
// keys of the form "<alias>.<key>", one map entry per key. Only maps with
// string keys are supported, and only values the decoder fills from text:
// scalars, types with a converter or encoding.TextUnmarshaler, pointers to
// them, and slices of them. A map with no converter of its own is not
// otherwise decodable, so any other map field is ignored, as before.
func (c *cache) createMapField(field reflect.StructField, alias, canonicalAlias string, options tagOptions, tag string) *fieldInfo {
	t := field.Type
	if t.Key().Kind() != reflect.String || t.Elem().Kind() == reflect.Map {
		return nil
	}
	// The value is described as a field of the value type, so decode handles
	// it the way it handles any field of that type.
	value := c.createField(reflect.StructField{Name: field.Name, Type: t.Elem()}, "", tag)
	if value == nil || value.isMultipart || !c.decodesFromText(value) {
		return nil
	}
	f := &fieldInfo{
		typ:            t,
		fastKind:       reflect.Invalid,
		name:           field.Name,
		alias:          alias,
		aliasLower:     utilstrings.ToLower(alias),
		canonicalAlias: canonicalAlias,
		isAnonymous:    field.Anonymous,
		isRequired:     options.Contains("required"),
		defaultValue:   options.getDefaultOptionValue(),
		mapValue:       []pathPart{{field: value, index: -1, soleIndex: -1}},
		entries:        newMapEntryPool(t),
	}
	f.resolveDefault()
	return f
}

// decodesFromText reports whether a map value described by f is filled from
// text: a struct, or a slice of them, needs a converter or
// encoding.TextUnmarshaler to be, since its fields have no keys of their own
// here.
func (c *cache) decodesFromText(f *fieldInfo) bool {
	t := indirectType(f.typ)
	// A type with a converter of its own, or that unmarshals itself, such as
	// a slice type implementing encoding.TextUnmarshaler, is decoded as a
	// whole.
	if c.converter(t) != nil || f.derefUnmarshaler.IsValid && !f.derefUnmarshaler.IsSliceElement {
		return true
	}
	if f.isSliceOfStructs {
		return f.elemUnmarshaler.IsValid || c.converter(indirectType(t.Elem())) != nil
	}
	if t.Kind() == reflect.Struct {
		return f.derefUnmarshaler.IsValid || c.converter(t) != nil
	}
	return true
}

// converter returns the converter registered for a type, or nil.
func (c *cache) converter(t reflect.Type) *valueConverter {
	reg := c.regconv.Load()
	if reg == nil || reg.slots&(1<<typeSlot(t)) == 0 {
		return nil
	}
	return reg.m[t]
}

// ----------------------------------------------------------------------------

type structInfo struct {
	fields             []*fieldInfo
	fieldsByName       map[string]*fieldInfo
	anonymousPtrFields []int
	// requiredGroups lists the required keys in a stable order, so each has
	// an index the satisfied-group bitset can address; requiredPrefixes maps
	// a nested-key prefix to the groups a key under it can satisfy.
	requiredGroups       []requiredGroup
	requiredPrefixes     map[string][]requiredPrefix
	maxRequiredPrefixLen int
	// The required paths and prefixes lowercased, for keys in another case.
	foldedRequiredPaths    map[string][]requiredPrefix
	foldedRequiredPrefixes map[string][]requiredPrefix
	// direct maps lowercase statically-resolvable keys to their precomputed
	// parsed paths; built once and immutable, see buildDirectPaths.
	direct map[string][]pathPart
	// directIndex is direct, indexed for parsePathInfo's case-insensitive
	// lookups
	directIndex foldIndex
	// paths caches parsed paths rooted at this struct type; keys are cloned
	// so they never alias reused request buffers.
	paths        pathCache
	hasMapFields bool
	// defaults is the setDefaults walk of this struct tree, see buildDefaults.
	defaults defaultsPlan
	// needsDefaultsWalk reports whether the setDefaults walk can have any
	// effect on this struct tree: it is set when a default tag option or an
	// anonymous embedded pointer field (which the walk allocates) exists
	// anywhere in the tree, letting the decoder skip the walk otherwise.
	needsDefaultsWalk bool
	// requiredMixedCase is set if a required path or prefix is not lowercase.
	requiredMixedCase bool
}

// addByName indexes fields by lowercase alias; the first field of an alias wins.
func (i *structInfo) addByName(fields []*fieldInfo) {
	for _, field := range fields {
		if _, exists := i.fieldsByName[field.aliasLower]; !exists {
			i.fieldsByName[field.aliasLower] = field
		}
	}
}

func (i *structInfo) get(alias string) *fieldInfo {
	// fieldsByName is keyed by lowercase alias, so the raw probe settles an
	// already-lowercase one without folding. A mixed-case alias is folded
	// into a stack buffer and probed from there: the compiler elides the
	// string conversion in a map index, so that allocates nothing either.
	if field, ok := i.fieldsByName[alias]; ok {
		return field
	}
	if n := len(alias); n <= maxDirectKeyLen {
		var buf [maxDirectKeyLen]byte
		if !foldASCIILower(buf[:], alias) {
			return nil // already lowercase: the probe above was conclusive
		}
		return i.fieldsByName[string(buf[:n])]
	}
	return i.fieldsByName[utilstrings.ToLower(alias)]
}

func (c *cache) buildRequiredFields(info *structInfo) []requiredGroup {
	var groups []requiredGroup
	var byKey map[string]int
	add := func(key string, f fieldWithPrefix) {
		if i, ok := byKey[key]; ok {
			groups[i].fields = append(groups[i].fields, f)
			return
		}
		if byKey == nil {
			byKey = make(map[string]int)
		}
		byKey[key] = len(groups)
		groups = append(groups, requiredGroup{key: key, fields: []fieldWithPrefix{f}})
	}
	for _, field := range info.fields {
		// Fields promoted through an embedded pointer are required as through a value.
		if ft := indirectType(field.typ); ft.Kind() == reflect.Struct && (ft == field.typ || field.isAnonymous) && !field.whole {
			nested := c.get(ft)
			for _, prefix := range field.paths("") {
				nestedPrefix := prefix + "."
				for _, group := range nested.requiredGroups {
					requiredKey := field.canonicalAlias + "." + group.key
					for _, nestedField := range group.fields {
						add(requiredKey, newFieldWithPrefix(nestedField.fieldInfo, nestedPrefix+nestedField.prefix))
					}
				}
			}
		}
		if field.isRequired {
			add(field.canonicalAlias, newFieldWithPrefix(field, ""))
		}
	}
	return groups
}

// buildRequiredIndex indexes info's required groups by the nested-key prefix
// that can satisfy them ("d." for required key "d"), so a source key can be
// resolved against all of them at once, and by their paths and prefixes
// lowercased, for keys in another case.
func buildRequiredIndex(info *structInfo) {
	groups := info.requiredGroups
	if len(groups) == 0 {
		return
	}
	info.requiredPrefixes = make(map[string][]requiredPrefix)
	info.foldedRequiredPaths = make(map[string][]requiredPrefix)
	info.foldedRequiredPrefixes = make(map[string][]requiredPrefix)
	for gi := range groups {
		for _, f := range groups[gi].fields {
			for _, path := range f.searchPaths {
				lower := utilstrings.ToLower(path)
				info.requiredMixedCase = info.requiredMixedCase || lower != path
				addRequiredOwner(info.foldedRequiredPaths, lower, gi, f.typ)
			}
			typ := f.typ
			if f.mapValue != nil {
				// A key under a map field gives one entry, so it is empty
				// when a value of the map's value type would be.
				typ = typ.Elem()
			}
			for _, dot := range f.searchPathDots {
				addRequiredOwner(info.requiredPrefixes, dot, gi, typ)
				lower := utilstrings.ToLower(dot)
				info.requiredMixedCase = info.requiredMixedCase || lower != dot
				addRequiredOwner(info.foldedRequiredPrefixes, lower, gi, typ)
				info.maxRequiredPrefixLen = max(info.maxRequiredPrefixLen, len(dot))
			}
		}
	}
}

// addRequiredOwner adds group gi to index[key], once per field type: several
// fields of a group can reach the same key.
func addRequiredOwner(index map[string][]requiredPrefix, key string, gi int, typ reflect.Type) {
	owners := index[key]
	for _, o := range owners {
		if o.group == gi && o.typ == typ {
			return
		}
	}
	index[key] = append(owners, requiredPrefix{group: gi, typ: typ})
}

func containsAlias(infos []*structInfo, alias string) bool {
	aliasKey := utilstrings.ToLower(alias)
	for _, info := range infos {
		if _, ok := info.fieldsByName[aliasKey]; ok {
			return true
		}
	}
	return false
}

type fieldInfo struct {
	typ reflect.Type
	// fastKind is the field's builtin scalar kind when decode can set it
	// directly (no pointer, unmarshaler, or custom converter); else Invalid.
	fastKind reflect.Kind
	// index is the field index chain relative to the struct type whose
	// structInfo holds this fieldInfo; promoted fields carry the full chain
	// through the embedded structs (a copy is made per promotion level).
	index []int
	// name is the field name in the struct.
	name  string
	alias string
	// aliasLower is the pre-computed lowercase alias for fast lookups.
	aliasLower string
	// canonicalAlias is almost the same as the alias, but is prefixed with
	// an embedded struct field alias in dotted notation if this field is
	// promoted from the struct.
	// For instance, if the alias is "N" and this field is an embedded field
	// in a struct "X", canonicalAlias will be "X.N".
	canonicalAlias string
	// unmarshalerInfo contains information regarding the
	// encoding.TextUnmarshaler implementation of the field type.
	unmarshalerInfo unmarshaler
	// derefUnmarshaler caches the encoding.TextUnmarshaler information for
	// the field type after one pointer dereference, which is what the
	// decoder sees after walking to the field. Only the type-level flags are
	// meaningful; the decoder binds instances to live values itself.
	derefUnmarshaler unmarshaler
	// elemUnmarshaler is like derefUnmarshaler but for the fully unwrapped
	// slice element type; the decoder uses it when a path terminates at a
	// slice index (e.g. "a.0") and the value at hand is an element rather
	// than the slice field itself.
	elemUnmarshaler unmarshaler
	// isMultipart indicates whether the field type is one of the supported
	// multipart file header shapes, precomputed so the decoder can skip the
	// type comparisons on every other field.
	isMultipart bool
	// whole marks a file field, or a struct (or slice of them) that a converter or
	// encoding.TextUnmarshaler fills: no key continues into the struct's own fields.
	whole bool
	// isSliceOfStructs indicates if the field type is a slice of structs.
	isSliceOfStructs bool
	// isAnonymous indicates whether the field is embedded in the struct.
	isAnonymous  bool
	isRequired   bool
	defaultValue string
	// def is the default option resolved at build time, nil for the fields
	// that have none; see resolveDefault.
	def *fieldDefault
	// mapValue is set for a map field decoded entry by entry: the path part
	// decode fills a map value with, describing the map's value type as a
	// field of its own. Nil for every other field.
	mapValue []pathPart
	// mapRoot is the path of an entry named by the bare alias, the same for every key.
	mapRoot []pathPart
	entries *sync.Pool
}

type defaultsPlan struct {
	allocs [][]int // nil embedded pointers to allocate, parents first
	steps  []defaultStep
}

// defaultStep is a defaulted field, or a named struct pointer walked when set,
// whose keys then end in the separator.
type defaultStep struct {
	field *fieldInfo
	index []int
	keys  []string
}

// fieldDefault is what setDefaults assigns for a field's default option:
// val for a scalar field (converted to the field type) or a pointer field
// (converted to the pointee type), and slice, a template a slice field's
// default is copied from. An element that failed to convert leaves the
// template short and sets err, which is reported each time the default
// applies. A kind that takes no default leaves all three unset.
type fieldDefault struct {
	val   reflect.Value
	slice reflect.Value
	err   error
}

func (f *fieldInfo) paths(prefix string) []string {
	if f.alias == f.canonicalAlias {
		return []string{prefix + f.alias}
	}
	return []string{prefix + f.alias, prefix + f.canonicalAlias}
}

type pathPart struct {
	field *fieldInfo
	hops  []pathHop // path to the field: walks structs using field indices.
	index int       // struct index in slices of structs.
	// soleIndex is the single struct field index hops walks, or -1; see
	// soleHopIndex.
	soleIndex int
	// elem marks a terminal part whose path ended at a slice index ("a.0"):
	// the decoder's value is then an element of the slice field rather than
	// the field itself.
	elem bool
	// keyStart is where, in the path, the key of the map entry that a
	// terminal part through a map field names begins. Such a part is never
	// cached, so the offset always refers to the path being decoded; an int32
	// sits in elem's padding and leaves every other part the size it was.
	// Zero marks a shared mapRoot part, whose key follows the first '.'.
	keyStart int32
}

// soleHopIndex returns the single struct field index hops walks, or -1 when
// the walk needs decode's general loop: several hops, a promoted field's
// index chain, or embedded pointers to allocate on the way.
func soleHopIndex(hops []pathHop) int {
	if len(hops) != 1 || len(hops[0].ensure) != 0 || len(hops[0].index) != 1 {
		return -1
	}
	return hops[0].index[0]
}

// pathHop describes one named-field lookup along a path. index is the field
// index chain relative to the struct at this level (more than one element
// when the field is promoted from embedded structs), and ensure lists the
// anonymous pointer fields of that struct which must be allocated before the
// walk so promoted fields stay reachable.
type pathHop struct {
	index  []int
	ensure []int
}

// ----------------------------------------------------------------------------

func indirectType(typ reflect.Type) reflect.Type {
	if typ.Kind() == reflect.Pointer {
		return typ.Elem()
	}
	return typ
}

// fieldAlias parses a field tag to get a field alias.
func fieldAlias(field reflect.StructField, tagName string) (alias string, options tagOptions) {
	if tag := field.Tag.Get(tagName); tag != "" {
		alias, options = parseTag(tag)
	}
	if alias == "" {
		alias = field.Name
	}
	return alias, options
}

// tagOptions is the string following a comma in a struct field's tag, or
// the empty string. It does not include the leading comma. Keeping the raw
// comma-separated string avoids the []string allocation of strings.Split on
// hot paths (the encoder parses tags on every Encode call).
type tagOptions string

// parseTag splits a struct field's url tag into its name and comma-separated
// options.
func parseTag(tag string) (string, tagOptions) {
	if before, after, ok := strings.Cut(tag, ","); ok {
		return before, tagOptions(after)
	}
	return tag, ""
}

// Contains checks whether the tagOptions contains the specified option.
func (o tagOptions) Contains(option string) bool {
	if o == "" {
		return false
	}
	for s := range strings.SplitSeq(string(o), ",") {
		if s == option {
			return true
		}
	}
	return false
}

func (o tagOptions) getDefaultOptionValue() string {
	if o == "" {
		return ""
	}
	for s := range strings.SplitSeq(string(o), ",") {
		if value, ok := strings.CutPrefix(s, "default:"); ok {
			return value
		}
	}
	return ""
}

// foldIndex is a read-only open-addressing index from lowercase keys to their
// parsed paths that is probed ASCII case-insensitively. A key's slot comes from
// a hash of its case-folded bytes and is confirmed with utils.EqualFold, so a
// mixed-case key, a header name say, costs one probe where a map needed a
// probe with the key as given and another with it folded, and a key naming no
// field is turned away just as fast.
type foldIndex struct {
	// slots is a power of two long, at most half full, so every probe
	// reaches a free slot
	slots []foldSlot
	mask  uint64
	len   int
}

type foldSlot struct {
	key   string // lowercase key; "" marks a free slot, since no path is empty
	parts []pathPart
}

// newFoldIndex indexes paths, whose keys are lowercase and non-empty.
func newFoldIndex(paths map[string][]pathPart) foldIndex {
	if len(paths) == 0 {
		return foldIndex{}
	}
	size := 8
	for size < 2*len(paths) {
		size *= 2
	}
	x := foldIndex{slots: make([]foldSlot, size), mask: uint64(size - 1), len: len(paths)}
	for key, parts := range paths {
		i := utils.HashFold(key) & x.mask
		for x.slots[i].key != "" {
			i = (i + 1) & x.mask
		}
		x.slots[i] = foldSlot{key: key, parts: parts}
	}
	return x
}

// lookup returns the parts of the key that equals k ASCII case-insensitively.
func (x *foldIndex) lookup(k string) ([]pathPart, bool) {
	for i := utils.HashFold(k) & x.mask; ; i = (i + 1) & x.mask {
		s := &x.slots[i]
		if s.key == "" {
			return nil, false
		}
		if len(s.key) == len(k) && utils.EqualFold(s.key, k) {
			return s.parts, true
		}
	}
}
