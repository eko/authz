package schema

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"sync"
	"sync/atomic"

	utils "github.com/gofiber/utils/v2"
)

type encoderFunc func(reflect.Value) string

// errNotStruct is returned by Encode for invalid sources; hoisted so the
// check does not allocate on every call.
var errNotStruct = errors.New("schema: interface must be a struct")

// errNilDst is returned by Encode when the destination map is nil, which
// would otherwise panic on the first map assignment.
var errNilDst = errors.New("schema: dst map must not be nil")

// maxScratchValueLen bounds values stored in encode's shared scratch array,
// capping how much data a surviving dst entry can keep reachable after the
// caller deletes neighboring keys. Only strings allocated by this package's
// own formatters are eligible at all (see encField.scratchSafe), so the cap
// bounds real bytes, not just headers.
const maxScratchValueLen = 64

// Encoder encodes values from a struct into url.Values.
type Encoder struct {
	cache  *cache
	regenc map[reflect.Type]encoderFunc
	// encCache memoizes the per-struct-type encoding plan
	// (map[reflect.Type][]encField) so tags are parsed and encoder
	// functions resolved once per type instead of on every Encode call.
	encCache sync.Map
	// encGen is bumped before encCache is cleared on configuration changes;
	// structInfo snapshots it before building a plan and refuses to store
	// the plan if it changed, so a build racing a reconfiguration cannot
	// re-insert a stale plan after the clear.
	encGen atomic.Uint64
	// encSlots is a direct-mapped cache in front of encCache, indexed like cache.slots.
	encSlots [typeCacheSlots]atomic.Pointer[encPlan]
}

// encPlan tags a per-type encoding plan with the configuration generation it
// was built under: hit-validation ignores plans from older generations, so
// any Encode starting after a reconfiguration returns observes the new
// configuration even if a racing build stored a stale plan after the clear.
type encPlan struct {
	typ    reflect.Type
	fields []encField
	gen    uint64
	// scratch is how many values encode can put in its shared scratch array.
	scratch int
	// freshKeys reports that encoding a value of this type into an empty map
	// writes every key exactly once: no two fields share a name, and no field
	// is recursed into — a nested struct's keys land in the same map, without
	// a prefix, so they could collide with these. encode can then assign each
	// key outright instead of reading it back to append.
	freshKeys bool
}

// encField is the precomputed encoding plan for one struct field.
type encField struct {
	name      string
	enc       encoderFunc // immediate encoder; nil for structs and slices
	elemEnc   encoderFunc // slice element encoder, when the field is a slice
	idx       int
	omitEmpty bool
	// recurseStructPtr marks pointer-to-struct fields without a custom
	// encoder: non-nil values are encoded by recursing into the element.
	recurseStructPtr bool
	// scratchSafe marks fields whose encoder output is allocated by this
	// package (numeric/bool/float formatters), so it can never alias a large
	// caller-owned buffer and may be batched in encode's shared scratch.
	scratchSafe bool
	// hasIsZero marks a struct-typed field whose type decides omitempty for
	// itself through an IsZero method. Settling that here keeps the check off
	// reflect.Value.Interface, which copies the struct to the heap every time
	// it is asked — whether or not the method is there.
	hasIsZero bool
	isStruct  bool
	// nilAsNull marks pointer fields whose element has no immediate
	// encoder (structs recursed via recurseStructPtr, or unsupported
	// types): nil values encode as "null", matching the closure behavior
	// pointer fields with encodable elements get.
	nilAsNull bool
	// elemPtrNil marks slice fields whose element is a pointer type with no
	// encoder (e.g. []*Struct): nil elements encode as "null" (as they did
	// historically), while a non-nil such element is an error.
	elemPtrNil bool
}

// NewEncoder returns a new Encoder with defaults.
func NewEncoder() *Encoder {
	return &Encoder{cache: newCache(), regenc: make(map[reflect.Type]encoderFunc)}
}

// Encode encodes a struct into map[string][]string.
//
// Intended for use with url.Values.
func (e *Encoder) Encode(src any, dst map[string][]string) (err error) {
	if dst == nil {
		return errNilDst
	}

	// Catch panics from reflection or user-registered encoders and return
	// them as an error instead of crashing the caller, mirroring Decode.
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(error); ok {
				err = e
			} else {
				err = fmt.Errorf("schema: panic while encoding: %v", r)
			}
		}
	}()

	v := reflect.ValueOf(src)

	return e.encode(v, dst)
}

