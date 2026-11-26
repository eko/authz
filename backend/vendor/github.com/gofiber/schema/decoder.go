// Copyright 2012 The Gorilla Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package schema

import (
	"encoding"
	"errors"
	"fmt"
	"maps"
	"math/bits"
	"mime/multipart"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"unsafe"

	utils "github.com/gofiber/utils/v2"
	utilstrings "github.com/gofiber/utils/v2/strings"
)

const (
	defaultMaxSize = 16000
	// maxAllocBytes bounds the slices and pointer targets one call allocates besides the destination's own
	// slice fields: keys reach new ones per parent element or pointer path, so no per-field limit bounds them.
	maxAllocBytes = 8 << 20
)

// errNotPointerToStruct is returned by Decode for invalid destinations;
// hoisted so the check does not allocate on every call.
var errNotPointerToStruct = errors.New("schema: interface must be a pointer to struct")

var errAllocLimit = errors.New("schema: keys allocate past the per-call limit")

// fileKeyValues stands in for the value of a multipart file's key in the
// decode view. Nothing writes through a source map's values and the map it
// goes into is this package's own copy, so one slice serves every file key.
var fileKeyValues = []string{""}

// NewDecoder returns a new Decoder.
func NewDecoder() *Decoder {
	return &Decoder{cache: newCache(), maxSize: defaultMaxSize}
}

// Decoder decodes values from a map[string][]string to a struct.
type Decoder struct {
	cache             *cache
	zeroEmpty         bool
	ignoreUnknownKeys bool
	maxSize           int
}

// SetAliasTag changes the tag used to locate custom field aliases.
// The default tag is "schema".
func (d *Decoder) SetAliasTag(tag string) {
	d.cache.l.Lock()
	d.cache.tag = tag
	d.cache.reset()
	d.cache.l.Unlock()
}

// ZeroEmpty controls the behaviour when the decoder encounters empty values
// in a map.
// If z is true and a key in the map has the empty string as a value
// then the corresponding struct field is set to the zero value, unless no
// value could set it, as for a struct without a converter or TextUnmarshaler.
// If z is false then empty strings are ignored.
//
// The default value is false, that is empty values do not change
// the value of the struct field.
func (d *Decoder) ZeroEmpty(z bool) {
	d.zeroEmpty = z
}

// IgnoreUnknownKeys controls the behaviour when the decoder encounters unknown
// keys in the map.
// If i is true and an unknown field is encountered, it is ignored. This is
// similar to how unknown keys are handled by encoding/json.
// If i is false then Decode will return an error. Note that any valid keys
// will still be decoded in to the target struct.
//
// To preserve backwards compatibility, the default value is false.
func (d *Decoder) IgnoreUnknownKeys(i bool) {
	d.ignoreUnknownKeys = i
}

// MaxSize limits the size of slices for URL nested arrays or object arrays.
// Choose MaxSize carefully; large values may create many zero-value slice elements.
// Example: "items.100000=apple" would create a slice with 100,000 empty strings.
func (d *Decoder) MaxSize(size int) {
	d.maxSize = size
}

// RegisterConverter registers a converter function for a custom type,
// replacing any converter registered for the type before, including one
// registered with RegisterConverterE.
func (d *Decoder) RegisterConverter(value any, converterFunc Converter) {
	d.cache.registerConverter(value, converterFunc)
}

// RegisterConverterE registers a converter function for a custom type that
// reports why a value does not convert: a non-nil error fails the conversion
// and is kept as the Err of the returned ConversionError. It replaces any
// converter registered for the type before, including one registered with
// RegisterConverter.
func (d *Decoder) RegisterConverterE(value any, converterFunc ConverterE) {
	d.cache.registerConverterE(value, converterFunc)
}

// Decode decodes a map[string][]string to a struct.
//
// The first parameter must be a pointer to a struct.
//
// The second parameter is a map, typically url.Values from an HTTP request.
// Keys are "paths" in dotted notation to the struct fields and nested structs.
//
// See the package documentation for a full explanation of the mechanics.
func (d *Decoder) Decode(dst any, src map[string][]string, files ...map[string][]*multipart.FileHeader) error {
	var multipartFiles map[string][]*multipart.FileHeader

	if len(files) > 0 {
		multipartFiles = files[0]
	}

	// Add files as empty string values to the decode view so path parsing
	// works uniformly. Work on a copy: the caller's src map must not be
	// mutated (and a caller-provided value under a file's key must not be
	// overwritten in it).
	if len(multipartFiles) > 0 {
		merged := make(map[string][]string, len(src)+len(multipartFiles))
		maps.Copy(merged, src)
		for path := range multipartFiles {
			merged[path] = fileKeyValues
		}
		src = merged
	}

	s := source{m: src}
	return d.decodeSource(dst, &s, multipartFiles)
}

// DecodeValues decodes key/value pairs to a struct as Decode decodes the map
// they make: values[i] is a value of the key keys[i], and a key given more
// than once has all of its values, in the order given. It saves a caller that
// has the pairs one at a time, as a parsed query string does, building that
// map, and saves the decode iterating it.
//
// Where several keys name the same field, as keys that differ only in case
// do, the keys decode in the order of their first pairs, so the key given
// first the latest wins; Decode lets whichever its map yields last win.
//
// The first parameter must be a pointer to a struct, and keys and values must
// be the same length. The slices are only read, and only during the call.
func (d *Decoder) DecodeValues(dst any, keys, values []string) error {
	if len(keys) != len(values) {
		return errValuesLength
	}
	// The index is set up field by field where it lies. Go puts a literal for
	// a variable whose address is taken, or a call's result, together
	// elsewhere and copies it over, and the copy reads back in wide loads the
	// words just stored in narrow ones, which stalls each load.
	var p pairIndex
	p.keys, p.values = keys, values
	if n := len(keys); n <= maxScanPairs {
		var next, last [maxScanPairs]int32
		p.next, p.last = next[:n], last[:n]
	} else if n <= maxInlinePairs {
		var table [2 * maxInlinePairs]int32
		var next, last [maxInlinePairs]int32
		p.table, p.next, p.last = table[:indexSize(n)], next[:n], last[:n]
	} else {
		p.table, p.next, p.last = make([]int32, indexSize(n)), make([]int32, n), make([]int32, n)
	}
	p.index()
	s := source{pairs: &p}
	return d.decodeSource(dst, &s, nil)
}

