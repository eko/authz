package utils

import (
	"errors"
	"math"
	"math/bits"
	"strconv"

	"github.com/gofiber/utils/v2/internal/unsafeconv"
	"github.com/gofiber/utils/v2/swar"
)

const (
	maxUint64Cutoff = math.MaxUint64 / 10
	maxUint64Cutlim = math.MaxUint64 % 10
)

type Signed interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}
type Unsigned interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

// ParseUint parses a decimal ASCII string or byte slice into a uint64, or returns 0 and an error.
func ParseUint[S byteSeq](s S) (uint64, error) {
	return parseUnsigned[S, uint64]("ParseUint", s, uint64(math.MaxUint64))
}

// ParseNativeUint is ParseUint for a uint.
func ParseNativeUint[S byteSeq](s S) (uint, error) {
	return parseUnsigned[S, uint]("ParseNativeUint", s, math.MaxUint)
}

// ParseInt is ParseUint for an int64, with an optional sign.
func ParseInt[S byteSeq](s S) (int64, error) {
	if len(s) > 0 && s[0] != '-' && s[0] != '+' && len(s) <= 19 {
		// At most 19 digits: two 8-digit steps plus a scalar tail cannot overflow before the range check.
		// Duplicates parseDigitsBig on purpose: calling it costs +22% on 9 digits (benchstat
		// -count=10, Go 1.25, Apple M2 Pro).
		var n uint64
		i := 0
		for ; len(s)-i >= 8; i += 8 {
			w := swar.Load8(s, i)
			if !isEightDigits(w) {
				break
			}
			n = n*100000000 + parse8Digits(w)
		}
		for ; i < len(s); i++ {
			c := s[i] - '0'
			if c > 9 {
				return 0, &strconv.NumError{Func: "ParseInt", Num: string(s), Err: strconv.ErrSyntax}
			}
			n = n*10 + uint64(c)
		}
		if n > uint64(math.MaxInt64) {
			return 0, &strconv.NumError{Func: "ParseInt", Num: string(s), Err: strconv.ErrRange}
		}
		return int64(n), nil
	}

	return parseSigned[S, int64]("ParseInt", s, math.MinInt64, math.MaxInt64)
}

// ParseNativeInt is ParseInt for an int.
func ParseNativeInt[S byteSeq](s S) (int, error) {
	return parseSigned[S, int]("ParseNativeInt", s, math.MinInt, math.MaxInt)
}

// ParseInt32 is ParseInt for an int32.
func ParseInt32[S byteSeq](s S) (int32, error) {
	return parseSigned[S, int32]("ParseInt32", s, math.MinInt32, math.MaxInt32)
}

// ParseInt16 is ParseInt for an int16.
func ParseInt16[S byteSeq](s S) (int16, error) {
	return parseSigned[S, int16]("ParseInt16", s, math.MinInt16, math.MaxInt16)
}

// ParseInt8 is ParseInt for an int8.
func ParseInt8[S byteSeq](s S) (int8, error) {
	if len(s) == 0 {
		return 0, &strconv.NumError{Func: "ParseInt8", Num: "", Err: strconv.ErrSyntax}
	}

	neg := false
	i := 0
	switch s[0] {
	case '-':
		neg = true
		i++
	case '+':
		i++
	}
	if i == len(s) {
		return 0, &strconv.NumError{Func: "ParseInt8", Num: string(s), Err: strconv.ErrSyntax}
	}

	if len(s)-i <= 3 {
		var n uint16
		for ; i < len(s); i++ {
			c := s[i] - '0'
			if c > 9 {
				return 0, &strconv.NumError{Func: "ParseInt8", Num: string(s), Err: strconv.ErrSyntax}
			}
			n = n*10 + uint16(c)
		}
		if neg {
			if n > 128 {
				return 0, &strconv.NumError{Func: "ParseInt8", Num: string(s), Err: strconv.ErrRange}
			}
			if n == 128 {
				return math.MinInt8, nil
			}
			return -int8(n), nil
		}
		if n > math.MaxInt8 {
			return 0, &strconv.NumError{Func: "ParseInt8", Num: string(s), Err: strconv.ErrRange}
		}
		return int8(n), nil
	}

	return parseSigned[S, int8]("ParseInt8", s, math.MinInt8, math.MaxInt8)
}