// RegisterEncoder registers a converter for encoding a custom type.
func (e *Encoder) RegisterEncoder(value any, encoder func(reflect.Value) string) {
	e.cache.l.Lock()
	e.regenc[reflect.TypeOf(value)] = encoder
	e.cache.l.Unlock()
	e.resetPlans()
}

// SetAliasTag changes the tag used to locate custom field aliases.
// The default tag is "schema".
func (e *Encoder) SetAliasTag(tag string) {
	e.cache.l.Lock()
	e.cache.tag = tag
	e.cache.l.Unlock()
	e.resetPlans()
}

func (e *Encoder) resetPlans() {
	e.encGen.Add(1)
	e.encCache.Clear()
	for i := range e.encSlots {
		e.encSlots[i].Store(nil)
	}
}

// structInfo returns the cached encoding plan for struct type t, building it
// on first use. The build reads the tag and registered encoders under the
// configuration lock; the generation re-checks around the cache store keep a
// build racing a reconfiguration from inserting a stale plan.
func (e *Encoder) structInfo(t reflect.Type) *encPlan {
	gen := e.encGen.Load()
	slot := &e.encSlots[typeSlot(t)]
	cur := slot.Load()
	if cur != nil && cur.typ == t && cur.gen == gen {
		return cur
	}
	if cached, ok := e.encCache.Load(t); ok {
		// Ignore plans built under an older configuration; fall through and
		// rebuild (the fresh plan overwrites the stale entry).
		if p := cached.(*encPlan); p.gen == gen {
			claimPlanSlot(slot, cur, p)
			return p
		}
	}
	e.cache.l.RLock()
	tag := e.cache.tag
	fields := make([]encField, 0, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		sf := t.Field(i)
		name, opts := fieldAlias(sf, tag)
		if name == "-" {
			continue
		}
		ft := sf.Type
		f := encField{
			idx:       i,
			name:      name,
			omitEmpty: opts.Contains("omitempty"),
			recurseStructPtr: ft.Kind() == reflect.Pointer &&
				ft.Elem().Kind() == reflect.Struct &&
				!e.hasCustomEncoder(ft),
			enc:         typeEncoder(ft, e.regenc),
			scratchSafe: scratchSafeEncoder(ft, e.regenc),
			hasIsZero:   ft.Kind() == reflect.Struct && ft.Implements(zeroerType),
		}
		if f.enc == nil {
			switch ft.Kind() {
			case reflect.Struct:
				f.isStruct = true
			case reflect.Slice:
				f.elemEnc = typeEncoder(ft.Elem(), e.regenc)
				if f.elemEnc == nil && ft.Elem().Kind() == reflect.Pointer {
					f.elemPtrNil = true
				}
			case reflect.Pointer:
				f.nilAsNull = true
			}
		}
		fields = append(fields, f)
	}
	e.cache.l.RUnlock()
	p := &encPlan{
		typ:       t,
		fields:    fields,
		gen:       gen,
		scratch:   scratchValues(fields),
		freshKeys: writesEachKeyOnce(fields),
	}
	// Don't cache a plan whose inputs (tag, registered encoders) changed
	// while it was being built; the next call rebuilds it fresh. Even if a
	// stale plan slips in after the clear, its generation tag keeps it from
	// ever being served.
	if e.encGen.Load() == gen {
		e.encCache.Store(t, p)
		claimPlanSlot(slot, cur, p)
	}
	return p
}

// claimPlanSlot keeps another type's live plan in slot, so types sharing it don't evict each other.
func claimPlanSlot(slot *atomic.Pointer[encPlan], cur, p *encPlan) {
	if cur == nil || cur.gen < p.gen {
		slot.CompareAndSwap(cur, p)
	}
}

func scratchValues(fields []encField) int {
	n := 0
	for i := range fields {
		if f := &fields[i]; f.scratchSafe || f.nilAsNull {
			n++
		}
	}
	return n
}

// writesEachKeyOnce reports whether the plan's fields write distinct keys and
// none recurses into a nested struct, whose keys land in the same map and
// could repeat one of these.
func writesEachKeyOnce(fields []encField) bool {
	names := make(map[string]struct{}, len(fields))
	for i := range fields {
		f := &fields[i]
		if f.isStruct || f.recurseStructPtr {
			return false
		}
		if _, dup := names[f.name]; dup {
			return false
		}
		names[f.name] = struct{}{}
	}
	return true
}

// zeroer is the optional method a type can provide to decide, for omitempty,
// whether a value of it counts as empty.
type zeroer interface{ IsZero() bool }

var zeroerType = reflect.TypeFor[zeroer]()