// decodeSource decodes src into dst, the pointer to a struct Decode and
// DecodeValues were given, with the files Decode was.
func (d *Decoder) decodeSource(dst any, src *source, multipartFiles map[string][]*multipart.FileHeader) (err error) {
	v := reflect.ValueOf(dst)
	if v.Kind() != reflect.Pointer || v.Elem().Kind() != reflect.Struct {
		return errNotPointerToStruct
	}

	// Catch panics from the decoder and return them as an error.
	// This is needed because the decoder calls reflect and reflect panics.
	// Installed before any other work so nothing can crash the caller.
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(error); ok {
				err = e
			} else {
				err = fmt.Errorf("schema: panic while decoding: %v", r)
			}
		}
	}()

	v = v.Elem()
	t := v.Type()
	rootInfo := d.cache.get(t)
	// Required-key bookkeeping rides along with the loop below so src is
	// walked once: the direct lookups settle almost every group here, and
	// the rest are answered by the nested keys the loop visits anyway.
	var satisfied, decided []uint64
	pending := 0
	if len(rootInfo.requiredGroups) > 0 {
		// Declared here so a struct with no required keys never pays to
		// zero it.
		var requiredBits [requiredBitWords]uint64
		satisfied, pending = markProvidedDirectly(rootInfo.requiredGroups, src, requiredBits[:])
		if d.zeroEmpty {
			// The groups an empty value has zeroed, which the decode order
			// settles from then on (see markAssigned).
			var decidedBits [requiredBitWords]uint64
			if len(satisfied) <= len(decidedBits) {
				decided = decidedBits[:len(satisfied)]
			} else {
				decided = make([]uint64, len(satisfied))
			}
		}
	}
	var grow growTracker
	// The state decodeKey carries from key to key is kept in variables of
	// its own rather than a struct: the compiler would let the stack values
	// above escape along with the errors a struct holding both returns.
	var errs MultiError
	if p := src.pairs; p != nil {
		// A key whose pairs are not adjacent gathers its values here.
		var scratch [8]string
		for i, path := range p.keys {
			if p.first(i) {
				errs, pending = d.decodeKey(v, rootInfo, multipartFiles, &grow, satisfied, decided, path, p.group(i, scratch[:0]), errs, pending)
			}
		}
	} else {
		for path, values := range src.m {
			errs, pending = d.decodeKey(v, rootInfo, multipartFiles, &grow, satisfied, decided, path, values, errs, pending)
		}
	}
	if pending > 0 {
		if anyBit(decided) {
			pending = markProvidedFoldedUndecided(d.cache, rootInfo, src, satisfied, decided, pending)
		} else {
			pending = markProvidedFolded(d.cache, rootInfo, src, satisfied, pending)
		}
	}
	if rootInfo.needsDefaultsWalk {
		errs = mergeErrors(errs, d.setDefaults(t, v, src, "", "", 0))
	}
	errs = mergeErrors(errs, missingRequired(rootInfo.requiredGroups, satisfied, pending))
	if len(errs) > 0 {
		return errs
	}
	return nil
}

// decodeKey decodes the values of one source key into v, the struct whose
// metadata is rootInfo, and returns errs with any error it met added, and
// pending less the required groups the key settles in satisfied.
func (d *Decoder) decodeKey(
	v reflect.Value, rootInfo *structInfo, files map[string][]*multipart.FileHeader, grow *growTracker,
	satisfied, decided []uint64, path string, values []string, errs MultiError, pending int,
) (MultiError, int) {
	if parts, err := d.cache.parsePathInfo(path, rootInfo); err == nil {
		// Only an empty value can zero a required field, and only a key
		// after one did can provide it again.
		if decided != nil && (len(values) == 0 || values[len(values)-1] == "" || anyBit(decided)) {
			pending = markAssigned(rootInfo, path, values, satisfied, decided, pending)
		}
		if pending > 0 {
			if i := strings.IndexByte(path, '.'); i >= 0 && len(values) > 0 {
				pending = markProvidedByNestedKey(d.cache, true, rootInfo, path, i+1, values, satisfied, pending)
			}
		}
		var filesSlice []*multipart.FileHeader
		if files != nil {
			filesSlice = files[path]
		}
		if err = d.decode(v, path, parts, values, filesSlice, grow, false); err != nil {
			errs = appendError(errs, path, err)
		}
		// A key that names no field is by far the most common failure, and
		// parsePathInfo returns its sentinels unwrapped; the others are limits.
	} else if err != errInvalidPath { //nolint:errorlint // see above
		errs = appendError(errs, path, err)
	} else if !d.ignoreUnknownKeys {
		errs = appendError(errs, path, UnknownKeyError{Key: path})
	}
	return errs, pending
}

// setDefaults sets the default values when the `default` tag is specified,
// default is supported on basic/primitive types and their pointers,
// nested structs can also have default tags
var (
	errUnsupportedDefault = errors.New("default option is supported only on: bool, float variants, string, unit variants types or their corresponding pointers or slices")
	errRequiredDefault    = errors.New("required fields cannot have a default value")
)

// resolveDefault converts the field's default option once, when the struct
// metadata is built, so setDefaults assigns a resolved value on every call
// instead of converting the string each time. Slice and pointer defaults are
// still instantiated per call, so no two decoded values share one.
func (f *fieldInfo) resolveDefault() {
	if f.defaultValue == "" {
		return
	}
	def := &fieldDefault{}
	f.def = def
	switch f.typ.Kind() {
	case reflect.Slice:
		elemT := f.typ.Elem()
		conv := getBuiltinConverter(elemT.Kind())
		if conv == nil {
			return
		}
		tmpl := reflect.MakeSlice(f.typ, 0, strings.Count(f.defaultValue, "|")+1)
		for val := range strings.SplitSeq(f.defaultValue, "|") {
			v := conv(val)
			if !v.IsValid() {
				def.err = fmt.Errorf("failed setting default: %s is not compatible with field %s type", val, f.name)
				break
			}
			// Builtin converters return the underlying kind; convert to the
			// (possibly named) element type, else Append panics for []MyInt.
			tmpl = reflect.Append(tmpl, v.Convert(elemT))
		}
		def.slice = tmpl
	case reflect.Pointer:
		t1 := f.typ.Elem()
		if conv := getBuiltinConverter(t1.Kind()); conv != nil {
			if v := conv(f.defaultValue); v.IsValid() {
				def.val = v.Convert(t1)
			}
		}
	default:
		if conv := getBuiltinConverter(f.typ.Kind()); conv != nil {
			if v := conv(f.defaultValue); v.IsValid() {
				def.val = v.Convert(f.typ)
			}
		}
	}
}

