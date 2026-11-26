// ⚡️ Fiber is an Express inspired web framework written in Go with ☕️
// 🤖 Github Repository: https://github.com/gofiber/fiber
// 📌 API Documentation: https://docs.gofiber.io

package utils

import (
	"fmt"
	"math/bits"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/utils/v2/internal/unsafeconv"
)

// UnsafeString returns a string sharing the memory of b, without allocation; changing b changes it.
func UnsafeString(b []byte) string {
	return unsafeconv.UnsafeString(b)
}

// UnsafeBytes returns a byte slice sharing the memory of s, without allocation, to be read only: a
// write faults for a literal and alters any other string. CopyBytes gives one to write.
func UnsafeBytes(s string) []byte {
	return unsafeconv.UnsafeBytes(s)
}

// CopyString returns a copy of s, unaffected by later changes to a buffer that s views.
func CopyString(s string) string {
	// #nosec G103
	return string(UnsafeBytes(s))
}

// #nosec G103
// CopyBytes returns a copy of b with capacity equal to its length;
// a nil or empty b gives a non-nil empty slice.
func CopyBytes(b []byte) []byte {
	tmp := make([]byte, len(b))
	copy(tmp, b)
	return tmp
}

// ByteSize returns a human-readable byte string such as 10MB or 12.5KB, in powers of 1024 and
// rounded half up to one decimal place: 1048575 bytes is 1MB, not 1024KB.
func ByteSize(bytes uint64) string {
	if bytes == 0 {
		return "0B"
	}
	return formatByteSize(bytes)
}

// formatByteSize is ByteSize for a size above zero, apart so that the zero case inlines.
func formatByteSize(bytes uint64) string {
	// Written right to left on the stack, so the string is allocated at its own length: the
	// longest, 1023.9KB, is 8 bytes.
	var buf [8]byte
	i := len(buf) - 1
	buf[i] = 'B'
	whole := bytes
	exp := (bits.Len64(bytes) - 1) / 10 // 0 for B, 1 for KB, ... 6 for EB
	if exp > 0 {
		div := uint64(1) << (10 * exp)
		// Tenths of the unit, rounded half up. No uint64 overflows:
		// (bytes%div)*10 + div/2 stays below 1.21e19.
		scaled := (bytes/div)*10 + ((bytes%div)*10+div/2)/div
		if scaled >= 1024*10 {
			// 1023.95 or more of a unit rounds to 1024.0, which is 1.0 of the next; EB tops out at 16.
			scaled = 10
			exp++
		}
		i--
		buf[i] = "KMGTPE"[exp-1]
		whole = scaled / 10
		if fractional := scaled % 10; fractional > 0 {
			i -= 2
			buf[i], buf[i+1] = '.', byte(fractional)+'0'
		}
	}
	for ; whole >= 10; whole /= 10 { // at most 1023
		i--
		buf[i] = byte(whole%10) + '0'
	}
	i--
	buf[i] = byte(whole) + '0'
	return string(buf[i:])
}