// ParseUint32 is ParseUint for a uint32.
func ParseUint32[S byteSeq](s S) (uint32, error) {
	return parseUnsigned[S, uint32]("ParseUint32", s, uint32(math.MaxUint32))
}

// ParseUint16 is ParseUint for a uint16.
func ParseUint16[S byteSeq](s S) (uint16, error) {
	return parseUnsigned[S, uint16]("ParseUint16", s, uint16(math.MaxUint16))
}

// ParseUint8 is ParseUint for a uint8.
func ParseUint8[S byteSeq](s S) (uint8, error) {
	if len(s) == 0 {
		return 0, &strconv.NumError{Func: "ParseUint8", Num: "", Err: strconv.ErrSyntax}
	}

	if len(s) <= 3 {
		var n uint16
		for i := range len(s) {
			c := s[i] - '0'
			if c > 9 {
				return 0, &strconv.NumError{Func: "ParseUint8", Num: string(s), Err: strconv.ErrSyntax}
			}
			n = n*10 + uint16(c)
		}
		if n > math.MaxUint8 {
			return 0, &strconv.NumError{Func: "ParseUint8", Num: string(s), Err: strconv.ErrRange}
		}
		return uint8(n), nil
	}

	return parseUnsigned[S, uint8]("ParseUint8", s, uint8(math.MaxUint8))
}

// isEightDigits is exact: a lane >= 0xFA can carry into its neighbor but fails its own high-nibble test.
func isEightDigits(w uint64) bool {
	const (
		highNibbles = 0xF0F0F0F0F0F0F0F0
		sixes       = 0x0606060606060606
		threes      = 0x3333333333333333
	)
	return (w&highNibbles)|((w+sixes)&highNibbles)>>4 == threes
}

// parse8Digits converts 8 digit bytes to a value; every lane must be a digit.
func parse8Digits(w uint64) uint64 {
	const (
		pairMask = 0x000000FF000000FF
		mul1     = 100 + (1000000 << 32)
		mul2     = 1 + (10000 << 32)
	)
	w -= digitZeros
	w = w*10 + w>>8
	return ((w&pairMask)*mul1 + (w>>16&pairMask)*mul2) >> 32
}

// parseDigits parses fewer than 8 digits. Callers send 8+ bytes to parseDigitsBig, ruling out
// overflow here and keeping this body inlinable.
func parseDigits[S byteSeq](s S, i int) (uint64, error) {
	var n uint64
	for ; i < len(s); i++ {
		c := s[i] - '0'
		if c > 9 {
			return 0, strconv.ErrSyntax
		}
		n = n*10 + uint64(c)
	}
	return n, nil
}

// parseDigitsBig parses 8+ digits, 8 per word while the value surely fits (19 digits do).
func parseDigitsBig[S byteSeq](s S, i int) (uint64, error) {
	var n uint64
	digits := 0
	for len(s)-i >= 8 && digits+8 <= 19 {
		w := swar.Load8(s, i)
		if !isEightDigits(w) {
			break
		}
		n = n*100000000 + parse8Digits(w)
		digits += 8
		i += 8
	}
	for ; i < len(s); i++ {
		c := s[i] - '0'
		if c > 9 {
			return 0, strconv.ErrSyntax
		}
		d := uint64(c)
		if digits >= 19 && (n > maxUint64Cutoff || (n == maxUint64Cutoff && d > maxUint64Cutlim)) {
			return 0, strconv.ErrRange
		}
		n = n*10 + d
		digits++
	}
	return n, nil
}

// classifyDigitsError returns strconv's error for a rejected digit run: digits above limit
// before a stray byte make a range error, not a syntax one.
func classifyDigitsError[S byteSeq](err error, s S, i int, limit uint64) error {
	if !errors.Is(err, strconv.ErrSyntax) {
		return err
	}
	var n uint64
	for ; i < len(s); i++ {
		c := s[i] - '0'
		if c > 9 {
			break
		}
		if n > (limit-uint64(c))/10 { // n*10+c > limit, without overflowing
			return strconv.ErrRange
		}
		n = n*10 + uint64(c)
	}
	return err
}