// setDefaults runs the planned walk of t, reached under prefix, the shortest
// path, and under long, the fully spelled one, depth levels below the root.
func (d *Decoder) setDefaults(t reflect.Type, v reflect.Value, src *source, prefix, long string, depth int) MultiError {
	struc := d.cache.get(t)
	// Skip the walk entirely when it can have no effect (no default tags and
	// no anonymous embedded pointers to allocate anywhere in the tree) — the
	// overwhelmingly common case.
	if !struc.needsDefaultsWalk {
		return nil
	}

	var errs MultiError

	// Allocate nil anonymous embedded pointer fields so their promoted
	// fields stay reachable.
	for _, chain := range struc.defaults.allocs {
		if field, _ := walkIndexChain(v, chain, nil); field.IsValid() && field.IsNil() {
			field.Set(reflect.New(field.Type().Elem()))
		}
	}

	for _, s := range struc.defaults.steps {
		f := s.field
		vCurrent, _ := walkIndexChain(v, s.index, nil) // nothing to refuse without a tracker
		if !vCurrent.IsValid() {
			// Unreachable behind an unsettable nil embedded pointer.
			continue
		}

		def := f.def
		if def == nil {
			if !vCurrent.IsNil() {
				// Mixed spellings would multiply per level, and recursive types nest without limit.
				np, nl := prefix, long
				// No key names more than maxPathDepth fields, so a deeper prefix could never match; a
				// destination filled by another decoder can be far deeper, and each prefix would stay live.
				if depth < maxPathDepth {
					np = prefix + s.keys[0]
					nl = np
					if long != prefix || len(s.keys) > 1 {
						nl = long + s.keys[len(s.keys)-1]
					}
				}
				errs = mergeErrors(errs, d.setDefaults(f.typ.Elem(), vCurrent.Elem(), src, np, nl, depth+1))
			}
			continue
		}
		if f.isRequired {
			errs = appendError(errs, "default-"+f.name, errRequiredDefault)
			continue
		}
		if !vCurrent.IsZero() || depth < maxPathDepth && keysProvided(src, prefix, long, s.keys) {
			continue
		}
		kind := f.typ.Kind()
		// The default itself was resolved when the metadata was built; see
		// resolveDefault. Only the per-call parts remain: the errors a kind
		// that takes no default reports, and a fresh pointer or slice so no
		// two decoded values share one.
		switch kind {
		case reflect.Struct:
			errs = appendError(errs, "default-"+f.name, errUnsupportedDefault)
		case reflect.Slice:
			if !def.slice.IsValid() {
				errs = appendError(errs, "default-"+f.name, errUnsupportedDefault)
				continue
			}
			if def.err != nil {
				errs = appendError(errs, "default-"+f.name, def.err)
			}
			fresh := reflect.MakeSlice(f.typ, def.slice.Len(), def.slice.Cap())
			reflect.Copy(fresh, def.slice)
			vCurrent.Set(fresh)
		case reflect.Pointer:
			t1 := f.typ.Elem()
			if t1.Kind() == reflect.Struct || t1.Kind() == reflect.Slice {
				errs = appendError(errs, "default-"+f.name, errUnsupportedDefault)
			}
			if def.val.IsValid() {
				// *elem is assignable to the field even when the field's
				// type is itself a named pointer type (type MyIntPtr *MyInt),
				// where converting a *int directly would panic.
				p := reflect.New(t1)
				p.Elem().Set(def.val)
				vCurrent.Set(p)
			}
		default:
			if getBuiltinConverter(kind) == nil {
				errs = appendError(errs, "default-"+f.name, errUnsupportedDefault)
			} else if def.val.IsValid() {
				vCurrent.Set(def.val)
			}
		}
	}

	return errs
}

// growCap is the capacity to give a slice being grown to n elements: room for
// another doubling when the slice is tracked, so later indices mostly land
// inside it, and exactly n when it is not, since the extra would go unused.
func (d *Decoder) growCap(n int, tracked bool) int {
	if !tracked {
		return n
	}
	c := 2 * n
	if c < n || c > d.maxSize+1 { // overflow, or past what maxSize admits
		c = d.maxSize + 1
	}
	if c < n {
		c = n
	}
	return c
}

func keysProvided(src *source, prefix, long string, keys []string) bool {
	for _, key := range keys {
		if keyProvided(src, prefix, key) || long != prefix && keyProvided(src, long, key) {
			return true
		}
	}
	return false
}

// keyProvided reports whether prefix+name is a key of src, assembling the key
// in a stack buffer when it fits so the probe allocates nothing.
func keyProvided(src *source, prefix, name string) bool {
	if prefix == "" {
		_, ok := src.lookup(name)
		return ok
	}
	if n := len(prefix) + len(name); n <= maxDirectKeyLen {
		var buf [maxDirectKeyLen]byte
		copy(buf[copy(buf[:], prefix):], name)
		return src.hasBytes(buf[:n])
	}
	_, ok := src.lookup(prefix + name)
	return ok
}

// The set of required fields (including those of nested structs) is
// precomputed once per struct type in structInfo.requiredGroups, so what
// follows only performs the per-request emptiness checks against src.
//
// A group is satisfied by a value under one of its own paths, or by any
// nested key below one of them ("d.e" satisfies required "d"). Direct paths
// are looked up first, since they settle almost every group; what is left is
// answered by walking the dotted prefixes of each source key, then by the keys
// in another case. Decode runs those steps around its own loop over src;
// checkRequired composes them for a standalone check.

// requiredBitWords sizes the inline satisfied-group bitset: 256 required
// keys, past which it is allocated.
const requiredBitWords = 4

// markProvidedDirectly marks every required group src answers through one of
// its own paths, returning the bitset — backed by inline when the groups fit
// — and how many are still pending.
func markProvidedDirectly(groups []requiredGroup, src *source, inline []uint64) ([]uint64, int) {
	if len(groups) == 0 {
		return nil, 0
	}
	var satisfied []uint64
	if w := (len(groups) + 63) >> 6; w <= len(inline) {
		satisfied = inline[:w]
	} else {
		satisfied = make([]uint64, w)
	}
	pending := len(groups)
	for gi := range groups {
		if directlyProvided(groups[gi].fields, src) {
			satisfied[gi>>6] |= 1 << uint(gi&63)
			pending--
		}
	}
	return satisfied, pending
}

// markProvidedByNestedKey walks the dotted prefixes of key from off, marking
// every required group they name that val is non-empty for, and returns the
// new pending count. Callers test for a dotted, non-empty key themselves, so
// the keys that can answer nothing do not pay for a call. Only a key the
// decoder can parse counts; unless parsed is set, c parses it on the first
// prefix that names a group.
func markProvidedByNestedKey(c *cache, parsed bool, info *structInfo, key string, off int, val []string, satisfied []uint64, pending int) int {
	// Each lookup hashes the whole prefix, so stopping at the longest one keeps a long key linear.
	for off <= info.maxRequiredPrefixLen {
		if owners := info.requiredPrefixes[key[:off]]; len(owners) > 0 {
			if !parsed {
				if _, err := c.parsePathInfo(key, info); err != nil {
					return pending
				}
				parsed = true
			}
			pending = markProvided(owners, val, satisfied, pending)
		}
		i := strings.IndexByte(key[off:], '.')
		if i < 0 {
			return pending
		}
		off += i + 1
	}
	return pending
}

// markAssigned settles the required groups key names through one of their own
// paths by the value the decoder assigns from it. With ZeroEmpty an empty value
// zeroes the field, so the last key to reach a field decides whether it is
// provided, not whether any of them was non-empty; the groups an empty value
// zeroed are recorded in decided.
func markAssigned(info *structInfo, key string, val []string, satisfied, decided []uint64, pending int) int {
	var buf [maxDirectKeyLen]byte
	lower := key
	if len(key) <= maxDirectKeyLen {
		if foldASCIILower(buf[:len(key)], key) {
			lower = utils.UnsafeString(buf[:len(key)])
		}
	} else {
		lower = utilstrings.ToLower(key)
	}
	for _, p := range info.foldedRequiredPaths[lower] {
		if !isFileValue(val) && isMultipartField(p.typ) {
			continue // a text key assigns nothing to a file field
		}
		word, bit := p.group>>6, uint64(1)<<uint(p.group&63)
		if isEmpty(p.typ, val) {
			decided[word] |= bit
			if satisfied[word]&bit != 0 {
				satisfied[word] &^= bit
				pending++
			}
		} else if satisfied[word]&bit == 0 {
			satisfied[word] |= bit
			pending--
		}
	}
	return pending
}