// ToString converts arg to a string. A time.Time, alone or in a slice or array, uses timeFormat[0],
// else "2006-01-02 15:04:05"; a *time.Time is formatted by its String method instead. An Error or
// String method that panics on a nil pointer gives "<nil>", as in fmt; any other panic propagates,
// where fmt prints %!v(PANIC=…).
func ToString(arg any, timeFormat ...string) string {
	switch v := arg.(type) {
	case int:
		return FormatInt(int64(v))
	case int8:
		return FormatInt8(v)
	case int16:
		return FormatInt16(v)
	case int32:
		return FormatInt32(v)
	case int64:
		return FormatInt(v)
	case uint:
		return FormatUint(uint64(v))
	case uint8:
		return FormatUint8(v)
	case uint16:
		return FormatUint16(v)
	case uint32:
		return FormatUint32(v)
	case uint64:
		return FormatUint(v)
	case string:
		return v
	case []byte:
		return string(v)
	case bool:
		return strconv.FormatBool(v)
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case time.Time:
		if len(timeFormat) > 0 {
			return v.Format(timeFormat[0])
		}
		return v.Format("2006-01-02 15:04:05")
	case reflect.Value:
		if !v.IsValid() || !v.CanInterface() {
			return ""
		}
		return ToString(v.Interface(), timeFormat...)
	// error before fmt.Stringer to match the fmt package's precedence
	case error:
		if mayBeNilPointer(arg) {
			return nilReceiverString(arg, v.Error)
		}
		return v.Error()
	case fmt.Stringer:
		if mayBeNilPointer(arg) {
			return nilReceiverString(arg, v.String)
		}
		return v.String()
	case *string:
		if v != nil {
			return *v
		}
		return ""
	case *int:
		if v != nil {
			return FormatInt(int64(*v))
		}
		return "0"
	case *int64:
		if v != nil {
			return FormatInt(*v)
		}
		return "0"
	case *uint64:
		if v != nil {
			return FormatUint(*v)
		}
		return "0"
	case *float64:
		if v != nil {
			return strconv.FormatFloat(*v, 'f', -1, 64)
		}
		return "0"
	case *bool:
		if v != nil {
			return strconv.FormatBool(*v)
		}
		return "false"
	case []string:
		if len(v) == 0 {
			return "[]"
		}
		size := len(v) + 1
		for _, s := range v {
			size += len(s)
		}
		buf := make([]byte, 0, size)
		buf = append(buf, '[')
		for i, s := range v {
			if i > 0 {
				buf = append(buf, ' ')
			}
			buf = append(buf, s...)
		}
		buf = append(buf, ']')
		return UnsafeString(buf)
	case []int:
		if len(v) == 0 {
			return "[]"
		}
		size := len(v) + 1
		for _, n := range v {
			size += intDigits(int64(n))
		}
		buf := make([]byte, 0, size)
		buf = append(buf, '[')
		for i, n := range v {
			if i > 0 {
				buf = append(buf, ' ')
			}
			buf = AppendInt(buf, int64(n))
		}
		buf = append(buf, ']')
		return UnsafeString(buf)
	default:
		return toStringReflect(arg, timeFormat)
	}
}

// toStringReflect is ToString's reflection fallback, apart so ToString's frame stays small: 336 bytes
// of locals against 472, and /int64 and /string run 15% and 13% faster (benchstat -count=10, Xeon 2.3 GHz).
func toStringReflect(arg any, timeFormat []string) string {
	rv := reflect.ValueOf(arg)
	if !rv.IsValid() {
		return ""
	}
	kind := rv.Kind()
	if kind == reflect.Pointer && !rv.IsNil() {
		if endlessPointerType(rv.Type()) {
			return fmt.Sprint(arg)
		}
		return ToString(rv.Elem().Interface(), timeFormat...)
	} else if kind == reflect.Slice || kind == reflect.Array {
		n := rv.Len()
		if n == 0 {
			return "[]"
		}
		var buf strings.Builder
		buf.Grow(n * 8)
		buf.WriteByte('[')
		for i := range n {
			if i > 0 {
				buf.WriteByte(' ')
			}
			buf.WriteString(ToString(rv.Index(i).Interface(), timeFormat...))
		}
		buf.WriteByte(']')
		return buf.String()
	}

	return fmt.Sprint(arg)
}

// endlessPointerType reports whether pointer type t leads only to pointers, as in type T *T,
// which would make ToString overflow the stack.
func endlessPointerType(t reflect.Type) bool {
	slow := t
	for {
		for range 2 {
			if t.Kind() != reflect.Pointer {
				return false
			}
			t = t.Elem()
		}
		if slow = slow.Elem(); slow == t {
			return true
		}
	}
}

// nilPointerReflect is mayBeNilPointer without the gc shortcut.
func nilPointerReflect(v any) bool {
	rv := reflect.ValueOf(v)
	return rv.Kind() == reflect.Pointer && rv.IsNil()
}

// nilReceiverString calls method, the Error or String of arg, whose data word is nil: if it
// panics, a nil pointer gives "<nil>", as in fmt, and any other value its panic back.
func nilReceiverString(arg any, method func() string) (s string) { //nolint:nonamedreturns // the deferred recover sets the result
	defer func() {
		if r := recover(); r != nil {
			if v := reflect.ValueOf(arg); v.Kind() != reflect.Pointer || !v.IsNil() {
				panic(r)
			}
			s = "<nil>"
		}
	}()
	return method()
}