func parseSigned[S byteSeq, T Signed](fn string, s S, minRange, maxRange T) (T, error) {
	if len(s) == 0 {
		return 0, &strconv.NumError{Func: fn, Num: "", Err: strconv.ErrSyntax}
	}

	neg := false
	i := 0
	switch s[0] {
	case '-':
		neg = true
		i++
	case '+':
		i++
	}
	if i == len(s) {
		return 0, &strconv.NumError{Func: fn, Num: string(s), Err: strconv.ErrSyntax}
	}

	var n uint64
	var err error
	if len(s)-i >= 8 {
		n, err = parseDigitsBig(s, i)
	} else {
		n, err = parseDigits(s, i)
	}
	if err != nil {
		// strconv parses digits as the unsigned type of the same width.
		err = classifyDigitsError(err, s, i, uint64(int64(maxRange))<<1|1)
		return 0, &strconv.NumError{Func: fn, Num: string(s), Err: err}
	}

	if !neg {
		if n > uint64(int64(maxRange)) {
			return 0, &strconv.NumError{Func: fn, Num: string(s), Err: strconv.ErrRange}
		}
		return T(n), nil
	}

	minAbs := uint64(-int64(minRange))
	if n > minAbs {
		return 0, &strconv.NumError{Func: fn, Num: string(s), Err: strconv.ErrRange}
	}

	return T(-int64(n)), nil
}

func parseUnsigned[S byteSeq, T Unsigned](fn string, s S, maxRange T) (T, error) {
	if len(s) == 0 {
		return 0, &strconv.NumError{Func: fn, Num: "", Err: strconv.ErrSyntax}
	}

	var n uint64
	var err error
	if len(s) >= 8 {
		n, err = parseDigitsBig(s, 0)
	} else {
		n, err = parseDigits(s, 0)
	}
	if err != nil {
		err = classifyDigitsError(err, s, 0, uint64(maxRange))
		return 0, &strconv.NumError{Func: fn, Num: string(s), Err: err}
	}
	if n > uint64(maxRange) {
		return 0, &strconv.NumError{Func: fn, Num: string(s), Err: strconv.ErrRange}
	}
	return T(n), nil
}

// Exact powers of ten per type (5^22 < 2^53, 5^10 < 2^24), for Clinger's fast path.
var (
	float64pow10 = [...]float64{1e0, 1e1, 1e2, 1e3, 1e4, 1e5, 1e6, 1e7, 1e8, 1e9, 1e10, 1e11, 1e12, 1e13, 1e14, 1e15, 1e16, 1e17, 1e18, 1e19, 1e20, 1e21, 1e22}
	float32pow10 = [...]float32{1e0, 1e1, 1e2, 1e3, 1e4, 1e5, 1e6, 1e7, 1e8, 1e9, 1e10}
)

const (
	maxExactMantissa64 = 1 << 53
	maxExactMantissa32 = 1 << 24
	maxParsedExp       = 1 << 50

	// A 64-bit mantissa times 10^exp is 0 below minDecimalExp and infinite above maxDecimalExp.
	minDecimalExp = -342
	maxDecimalExp = 308

	// strconv reads up to maxRawLen bytes right (it keeps 800 digits).
	maxRawLen = 700
)

// parseExp reads the exponent in s[i:], clamped to ±maxParsedExp, which no digit count can offset.
func parseExp[S byteSeq](s S, i int) (int64, bool) {
	sign := int64(1)
	if i < len(s) {
		switch s[i] {
		case '-':
			sign = -1
			i++
		case '+':
			i++
		}
	}
	if i == len(s) {
		return 0, false
	}
	var exp int64
	for ; i < len(s); i++ {
		c := s[i] - '0'
		if c > 9 {
			return 0, false
		}
		exp = min(exp*10+int64(c), maxParsedExp)
	}
	return sign * exp, true
}