// isZeroValue applies omitempty to one of this field's values: isZero with
// the struct case answered from the plan, so neither outcome goes through
// reflect.Value.Interface. An addressable struct becomes an interface through
// its address, without a copy, and a type known to lack IsZero skips the
// conversion altogether.
func (f *encField) isZeroValue(v reflect.Value) bool {
	if v.Kind() != reflect.Struct || !v.CanInterface() {
		return isZero(v)
	}
	if f.hasIsZero {
		if v.CanAddr() {
			// A value receiver's IsZero is in the pointer's method set too.
			iz, _ := reflect.TypeAssert[zeroer](v.Addr())
			return iz.IsZero()
		}
		iz, _ := reflect.TypeAssert[zeroer](v)
		return iz.IsZero()
	}
	for i := 0; i < v.NumField(); i++ {
		if !isZero(v.Field(i)) {
			return false
		}
	}
	return true
}

func isZero(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Func:
	case reflect.Map, reflect.Slice:
		return v.IsNil() || v.Len() == 0
	case reflect.Array:
		for i := 0; i < v.Len(); i++ {
			if !isZero(v.Index(i)) {
				return false
			}
		}
		return true
	case reflect.Struct:
		if v.CanInterface() {
			if iz, ok := reflect.TypeAssert[interface{ IsZero() bool }](v); ok {
				return iz.IsZero()
			}
		}
		for i := 0; i < v.NumField(); i++ {
			if !isZero(v.Field(i)) {
				return false
			}
		}
		return true
	}
	// Compare other types directly:
	return v.IsZero()
}

func (e *Encoder) encode(v reflect.Value, dst map[string][]string) error {
	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return errNotStruct
	}

	var errs MultiError

	plan := e.structInfo(v.Type())
	fields := plan.fields
	// When dst starts empty (fresh url.Values), single short package-allocated
	// values of distinct keys share one backing array instead of allocating a
	// 1-element slice each; the three-index slice caps entries so later
	// appends cannot overwrite a neighbor. Caller-derived strings (string
	// fields, custom encoders) and long values get their own independently
	// collectible slice so a surviving entry cannot keep a deleted neighbor's
	// allocation alive; a non-empty dst keeps the single-map-op append pattern.
	useScratch := len(dst) == 0
	// An empty dst plus a plan that writes every key once means each key is
	// new, so the read append needs — a second hash of the key — is skipped.
	fresh := useScratch && plan.freshKeys
	var scratch []string
	appendValue := func(name, s string, scratchSafe bool) {
		if !useScratch || !scratchSafe || len(s) > maxScratchValueLen {
			if fresh {
				dst[name] = []string{s}
				return
			}
			dst[name] = append(dst[name], s)
			return
		}
		if !fresh {
			if old := dst[name]; old != nil {
				dst[name] = append(old, s)
				return
			}
		}
		if scratch == nil {
			scratch = make([]string, 0, plan.scratch)
		}
		scratch = append(scratch, s)
		dst[name] = scratch[len(scratch)-1 : len(scratch) : len(scratch)]
	}
	for i := range fields {
		f := &fields[i]
		fieldValue := v.Field(f.idx)

		// Encode struct pointer types if the field is a valid pointer and a struct.
		if f.recurseStructPtr && !fieldValue.IsNil() {
			if err := e.encode(fieldValue.Elem(), dst); err != nil {
				errs = setError(errs, fieldValue.Elem().Type().String(), err)
			}
			continue
		}

		// Encode non-slice types and custom implementations immediately.
		if f.enc != nil {
			if f.omitEmpty && f.isZeroValue(fieldValue) {
				continue
			}
			appendValue(f.name, f.enc(fieldValue), f.scratchSafe)
			continue
		}

		if f.nilAsNull && fieldValue.IsNil() {
			if f.omitEmpty {
				continue
			}
			// The "null" literal is static data; always scratch-safe.
			appendValue(f.name, "null", true)
			continue
		}

		if f.isStruct {
			if err := e.encode(fieldValue, dst); err != nil {
				errs = setError(errs, fieldValue.Type().String(), err)
			}
			continue
		}

		// A non-slice field with no encoder (map, chan, array, or a non-nil
		// pointer to an unencodable type), or a slice whose element type is
		// itself unencodable and not a pointer (e.g. []Struct), cannot be
		// encoded — historically this errored unconditionally.
		if fieldValue.Kind() != reflect.Slice || (f.elemEnc == nil && !f.elemPtrNil) {
			errs = setError(errs, fieldValue.Type().String(), fmt.Errorf("schema: encoder not found for %v", fieldValue))
			continue
		}

		// Encode a slice. An empty slice has nothing to encode, so it is
		// skipped under omitempty (and otherwise emitted empty).
		n := fieldValue.Len()
		if n == 0 && f.omitEmpty {
			continue
		}

		values := make([]string, n)
		if f.elemEnc == nil {
			// Pointer elements with no encoder (elemPtrNil): nil encodes as
			// "null" (as historically), a non-nil such element is an error.
			bad := false
			for j := range n {
				if fieldValue.Index(j).IsNil() {
					values[j] = "null"
					continue
				}
				errs = setError(errs, fieldValue.Type().String(), fmt.Errorf("schema: encoder not found for %v", fieldValue))
				bad = true
				break
			}
			if bad {
				continue
			}
		} else {
			for j := range n {
				values[j] = f.elemEnc(fieldValue.Index(j))
			}
		}
		dst[f.name] = values
	}

	if len(errs) > 0 {
		return errs
	}
	return nil
}