// markProvidedFoldedUndecided is markProvidedFolded for a decode in which
// markAssigned zeroed some groups: it leaves alone the ones the decode left
// zeroed, since a key in another case seen before the empty value no longer
// counts for them. decided is reused to hold those groups, so it is not
// meaningful afterwards.
func markProvidedFoldedUndecided(c *cache, info *structInfo, src *source, satisfied, decided []uint64, pending int) int {
	for w := range satisfied {
		decided[w] &^= satisfied[w]
	}
	pending = markProvidedFolded(c, info, src, satisfied, pending)
	for w := range satisfied {
		if again := satisfied[w] & decided[w]; again != 0 {
			satisfied[w] &^= again
			pending += bits.OnesCount64(again)
		}
	}
	return pending
}

// anyBit reports whether any bit of set is set.
func anyBit(set []uint64) bool {
	for _, w := range set {
		if w != 0 {
			return true
		}
	}
	return false
}

// markProvidedFolded matches src's keys case-insensitively against the groups still pending.
func markProvidedFolded(c *cache, info *structInfo, src *source, satisfied []uint64, pending int) int {
	if p := src.pairs; p != nil {
		var scratch [8]string
		for i, key := range p.keys {
			if !p.first(i) {
				continue
			}
			if pending = markProvidedByFoldedKey(c, info, key, p.group(i, scratch[:0]), satisfied, pending); pending == 0 {
				break
			}
		}
		return pending
	}
	for key, val := range src.m {
		if len(val) == 0 {
			continue
		}
		if pending = markProvidedByFoldedKey(c, info, key, val, satisfied, pending); pending == 0 {
			break
		}
	}
	return pending
}

// markProvidedByFoldedKey marks the groups key names in lower case; a lowercase
// key was looked up exactly, so it is skipped unless a required path is not.
func markProvidedByFoldedKey(c *cache, info *structInfo, key string, val []string, satisfied []uint64, pending int) int {
	var buf [maxDirectKeyLen]byte
	var lower string
	if len(key) <= maxDirectKeyLen {
		if !foldASCIILower(buf[:len(key)], key) && !info.requiredMixedCase {
			return pending
		}
		lower = utils.UnsafeString(buf[:len(key)])
	} else if lower = utilstrings.ToLower(key); lower == key && !info.requiredMixedCase {
		return pending
	}
	pending = markProvided(info.foldedRequiredPaths[lower], val, satisfied, pending)
	parsed := false
	for off := 0; pending > 0; {
		i := strings.IndexByte(lower[off:], '.')
		if i < 0 {
			break
		}
		// Lowercasing is ASCII-only, so the folded prefixes are no longer than maxRequiredPrefixLen either.
		if off += i + 1; off > info.maxRequiredPrefixLen {
			break
		}
		if owners := info.foldedRequiredPrefixes[lower[:off]]; len(owners) > 0 {
			if !parsed {
				if _, err := c.parsePathInfo(key, info); err != nil {
					return pending
				}
				parsed = true
			}
			pending = markProvided(owners, val, satisfied, pending)
		}
	}
	return pending
}

// markProvided marks the owners' groups that val is non-empty for; group g is bit g&63 of satisfied[g>>6].
func markProvided(owners []requiredPrefix, val []string, satisfied []uint64, pending int) int {
	for _, p := range owners {
		if satisfied[p.group>>6]&(1<<uint(p.group&63)) != 0 {
			continue
		}
		if !isEmpty(p.typ, val) {
			satisfied[p.group>>6] |= 1 << uint(p.group&63)
			pending--
		}
	}
	return pending
}

// missingRequired reports the required keys nothing in src answered.
func missingRequired(groups []requiredGroup, satisfied []uint64, pending int) MultiError {
	if pending == 0 {
		return nil
	}
	var errs MultiError
	for gi := range groups {
		if satisfied[gi>>6]&(1<<uint(gi&63)) == 0 {
			errs = appendError(errs, groups[gi].key, EmptyFieldError{Key: groups[gi].key})
		}
	}
	return errs
}

// checkRequired reports which of info's required keys src leaves empty: the
// standalone form of what Decode folds into its own loop.
func (d *Decoder) checkRequired(info *structInfo, src map[string][]string) MultiError {
	var inline [requiredBitWords]uint64
	s := source{m: src}
	satisfied, pending := markProvidedDirectly(info.requiredGroups, &s, inline[:])
	if pending > 0 {
		for key, val := range src {
			i := strings.IndexByte(key, '.')
			if i < 0 || len(val) == 0 {
				continue
			}
			if pending = markProvidedByNestedKey(d.cache, false, info, key, i+1, val, satisfied, pending); pending == 0 {
				break
			}
		}
	}
	if pending > 0 {
		pending = markProvidedFolded(d.cache, info, &s, satisfied, pending)
	}
	return missingRequired(info.requiredGroups, satisfied, pending)
}

// requiredGroup is one required key together with the fields that can
// satisfy it; its position in structInfo.requiredGroups is its bitset index.
type requiredGroup struct {
	key    string
	fields []fieldWithPrefix
}

// requiredPrefix names a group a nested source key can satisfy, with the
// field type that judges whether the key's value counts as non-empty.
type requiredPrefix struct {
	typ   reflect.Type
	group int
}

type fieldWithPrefix struct {
	*fieldInfo
	prefix string
	// searchPaths lists the src keys this required field answers to, and
	// searchPathDots the corresponding nested-key prefixes; both are
	// precomputed at cache-build time so per-request checks allocate
	// nothing. A field no path may end at has no searchPaths.
	searchPaths    []string
	searchPathDots []string
}

func newFieldWithPrefix(f *fieldInfo, prefix string) fieldWithPrefix {
	paths := f.paths(prefix)
	dots := make([]string, len(paths))
	for i, p := range paths {
		dots[i] = p + "."
	}
	if !f.endsPath() {
		paths = nil
	}
	return fieldWithPrefix{
		fieldInfo:      f,
		prefix:         prefix,
		searchPaths:    paths,
		searchPathDots: dots,
	}
}

// directlyProvided reports whether any of the group's own paths carries a
// non-empty value in src.
func directlyProvided(fields []fieldWithPrefix, src *source) bool {
	for _, f := range fields {
		for _, path := range f.searchPaths {
			if v, ok := src.lookup(path); ok && !isEmpty(f.typ, v) {
				return true
			}
		}
	}
	return false
}

// isEmpty returns true if value is empty for specific type, judging the last value as a scalar decodes it.
func isEmpty(t reflect.Type, value []string) bool {
	if len(value) == 0 {
		return true
	}
	switch t.Kind() {
	case boolType, float32Type, float64Type,
		intType, int8Type, int16Type, int32Type, int64Type,
		stringType,
		uintType, uint8Type, uint16Type, uint32Type, uint64Type:
		return len(value[len(value)-1]) == 0
	}
	if isMultipartField(t) {
		return !isFileValue(value) // only a file fills a file field
	}
	return false
}

// isFileValue reports whether value is the fileKeyValues Decode gives a file's key itself.
func isFileValue(value []string) bool {
	return len(value) > 0 && &value[0] == &fileKeyValues[0]
}

var (
	multipartFileHeaderPointerType      = reflect.TypeFor[*multipart.FileHeader]()
	sliceMultipartFileHeaderPointerType = reflect.TypeFor[[]*multipart.FileHeader]()
)