func parseFloat[S byteSeq](fn string, bitSize int, s S) (float64, error) {
	if len(s) == 0 {
		return 0, &strconv.NumError{Func: fn, Num: "", Err: strconv.ErrSyntax}
	}
	i := 0
	neg := false
	switch s[0] {
	case '-':
		neg = true
		i++
	case '+':
		i++
	}
	if i == len(s) {
		return 0, &strconv.NumError{Func: fn, Num: string(s), Err: strconv.ErrSyntax}
	}
	start := i

	// Once the mantissa reaches maxUint64Cutoff it is a prefix: later integer digits shift the
	// exponent, fractional ones are ignored.
	var mantissa uint64
	var exp10 int
	digits := 0
	truncated := false
	for i < len(s) {
		c := s[i] - '0'
		if c > 9 {
			break
		}
		if mantissa < maxUint64Cutoff {
			mantissa = mantissa*10 + uint64(c)
		} else {
			exp10++
			truncated = truncated || c != 0
		}
		digits++
		i++
	}

	if i < len(s) && s[i] == '.' {
		i++
		for i < len(s) {
			c := s[i] - '0'
			if c > 9 {
				break
			}
			if mantissa < maxUint64Cutoff {
				mantissa = mantissa*10 + uint64(c)
				exp10--
			} else {
				truncated = truncated || c != 0
			}
			digits++
			i++
		}
	}

	if digits == 0 {
		return 0, &strconv.NumError{Func: fn, Num: string(s), Err: strconv.ErrSyntax}
	}

	var exp int64
	if i < len(s) {
		var ok bool
		if s[i] == 'e' || s[i] == 'E' {
			exp, ok = parseExp(s, i+1)
		}
		if !ok {
			return 0, &strconv.NumError{Func: fn, Num: string(s), Err: strconv.ErrSyntax}
		}
	}
	exp += int64(exp10)

	if mantissa == 0 {
		if neg {
			return math.Copysign(0, -1), nil
		}
		return 0, nil
	}

	var f float64
	switch {
	case bitSize == 64 && exactFloat64(mantissa, exp):
		f = scale64(mantissa, exp)
	case exp < minDecimalExp:
		f = 0
	case exp > maxDecimalExp:
		f = math.Inf(1)
	case truncated:
		f = parseSlow(bitSize, s, start)
	default:
		var ok bool
		if f, ok = roundExact(mantissa, exp, bitSize); !ok {
			f = parseSlow(bitSize, s, start)
		}
	}
	if f > math.MaxFloat64 {
		return 0, &strconv.NumError{Func: fn, Num: string(s), Err: strconv.ErrRange}
	}
	if neg {
		f = -f
	}
	return f, nil
}

func roundExact(mantissa uint64, exp int64, bitSize int) (float64, bool) {
	if bitSize == 32 {
		if f, ok := scale32(mantissa, exp); ok {
			return f, true
		}
	}
	return exactDyadic(mantissa, exp, bitSize)
}

func parseSlow[S byteSeq](bitSize int, s S, start int) float64 {
	if len(s) <= maxRawLen {
		return parseNormalized(bitSize, unsafeconv.Bytes(s[start:]))
	}
	return parseLong(bitSize, s)
}

// A halfway point between floats has at most 767 digits, so 768 plus a sticky one settle the rounding.
const maxSlowDigits = 768

func parseLong[S byteSeq](bitSize int, s S) float64 {
	var buf [maxSlowDigits + 32]byte
	k, sticky, exp := longDigits(s, buf[:maxSlowDigits])
	text := buf[:k]
	if sticky {
		text, exp = append(text, '1'), exp-1
	}
	text = strconv.AppendInt(append(text, 'e'), exp, 10)
	return parseNormalized(bitSize, text)
}

// parseNormalized parses text strconv reads right; strconv keeps no reference to it,
// and a range error gives +Inf.
func parseNormalized(bitSize int, text []byte) float64 {
	f, err := strconv.ParseFloat(unsafeconv.UnsafeString(text), bitSize)
	if err != nil {
		return math.Inf(1)
	}
	return f
}