// setError lazily allocates m and stores err under key, overwriting any
// previous entry (matching the historical encoder error semantics).
func setError(m MultiError, key string, err error) MultiError {
	if m == nil {
		m = make(MultiError)
	}
	m[key] = err
	return m
}

func (e *Encoder) hasCustomEncoder(t reflect.Type) bool {
	_, exists := e.regenc[t]
	return exists
}

// scratchSafeEncoder reports whether typeEncoder(t, reg)'s output strings are
// always allocated by this package's own formatters: string fields return the
// caller's string and custom encoders may return anything — either could be a
// small substring aliasing a large buffer, so neither may enter the scratch.
func scratchSafeEncoder(t reflect.Type, reg map[reflect.Type]encoderFunc) bool {
	if _, ok := reg[t]; ok {
		return false
	}
	switch t.Kind() {
	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return true
	case reflect.Pointer:
		return scratchSafeEncoder(t.Elem(), reg)
	default:
		return false
	}
}

func typeEncoder(t reflect.Type, reg map[reflect.Type]encoderFunc) encoderFunc {
	if f, ok := reg[t]; ok {
		return f
	}

	switch t.Kind() {
	case reflect.Bool:
		return encodeBool
	case reflect.Int8:
		return encodeInt8
	case reflect.Int16:
		return encodeInt16
	case reflect.Int32:
		return encodeInt32
	case reflect.Int, reflect.Int64:
		return encodeInt
	case reflect.Uint8:
		return encodeUint8
	case reflect.Uint16:
		return encodeUint16
	case reflect.Uint32:
		return encodeUint32
	case reflect.Uint, reflect.Uint64:
		return encodeUint
	case reflect.Float32:
		return encodeFloat32
	case reflect.Float64:
		return encodeFloat64
	case reflect.Pointer:
		f := typeEncoder(t.Elem(), reg)
		if f == nil {
			// No encoder for the element: report unsupported instead of
			// returning a closure that would panic on non-nil values.
			// Nil handling for such fields is done by encField.nilAsNull.
			return nil
		}
		return func(v reflect.Value) string {
			if v.IsNil() {
				return "null"
			}
			return f(v.Elem())
		}
	case reflect.String:
		return encodeString
	default:
		return nil
	}
}

func encodeBool(v reflect.Value) string {
	return strconv.FormatBool(v.Bool())
}

// Integers go through the width-specific utils helpers: the 8-bit ones are
// table lookups, allocating nothing, and the 32-bit ones skip a digit group
// the 64-bit formatter has to consider.

func encodeInt(v reflect.Value) string {
	return utils.FormatInt(v.Int())
}

func encodeInt8(v reflect.Value) string {
	return utils.FormatInt8(int8(v.Int()))
}

func encodeInt16(v reflect.Value) string {
	return utils.FormatInt16(int16(v.Int()))
}

func encodeInt32(v reflect.Value) string {
	return utils.FormatInt32(int32(v.Int()))
}

func encodeUint(v reflect.Value) string {
	return utils.FormatUint(v.Uint())
}

func encodeUint8(v reflect.Value) string {
	return utils.FormatUint8(uint8(v.Uint()))
}

func encodeUint16(v reflect.Value) string {
	return utils.FormatUint16(uint16(v.Uint()))
}

func encodeUint32(v reflect.Value) string {
	return utils.FormatUint32(uint32(v.Uint()))
}

func encodeFloat(v reflect.Value, bits int) string {
	return formatFloatFixed(v.Float(), bits)
}

func encodeFloat32(v reflect.Value) string {
	return encodeFloat(v, 32)
}

func encodeFloat64(v reflect.Value) string {
	return encodeFloat(v, 64)
}

func encodeString(v reflect.Value) string {
	return v.String()
}