// Supported multiple types:
// *multipart.FileHeader, *[]multipart.FileHeader, []*multipart.FileHeader
func handleMultipartField(field reflect.Value, files []*multipart.FileHeader) bool {
	fieldType := field.Type()
	if !isMultipartField(fieldType) {
		return false
	}

	// Skip if files are empty and field is multipart
	if len(files) == 0 {
		return true
	}

	// Check for *multipart.FileHeader
	if fieldType == multipartFileHeaderPointerType {
		field.Set(reflect.ValueOf(files[0]))
		return true
	}

	// Check for []*multipart.FileHeader
	if fieldType == sliceMultipartFileHeaderPointerType {
		field.Set(reflect.ValueOf(files))
		return true
	}

	// Check for *[]*multipart.FileHeader
	if fieldType.Kind() == reflect.Pointer {
		fieldType = fieldType.Elem()

		if field.IsNil() {
			field.Set(reflect.New(fieldType))
		}

		if fieldType == sliceMultipartFileHeaderPointerType {
			field.Elem().Set(reflect.ValueOf(files))
			return true
		}
	}

	return false
}

// Supported multiple types:
// *multipart.FileHeader, *[]multipart.FileHeader, []*multipart.FileHeader
func isMultipartField(typ reflect.Type) bool {
	// Check for *multipart.FileHeader
	if typ == multipartFileHeaderPointerType {
		return true
	}

	// Check for []*multipart.FileHeader
	if typ == sliceMultipartFileHeaderPointerType {
		return true
	}

	// Check for *[]*multipart.FileHeader
	if typ.Kind() == reflect.Pointer {
		typ = typ.Elem()

		if typ == sliceMultipartFileHeaderPointerType {
			return true
		}
	}

	return false
}

// walkHops follows a path to the target field, dereferencing pointers and
// allocating the embedded pointers promoted fields need on the way. It
// returns the zero Value when an unsettable nil embedded pointer blocks the
// walk, and an error once the call's allocation budget is spent.
func walkHops(v reflect.Value, hops []pathHop, grow *growTracker) (reflect.Value, error) {
	var err error
	for _, hop := range hops {
		// A previous hop may have been blocked by an unsettable nil
		// embedded pointer; the field is unreachable then.
		if !v.IsValid() {
			return v, nil
		}
		if v.Kind() == reflect.Pointer {
			if v.IsNil() {
				if err = grow.alloc(v, v.Type().Elem()); err != nil {
					return v, err
				}
			}
			v = v.Elem()
		}

		// Allocate embedded anonymous pointers required for promoted fields.
		for _, idx := range hop.ensure {
			if f := v.Field(idx); f.IsNil() {
				if err = grow.alloc(f, f.Type().Elem()); err != nil {
					return v, err
				}
			}
		}

		// A field of v itself needs none of the chain's pointer handling.
		if len(hop.index) == 1 {
			v = v.Field(hop.index[0])
		} else if v, err = walkIndexChain(v, hop.index, grow); err != nil {
			return v, err
		}
	}
	return v, nil
}

// walkIndexChain walks v along a struct field index chain. Chains longer
// than one element traverse embedded structs; intermediate nil pointers are
// allocated so promoted fields stay reachable. It returns the zero Value
// when the chain is blocked by a nil pointer that cannot be set (an
// unexported embedded pointer), which callers treat as an unreachable
// field, and an error once the call's allocation budget is spent.
func walkIndexChain(v reflect.Value, chain []int, grow *growTracker) (reflect.Value, error) {
	for j, fi := range chain {
		if j > 0 && v.Kind() == reflect.Pointer {
			if v.IsNil() {
				if !v.CanSet() {
					return reflect.Value{}, nil
				}
				if err := grow.alloc(v, v.Type().Elem()); err != nil {
					return v, err
				}
			}
			v = v.Elem()
		}
		v = v.Field(fi)
	}
	return v, nil
}

// growTracker counts what a Decode call allocates toward maxAllocBytes and lets the
// first eight arrays it grows at a key's first index grow geometrically.
type growTracker struct {
	arrays [8]unsafe.Pointer
	n      int
	used   int
}

// charge counts n values of size sz toward maxAllocBytes. Checked before
// adding, so the count cannot wrap and what is refused costs nothing.
func (g *growTracker) charge(n int, sz uintptr) error {
	if s := int(sz); s != 0 && n > (maxAllocBytes-g.used)/s {
		return errAllocLimit
	}
	g.used += n * int(sz)
	return nil
}

// alloc points the nil pointer v at a new t, counted unless g is nil.
func (g *growTracker) alloc(v reflect.Value, t reflect.Type) error {
	if g != nil {
		if err := g.charge(1, t.Size()); err != nil {
			return err
		}
	}
	v.Set(reflect.New(t))
	return nil
}

// slot returns the entry of the array p if this call allocated it, else the
// next free entry, or -1 when every entry is taken.
func (g *growTracker) slot(p unsafe.Pointer) int {
	for i := range g.n {
		if g.arrays[i] == p {
			return i
		}
	}
	if g.n == len(g.arrays) {
		return -1
	}
	return g.n
}

func (g *growTracker) record(i int, p unsafe.Pointer) {
	g.arrays[i] = p
	g.n = max(g.n, i+1)
}