// longDigits copies the first len(buf) significant digits of the valid number s into buf,
// returning their count, whether a nonzero digit followed, and the exponent.
func longDigits[S byteSeq](s S, buf []byte) (int, bool, int64) {
	i := 0
	if s[0] == '-' || s[0] == '+' {
		i++
	}
	n, seen, frac := 0, 0, 0
	sticky, point := false, false
	for ; i < len(s); i++ {
		c := s[i]
		if c == '.' {
			point = true
			continue
		}
		if c-'0' > 9 {
			break
		}
		if point {
			frac++
		}
		if seen == 0 && c == '0' {
			continue
		}
		seen++
		if n < len(buf) {
			buf[n] = c
			n++
		} else if c != '0' {
			sticky = true
		}
	}
	var exp int64
	if i < len(s) {
		exp, _ = parseExp(s, i+1)
	}
	return n, sticky, exp - int64(frac) + int64(seen-n)
}

func exactFloat64(mantissa uint64, exp int64) bool {
	return mantissa <= maxExactMantissa64 && exp > -int64(len(float64pow10)) && exp < int64(len(float64pow10))
}

func scale64(mantissa uint64, exp int64) float64 {
	f := float64(mantissa)
	switch {
	case exp > 0:
		f *= float64pow10[exp]
	case exp < 0:
		f /= float64pow10[-exp]
	}
	return f
}

func scale32(mantissa uint64, exp int64) (float64, bool) {
	if mantissa <= maxExactMantissa32 && exp > -int64(len(float32pow10)) && exp < int64(len(float32pow10)) {
		f := float32(mantissa)
		switch {
		case exp > 0:
			f *= float32pow10[exp]
		case exp < 0:
			f /= float32pow10[-exp]
		}
		return float64(f), true
	}
	if exactFloat64(mantissa, exp) {
		// A second rounding is exact, except at a float32 halfway point: the decimal may lie either side.
		f := scale64(mantissa, exp)
		const halfway = 1 << 28 // the 29 mantissa bits a float32 drops, half set
		if math.Float64bits(f)&(1<<29-1) != halfway {
			return float64(float32(f)), true
		}
	}
	return 0, false
}

// exactDyadic rounds mantissa * 10^exp when it is an integer times 2^exp: float halfway points,
// integers past 2^53, zero-padded decimals.
func exactDyadic(mantissa uint64, exp int64, bitSize int) (float64, bool) {
	var u uint64
	switch {
	case exp >= 0 && exp < int64(len(pow10)):
		hi, lo := bits.Mul64(mantissa, pow10[exp]>>exp)
		if hi != 0 {
			return 0, false
		}
		u = lo
	case exp < 0 && -exp < int64(len(pow10)):
		p5 := pow10[-exp] >> -exp
		if mantissa%5 != 0 || mantissa%p5 != 0 {
			return 0, false
		}
		u = mantissa / p5
	default:
		return 0, false
	}
	f := float64(u)
	if bitSize == 32 {
		f = roundBits(u, 24) // float32(u) rounds twice on 32-bit targets
	}
	if exp >= 0 {
		return f * float64(uint64(1)<<exp), true
	}
	return f / float64(uint64(1)<<-exp), true
}

func roundBits(u uint64, n int) float64 {
	extra := bits.Len64(u) - n
	if extra <= 0 {
		return float64(u)
	}
	m := u >> extra
	rest := u & (1<<extra - 1)
	if half := uint64(1) << (extra - 1); rest > half || (rest == half && m&1 != 0) {
		m++
	}
	return float64(m) * float64(uint64(1)<<extra)
}

// ParseFloat64 is strconv.ParseFloat for decimal input only; on overflow it returns 0 and an error,
// not ±Inf. With over 800 digits before the point it rounds correctly where strconv does not.
func ParseFloat64[S byteSeq](s S) (float64, error) {
	return parseFloat[S]("ParseFloat64", 64, s)
}

// ParseFloat32 is ParseFloat64 for float32, rounded once from the decimal.
func ParseFloat32[S byteSeq](s S) (float32, error) {
	f, err := parseFloat[S]("ParseFloat32", 32, s)
	return float32(f), err
}
