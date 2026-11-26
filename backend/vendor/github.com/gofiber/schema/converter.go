// Copyright 2012 The Gorilla Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package schema

import (
	"reflect"
	"strconv"

	utils "github.com/gofiber/utils/v2"
)

// Converter converts a source value to a custom type. An invalid (zero)
// reflect.Value reports that the value does not convert; the resulting
// ConversionError carries no Err, as Converter has no way to say why. Use
// ConverterE to report the reason.
type Converter func(string) reflect.Value

// ConverterE is a Converter that reports why a value does not convert. A
// non-nil error fails the conversion, whatever value is returned with it,
// and is kept as the Err of the resulting ConversionError, whose message
// then ends in its text, as for an encoding.TextUnmarshaler that fails. An
// invalid reflect.Value with a nil error fails as a Converter's does.
type ConverterE func(string) (reflect.Value, error)

// valueConverter is a converter of either kind; exactly one field is set.
// Plain converters are called directly rather than wrapped in a ConverterE,
// which would add an indirect call to every conversion.
type valueConverter struct {
	plain   Converter
	withErr ConverterE
}

// convert runs the converter. The conversion failed when the error is
// non-nil or the value is invalid; callers check both.
func (vc valueConverter) convert(s string) (reflect.Value, error) {
	if vc.withErr != nil {
		return vc.withErr(s)
	}
	return vc.plain(s), nil
}

var (
	invalidValue = reflect.Value{}
	boolType     = reflect.Bool
	float32Type  = reflect.Float32
	float64Type  = reflect.Float64
	intType      = reflect.Int
	int8Type     = reflect.Int8
	int16Type    = reflect.Int16
	int32Type    = reflect.Int32
	int64Type    = reflect.Int64
	stringType   = reflect.String
	uintType     = reflect.Uint
	uint8Type    = reflect.Uint8
	uint16Type   = reflect.Uint16
	uint32Type   = reflect.Uint32
	uint64Type   = reflect.Uint64
)

// builtinConverters is the single source of truth for type converters.
var builtinConverters = map[reflect.Kind]Converter{
	boolType:    convertBool,
	float32Type: convertFloat32,
	float64Type: convertFloat64,
	intType:     convertInt,
	int8Type:    convertInt8,
	int16Type:   convertInt16,
	int32Type:   convertInt32,
	int64Type:   convertInt64,
	stringType:  convertString,
	uintType:    convertUint,
	uint8Type:   convertUint8,
	uint16Type:  convertUint16,
	uint32Type:  convertUint32,
	uint64Type:  convertUint64,
}

// builtinConvertersArray is an array indexed by reflect.Kind for O(1) lookup.
var builtinConvertersArray [reflect.UnsafePointer + 1]Converter

func init() {
	for k, conv := range builtinConverters {
		builtinConvertersArray[k] = conv
	}
}

// getBuiltinConverter returns the converter for a kind using O(1) array lookup.
func getBuiltinConverter(k reflect.Kind) Converter {
	if k <= reflect.UnsafePointer {
		return builtinConvertersArray[k]
	}
	return nil
}

func convertBool(value string) reflect.Value {
	if value == "on" {
		return reflect.ValueOf(true)
	} else if v, err := strconv.ParseBool(value); err == nil {
		return reflect.ValueOf(v)
	}
	return invalidValue
}

func convertFloat32(value string) reflect.Value {
	if v, err := utils.ParseFloat32(value); err == nil {
		return reflect.ValueOf(v)
	}
	return invalidValue
}

func convertFloat64(value string) reflect.Value {
	if v, err := utils.ParseFloat64(value); err == nil {
		return reflect.ValueOf(v)
	}
	return invalidValue
}

// Native int/uint parsing goes through utils.ParseInt/ParseUint (which have
// a SWAR fast path) with an inline native-fit guard: `int64(int(v)) == v`
// compiles away on 64-bit and rejects values on 32-bit that a plain int()
// conversion would silently truncate. The guard is written inline at each
// call site so the ~8ns parse doesn't pay a wrapper call frame.

func convertInt(value string) reflect.Value {
	if v, err := utils.ParseInt(value); err == nil && int64(int(v)) == v {
		return reflect.ValueOf(int(v))
	}
	return invalidValue
}

func convertInt8(value string) reflect.Value {
	if v, err := utils.ParseInt8(value); err == nil {
		return reflect.ValueOf(v)
	}
	return invalidValue
}

func convertInt16(value string) reflect.Value {
	if v, err := utils.ParseInt16(value); err == nil {
		return reflect.ValueOf(v)
	}
	return invalidValue
}

func convertInt32(value string) reflect.Value {
	if v, err := utils.ParseInt32(value); err == nil {
		return reflect.ValueOf(v)
	}
	return invalidValue
}