// decode fills a struct field using a parsed path. nested is set below a key's
// first index; the outermost call is the one place slice growth can be tracked.
func (d *Decoder) decode(v reflect.Value, path string, parts []pathPart, values []string, files []*multipart.FileHeader, grow *growTracker, nested bool) error {
	// Get the field walking the struct fields by index. Almost every path is
	// one hop into a field of v, which needs none of the loop's bookkeeping.
	if idx := parts[0].soleIndex; idx >= 0 && v.Kind() == reflect.Struct {
		v = v.Field(idx)
	} else {
		var err error
		if v, err = walkHops(v, parts[0].hops, grow); err != nil || !v.IsValid() {
			// Past the budget, or unreachable behind an unsettable nil embedded pointer.
			return err
		}
	}

	// Don't even bother for unexported fields.
	if !v.CanSet() {
		return nil
	}

	// Fast path: plain builtin scalar fields skip the converter/unmarshaler
	// dispatch below; fastKind was validated at cache-build time.
	if k := parts[0].field.fastKind; k != reflect.Invalid && len(parts) == 1 && !parts[0].elem {
		val := ""
		if len(values) > 0 {
			val = values[len(values)-1]
		}
		if val == "" {
			if d.zeroEmpty {
				v.SetZero()
			}
			return nil
		}
		if _, ok := setBuiltinKind(v, k, val); !ok {
			return ConversionError{
				Key:   path,
				Type:  parts[0].field.typ,
				Index: -1,
			}
		}
		return nil
	}

	// Check multipart files
	if parts[0].field.isMultipart && handleMultipartField(v, files) {
		return nil
	}

	if parts[0].field.mapValue != nil {
		return d.decodeMapEntry(v, path, &parts[0], values)
	}

	// Dereference if needed.
	t := v.Type()
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
		if v.IsNil() {
			if err := grow.alloc(v, t); err != nil {
				return err
			}
		}
		v = v.Elem()
	}

	// Slice of structs. Let's go recursive.
	if len(parts) > 1 {
		idx := parts[0].index
		// a defensive check to avoid creating a large slice based on user input index
		if idx > d.maxSize {
			return fmt.Errorf("%v index %d is larger than the configured maxSize %d", v.Kind(), idx, d.maxSize)
		}
		if n := idx + 1; v.IsNil() || v.Len() < n {
			// Indices arrive in map order, so a slice is typically grown
			// several times per call. Extending one this call allocated is
			// free: the room past its length is ours and freshly zeroed.
			slot, root := -1, false
			if !nested {
				slot, root = grow.slot(v.UnsafePointer()), parts[0].soleIndex >= 0
			}
			ours := slot >= 0 && slot < grow.n
			if ours && n <= v.Cap() {
				v.SetLen(n)
			} else {
				// Otherwise grow into a fresh array, as the caller may still alias the old one;
				// a root field reserves room from its first growth, other slices from their second.
				c := d.growCap(n, slot >= 0 && (root || ours))
				// Only a tracked slice field of the destination struct is left to its index limit.
				if slot < 0 || !root {
					if err := grow.charge(c, t.Elem().Size()); err != nil {
						return err
					}
				}
				value := reflect.MakeSlice(t, n, c)
				if v.Len() > 0 {
					// Resize it.
					reflect.Copy(value, v)
				}
				v.Set(value)
				if slot >= 0 {
					grow.record(slot, value.UnsafePointer())
				}
			}
		}
		return d.decode(v.Index(idx), path, parts[1:], values, files, grow, true)
	}

	// Get the converter early in case there is one for a slice type.
	conv := d.cache.converter(t)
	// The encoding.TextUnmarshaler facts for v's type are precomputed per
	// field; instances are bound to live values where needed below. For an
	// elem part (path terminated at a slice index) v is an element of the
	// slice field, so the element type's facts apply.
	m := parts[0].field.derefUnmarshaler
	if parts[0].elem {
		m = parts[0].field.elemUnmarshaler
	}
	if conv == nil && t.Kind() == reflect.Slice && m.IsSliceElement {
		elemT := t.Elem()
		isPtrElem := elemT.Kind() == reflect.Pointer
		if isPtrElem {
			elemT = elemT.Elem()
		}

		// Try to get a converter for the element type.
		customConv := d.cache.converter(elemT)
		var conv valueConverter
		if customConv != nil {
			conv = *customConv
		} else if builtin := getBuiltinConverter(elemT.Kind()); builtin != nil {
			conv.plain = builtin
		} else {
			// As we are not dealing with slice of structs here, we don't need to check if the type
			// implements TextUnmarshaler interface
			return fmt.Errorf("schema: converter not found for %v", elemT)
		}

		// Fast path: builtin element kinds without unmarshalers, custom
		// converters or pointer elements decode straight into a fresh slice,
		// avoiding one reflect.Value allocation per element.
		if customConv == nil && !m.IsValid && !isPtrElem {
			return d.decodeBuiltinSlice(v, t, path, values)
		}

		if m.IsValid && !isPtrElem {
			return d.decodeUnmarshalerSlice(v, t, path, values)
		}

		return d.decodeBoxedSlice(v, t, elemT, path, values, conv, m, isPtrElem)
	} else {
		val := ""
		// Use the last value provided if any values were provided
		if len(values) > 0 {
			val = values[len(values)-1]
		}

		if conv != nil {
			value, err := conv.convert(val)
			if err != nil || !value.IsValid() {
				return ConversionError{
					Key:   path,
					Type:  t,
					Index: -1,
					Err:   err,
				}
			}
			v.Set(value.Convert(t))
		} else if m.IsValid {
			if m.IsPtr {
				// A zero v is what reflect.New would hand UnmarshalText, so
				// decode into it in place and save the allocation; a failed
				// call leaves it zero again, as it found it.
				inPlace := v.IsZero()
				u := v.Addr()
				if !inPlace {
					u = reflect.New(v.Type())
				}
				um, _ := reflect.TypeAssert[encoding.TextUnmarshaler](u)
				if err := um.UnmarshalText([]byte(val)); err != nil {
					if inPlace {
						v.SetZero()
					}
					return ConversionError{
						Key:   path,
						Type:  t,
						Index: -1,
						Err:   err,
					}
				}
				if !inPlace {
					v.Set(reflect.Indirect(u))
				}
			} else {
				// If the value implements the encoding.TextUnmarshaler interface
				// apply UnmarshalText as the converter, binding it to the
				// live value.
				um, _ := reflect.TypeAssert[encoding.TextUnmarshaler](v)
				if err := um.UnmarshalText([]byte(val)); err != nil {
					return ConversionError{
						Key:   path,
						Type:  t,
						Index: -1,
						Err:   err,
					}
				}
			}
		} else if val == "" {
			// Clear only what a value could set: a struct's own key must not
			// reach the fields no key names.
			if d.zeroEmpty && getBuiltinConverter(t.Kind()) != nil {
				v.Set(reflect.Zero(t))
			}
		} else if handled, ok := setBuiltinKind(v, t.Kind(), val); handled {
			if !ok {
				return ConversionError{
					Key:   path,
					Type:  t,
					Index: -1,
				}
			}
		} else {
			return fmt.Errorf("schema: converter not found for %v", t)
		}
	}
	return nil
}

type mapEntry struct {
	key, val reflect.Value
}

func newMapEntryPool(t reflect.Type) *sync.Pool {
	return &sync.Pool{New: func() any {
		return &mapEntry{key: reflect.New(t.Key()).Elem(), val: reflect.New(t.Elem()).Elem()}
	}}
}

// decodeMapEntry decodes values into the entry of the map field v that the
// rest of path from part.keyStart names, allocating the map on first use. The value is decoded as
// a field of the map's value type, so it converts, reports errors and takes
// the last of several values the way such a field does; the key is stored
// exactly as sent.
func (d *Decoder) decodeMapEntry(v reflect.Value, path string, part *pathPart, values []string) error {
	f := part.field
	// A plain scalar field ignores an empty value unless ZeroEmpty is set, so
	// such a value leaves the map as it is rather than storing a zero value.
	if k := f.mapValue[0].field.fastKind; k != reflect.Invalid && !d.zeroEmpty &&
		(len(values) == 0 || values[len(values)-1] == "") {
		return nil
	}
	keyStart := int(part.keyStart)
	if keyStart == 0 {
		keyStart = strings.IndexByte(path, '.') + 1
	}
	e, _ := f.entries.Get().(*mapEntry)
	err := d.decode(e.val, path, f.mapValue, values, nil, nil, false)
	if err == nil {
		if v.IsNil() {
			v.Set(reflect.MakeMap(v.Type()))
		}
		// The key may alias the caller's buffer.
		e.key.SetString(strings.Clone(path[keyStart:]))
		v.SetMapIndex(e.key, e.val)
	}
	// The map holds copies; zero the pair so the pool keeps none of the request's data.
	e.val.SetZero()
	e.key.SetZero()
	f.entries.Put(e)
	return err
}

// decodeBoxedSlice decodes values into the slice field v whose elements have
// to round-trip through a reflect.Value: a custom converter, or pointer
// elements. Each item goes straight into a slot of a detached slice, set on v
// only if all convert.
func (d *Decoder) decodeBoxedSlice(v reflect.Value, t, elemT reflect.Type, path string, values []string, conv valueConverter, m unmarshaler, isPtrElem bool) error {
	limit := len(values) // a value yields at most one item more than it has commas
	for _, value := range values {
		limit += strings.Count(value, ",")
	}
	var sl reflect.Value
	room, n := 0, 0
	slot := func() reflect.Value {
		if n == room {
			sl = growSlice(sl, t, limit)
			room = sl.Len()
		}
		n++
		return sl.Index(n - 1)
	}

	for key, value := range values {
		if value == "" {
			if d.zeroEmpty {
				slot() // stays zero
			}
			continue
		}
		if m.IsValid {
			// Only pointer elements get here, so u itself is the element.
			u := reflect.New(elemT)
			um, _ := reflect.TypeAssert[encoding.TextUnmarshaler](u)
			if err := um.UnmarshalText([]byte(value)); err != nil {
				return ConversionError{
					Key:   path,
					Type:  t,
					Index: key,
					Err:   err,
				}
			}
			slot().Set(u)
			continue
		}
		// conv.convert, spelled out: it is too big to inline, and this
		// runs once per element.
		var item reflect.Value
		var err error
		if conv.withErr == nil {
			item = conv.plain(value)
		} else {
			item, err = conv.withErr(value)
		}
		if err == nil && item.IsValid() {
			setConvertedItem(slot(), item, elemT, isPtrElem)
			continue
		}
		if strings.IndexByte(value, ',') == -1 {
			return ConversionError{
				Key:   path,
				Type:  elemT,
				Index: key,
				Err:   err,
			}
		}
		// The value did not convert as a whole, so it is retried as a
		// comma-separated list; the whole value's error is dropped, and an
		// item that fails reports its own.
		for value := range strings.SplitSeq(value, ",") {
			if value == "" {
				if d.zeroEmpty {
					slot()
				}
			} else if item, err := conv.convert(value); err == nil && item.IsValid() {
				setConvertedItem(slot(), item, elemT, isPtrElem)
			} else {
				return ConversionError{
					Key:   path,
					Type:  elemT,
					Index: key,
					Err:   err,
				}
			}
		}
	}
	if room == 0 {
		sl = reflect.MakeSlice(t, 0, 0)
	} else if n < room {
		sl = sl.Slice(0, n)
	}
	v.Set(sl)
	return nil
}

// setConvertedItem converts a builtin/custom converter result to the slice
// element type and stores it in slot, wrapping it in a freshly allocated pointer for
// pointer-element slices. The conversion must happen before the pointer wrap:
// builtin converters return the underlying kind (e.g. int for a named
// `type MyInt int`), which is not assignable to the named element type, so
// Set-ing it into a *MyInt without converting first panics.
func setConvertedItem(slot, item reflect.Value, elemT reflect.Type, isPtrElem bool) {
	if item.Type() != elemT {
		item = item.Convert(elemT)
	}
	if isPtrElem {
		ptr := reflect.New(elemT)
		ptr.Elem().Set(item)
		item = ptr
	}
	slot.Set(item)
}

// sliceBudget bounds the first allocation: the client chooses how many values arrive.
const sliceBudget = 4 << 10

// decodeUnmarshalerSlice unmarshals each value into an element of a detached slice,
// set on v only if all succeed; empty values are skipped unless ZeroEmpty is set.
func (d *Decoder) decodeUnmarshalerSlice(v reflect.Value, t reflect.Type, path string, values []string) error {
	var sl reflect.Value
	room, i := 0, 0
	for key, value := range values {
		if value == "" && !d.zeroEmpty {
			continue
		}
		if i == room {
			sl = growSlice(sl, t, len(values))
			room = sl.Len()
		}
		if value != "" { // an empty value kept for ZeroEmpty leaves its slot zero
			um, _ := reflect.TypeAssert[encoding.TextUnmarshaler](sl.Index(i).Addr())
			if err := um.UnmarshalText([]byte(value)); err != nil {
				return ConversionError{
					Key:   path,
					Type:  t,
					Index: key,
					Err:   err,
				}
			}
		}
		i++
	}
	if room == 0 {
		sl = reflect.MakeSlice(t, 0, 0)
	} else if i < room {
		sl = sl.Slice(0, i)
	}
	v.Set(sl)
	return nil
}

// growSlice returns sl with more room: the budget at first, then doubled, up to limit.
func growSlice(sl reflect.Value, t reflect.Type, limit int) reflect.Value {
	var n int
	if sl.IsValid() {
		n = 2 * sl.Len()
	} else if size := int(t.Elem().Size()); size == 0 {
		n = limit
	} else {
		n = max(1, sliceBudget/size)
	}
	n = min(n, limit)
	grown := reflect.MakeSlice(t, n, n)
	if sl.IsValid() {
		reflect.Copy(grown, sl)
	}
	return grown
}

// decodeBuiltinSlice decodes values into the slice field v of type t whose
// elements are builtin-convertible kinds, parsing directly into slice slots
// instead of boxing every element in a reflect.Value. The slice is built
// detached and only assigned to v when every value parsed, matching the
// all-or-nothing behavior of the generic path.
//
// A value that fails to parse as a whole is retried as a comma-separated
// list. For non-string kinds a value containing a comma can never parse as a
// whole (no builtin syntax admits commas — pinned by a test), and string
// values always parse, so item boundaries are knowable upfront: the slice is
// sized by a cheap comma count (an upper bound, since empty items may be
// skipped) and truncated to the filled length at the end.
func (d *Decoder) decodeBuiltinSlice(v reflect.Value, t reflect.Type, path string, values []string) error {
	elemT := t.Elem()
	k := elemT.Kind()
	split := k != reflect.String

	n := len(values)
	if split {
		for _, value := range values {
			n += strings.Count(value, ",")
		}
		// The counts already say whether any value holds a separator, so
		// when none does the per-value scans below have nothing to find.
		split = n > len(values)
	}

	// Exact builtin slice types decode without per-element reflect calls;
	// named slice or element types fall through to the generic path.
	switch t {
	case typSliceString:
		return decodeNativeSlice(d.zeroEmpty, v, path, values, elemT, n, split, parseNativeString)
	case typSliceInt:
		return decodeNativeSlice(d.zeroEmpty, v, path, values, elemT, n, split, parseNativeInt)
	case typSliceInt64:
		return decodeNativeSlice(d.zeroEmpty, v, path, values, elemT, n, split, parseNativeInt64)
	case typSliceUint:
		return decodeNativeSlice(d.zeroEmpty, v, path, values, elemT, n, split, parseNativeUint)
	case typSliceUint64:
		return decodeNativeSlice(d.zeroEmpty, v, path, values, elemT, n, split, parseNativeUint64)
	case typSliceFloat64:
		return decodeNativeSlice(d.zeroEmpty, v, path, values, elemT, n, split, parseNativeFloat64)
	case typSliceBool:
		return decodeNativeSlice(d.zeroEmpty, v, path, values, elemT, n, split, parseNativeBool)
	}

	sl := reflect.MakeSlice(t, n, n)
	i := 0
	for key, value := range values {
		switch {
		case value == "":
			if d.zeroEmpty {
				i++ // slot stays zero
			}
		case split && strings.IndexByte(value, ',') != -1:
			for rest, more := value, true; more; {
				var item string
				item, rest, more = utils.CutByte(rest, ',')
				if item == "" {
					if d.zeroEmpty {
						i++ // slot stays zero
					}
					continue
				}
				if _, ok := setBuiltinKind(sl.Index(i), k, item); !ok {
					return ConversionError{
						Key:   path,
						Type:  elemT,
						Index: key,
					}
				}
				i++
			}
		default:
			if _, ok := setBuiltinKind(sl.Index(i), k, value); !ok {
				return ConversionError{
					Key:   path,
					Type:  elemT,
					Index: key,
				}
			}
			i++
		}
	}
	if i < n {
		sl = sl.Slice(0, i)
	}
	v.Set(sl)
	return nil
}