func convertInt64(value string) reflect.Value {
	if v, err := utils.ParseInt(value); err == nil {
		return reflect.ValueOf(v)
	}
	return invalidValue
}

func convertString(value string) reflect.Value {
	return reflect.ValueOf(value)
}

func convertUint(value string) reflect.Value {
	if v, err := utils.ParseUint(value); err == nil && uint64(uint(v)) == v {
		return reflect.ValueOf(uint(v))
	}
	return invalidValue
}

func convertUint8(value string) reflect.Value {
	if v, err := utils.ParseUint8(value); err == nil {
		return reflect.ValueOf(v)
	}
	return invalidValue
}

func convertUint16(value string) reflect.Value {
	if v, err := utils.ParseUint16(value); err == nil {
		return reflect.ValueOf(v)
	}
	return invalidValue
}

func convertUint32(value string) reflect.Value {
	if v, err := utils.ParseUint32(value); err == nil {
		return reflect.ValueOf(v)
	}
	return invalidValue
}

func convertUint64(value string) reflect.Value {
	if v, err := utils.ParseUint(value); err == nil {
		return reflect.ValueOf(v)
	}
	return invalidValue
}

// parseNative* adapt the builtin parsers to the (T, bool) shape used by the
// native slice decode path, mirroring setBuiltinKind's per-kind behavior.

func parseNativeString(s string) (string, bool) { return s, true }

func parseNativeInt(s string) (int, bool) {
	n, err := utils.ParseInt(s)
	return int(n), err == nil && int64(int(n)) == n
}

func parseNativeInt64(s string) (int64, bool) {
	n, err := utils.ParseInt(s)
	return n, err == nil
}

func parseNativeUint(s string) (uint, bool) {
	n, err := utils.ParseUint(s)
	return uint(n), err == nil && uint64(uint(n)) == n
}

func parseNativeUint64(s string) (uint64, bool) {
	n, err := utils.ParseUint(s)
	return n, err == nil
}

func parseNativeFloat64(s string) (float64, bool) {
	f, err := utils.ParseFloat64(s)
	return f, err == nil
}

func parseNativeBool(s string) (bool, bool) {
	if s == "on" {
		return true, true
	}
	b, err := strconv.ParseBool(s)
	return b, err == nil
}

// setBuiltinKind parses val and assigns it directly into v for builtin
// convertible kinds, avoiding the reflect.Value boxing of the Converter API.
// handled reports whether the kind is builtin-convertible; ok reports whether
// val parsed successfully. v is only modified on success.
func setBuiltinKind(v reflect.Value, k reflect.Kind, val string) (handled, ok bool) {
	switch k {
	case boolType:
		if val == "on" {
			v.SetBool(true)
			return true, true
		}
		b, err := strconv.ParseBool(val)
		if err != nil {
			return true, false
		}
		v.SetBool(b)
	case stringType:
		v.SetString(val)
	case intType:
		n, err := utils.ParseInt(val)
		if err != nil || int64(int(n)) != n {
			return true, false
		}
		v.SetInt(n)
	case int8Type:
		n, err := utils.ParseInt8(val)
		if err != nil {
			return true, false
		}
		v.SetInt(int64(n))
	case int16Type:
		n, err := utils.ParseInt16(val)
		if err != nil {
			return true, false
		}
		v.SetInt(int64(n))
	case int32Type:
		n, err := utils.ParseInt32(val)
		if err != nil {
			return true, false
		}
		v.SetInt(int64(n))
	case int64Type:
		n, err := utils.ParseInt(val)
		if err != nil {
			return true, false
		}
		v.SetInt(n)
	case uintType:
		n, err := utils.ParseUint(val)
		if err != nil || uint64(uint(n)) != n {
			return true, false
		}
		v.SetUint(n)
	case uint8Type:
		n, err := utils.ParseUint8(val)
		if err != nil {
			return true, false
		}
		v.SetUint(uint64(n))
	case uint16Type:
		n, err := utils.ParseUint16(val)
		if err != nil {
			return true, false
		}
		v.SetUint(uint64(n))
	case uint32Type:
		n, err := utils.ParseUint32(val)
		if err != nil {
			return true, false
		}
		v.SetUint(uint64(n))
	case uint64Type:
		n, err := utils.ParseUint(val)
		if err != nil {
			return true, false
		}
		v.SetUint(n)
	case float32Type:
		f, err := utils.ParseFloat32(val)
		if err != nil {
			return true, false
		}
		v.SetFloat(float64(f))
	case float64Type:
		f, err := utils.ParseFloat64(val)
		if err != nil {
			return true, false
		}
		v.SetFloat(f)
	default:
		return false, false
	}
	return true, true
}