// Exact (unnamed) builtin slice types eligible for the native decode path.
var (
	typSliceString  = reflect.TypeFor[[]string]()
	typSliceInt     = reflect.TypeFor[[]int]()
	typSliceInt64   = reflect.TypeFor[[]int64]()
	typSliceUint    = reflect.TypeFor[[]uint]()
	typSliceUint64  = reflect.TypeFor[[]uint64]()
	typSliceFloat64 = reflect.TypeFor[[]float64]()
	typSliceBool    = reflect.TypeFor[[]bool]()
)

// decodeNativeSlice mirrors the generic decodeBuiltinSlice loop for a field
// typed exactly []T, parsing into a native slice assigned only if every value
// parsed (all-or-nothing). n is the precomputed element upper bound.
func decodeNativeSlice[T any](zeroEmpty bool, v reflect.Value, path string, values []string, elemT reflect.Type, n int, split bool, parse func(string) (T, bool)) error {
	out := make([]T, 0, n)
	var zero T
	for key, value := range values {
		switch {
		case value == "":
			if zeroEmpty {
				out = append(out, zero)
			}
		case split && strings.IndexByte(value, ',') != -1:
			for rest, more := value, true; more; {
				var item string
				item, rest, more = utils.CutByte(rest, ',')
				if item == "" {
					if zeroEmpty {
						out = append(out, zero)
					}
					continue
				}
				ev, ok := parse(item)
				if !ok {
					return ConversionError{
						Key:   path,
						Type:  elemT,
						Index: key,
					}
				}
				out = append(out, ev)
			}
		default:
			ev, ok := parse(value)
			if !ok {
				return ConversionError{
					Key:   path,
					Type:  elemT,
					Index: key,
				}
			}
			out = append(out, ev)
		}
	}
	// v's type is exactly []T here (the dispatch switch guarantees it), and v
	// is settable (decode returns before reaching here for a field it cannot
	// set), so assign through a typed pointer to it: no reflect.ValueOf
	// escape and no Set assignability checks. The pointer comes from
	// UnsafeAddr rather than Addr().Interface(), which resolved the pointer
	// type on every call and measured a fifth of the whole assignment's cost.
	*(*[]T)(unsafe.Pointer(v.UnsafeAddr())) = out
	return nil
}

func isTextUnmarshaler(v reflect.Value) unmarshaler {
	// Create a new unmarshaller instance
	m := unmarshaler{}
	if _, m.IsValid = reflect.TypeAssert[encoding.TextUnmarshaler](v); m.IsValid {
		return m
	}
	// As the UnmarshalText function should be applied to the pointer of the
	// type, we check that type to see if it implements the necessary
	// method.
	if _, m.IsValid = reflect.TypeAssert[encoding.TextUnmarshaler](reflect.New(v.Type())); m.IsValid {
		m.IsPtr = true
		return m
	}

	// if v is []T or *[]T create new T
	t := v.Type()
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	if t.Kind() == reflect.Slice {
		// The slice type itself cannot implement encoding.TextUnmarshaler
		// here: the value-level assert above already covered it. Check
		// whether the elements do.
		m.IsSliceElement = true
		if t = t.Elem(); t.Kind() == reflect.Pointer {
			t = reflect.PointerTo(t.Elem())
			m.IsSliceElementPtr = true
			_, m.IsValid = reflect.TypeAssert[encoding.TextUnmarshaler](reflect.Zero(t))
			return m
		}
	}

	_, m.IsValid = reflect.TypeAssert[encoding.TextUnmarshaler](reflect.New(t))
	return m
}

// TextUnmarshaler helpers ----------------------------------------------------
// unmarshaler describes how a type relates to encoding.TextUnmarshaler.
// It carries type-level facts only; the decoder binds instances to live
// values at the point of use.
type unmarshaler struct {
	// IsValid indicates whether the resolved type indicated by the other
	// flags implements the encoding.TextUnmarshaler interface.
	IsValid bool
	// IsPtr indicates that the resolved type is the pointer of the original
	// type.
	IsPtr bool
	// IsSliceElement indicates that the resolved type is a slice element of
	// the original type.
	IsSliceElement bool
	// IsSliceElementPtr indicates that the resolved type is a pointer to a
	// slice element of the original type.
	IsSliceElementPtr bool
}

// Errors ---------------------------------------------------------------------

// ConversionError stores information about a failed conversion.
type ConversionError struct {
	Key   string       // key from the source map.
	Type  reflect.Type // expected type of elem
	Index int          // index for multi-value fields; -1 for single-value fields.
	Err   error        // low-level error (when it exists)
}

// The error strings below are assembled by concatenation instead of
// fmt.Sprintf: %q is strconv.Quote and %d is utils.FormatInt, so the messages
// are byte-identical without the reflection-based formatter.
func (e ConversionError) Error() string {
	var output string

	if e.Index < 0 {
		output = "schema: error converting value for " + strconv.Quote(e.Key)
	} else {
		output = "schema: error converting value for index " +
			utils.FormatInt(int64(e.Index)) + " of " + strconv.Quote(e.Key)
	}

	if e.Err != nil {
		output += ". Details: " + e.Err.Error()
	}

	return output
}

// Unwrap returns Err, so errors.Is and errors.As reach the converter's error.
func (e ConversionError) Unwrap() error {
	return e.Err
}

// UnknownKeyError stores information about an unknown key in the source map.
type UnknownKeyError struct {
	Key string // key from the source map.
}

func (e UnknownKeyError) Error() string {
	return "schema: invalid path " + strconv.Quote(e.Key)
}

// EmptyFieldError stores information about an empty required field.
type EmptyFieldError struct {
	Key string // required key in the source map.
}

func (e EmptyFieldError) Error() string {
	return e.Key + " is empty"
}

// MultiError stores multiple decoding errors.
//
// Borrowed from the App Engine SDK.
type MultiError map[string]error

func (e MultiError) Error() string {
	s := ""
	for _, err := range e {
		s = err.Error()
		break
	}
	switch len(e) {
	case 0:
		return "(0 errors)"
	case 1:
		return s
	case 2:
		return s + " (and 1 other error)"
	}
	return s + " (and " + utils.FormatInt(int64(len(e)-1)) + " other errors)"
}

func appendError(m MultiError, key string, err error) MultiError {
	if err == nil {
		return m
	}
	if m == nil {
		m = make(MultiError)
	}
	if m[key] == nil {
		m[key] = err
	}
	return m
}

func mergeErrors(dst, src MultiError) MultiError {
	if len(src) == 0 {
		return dst
	}
	if dst == nil {
		dst = make(MultiError, len(src))
	}
	for key, err := range src {
		if dst[key] == nil {
			dst[key] = err
		}
	}
	return dst
}
