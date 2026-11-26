package schema

import (
	"math"
	"math/bits"
	"strconv"

	utils "github.com/gofiber/utils/v2"
)

// encFloatPrec is the number of fractional digits the encoder emits for
// float fields.
const encFloatPrec = 6

// scaleFixed is 10^encFloatPrec, the factor that turns a value into the
// integer the formatter prints.
const scaleFixed = 1_000_000

// fastFixedLimit bounds the magnitude the exact-integer path accepts: below
// it |value|*10^6 stays under 10^18, so the scaled result and a round-up's +1
// both fit a uint64.
const fastFixedLimit = 1e12

// formatFloatFixed is strconv.FormatFloat(f, 'f', encFloatPrec, bitSize).
//
// strconv has no Ryū fast path for the 'f' verb with a fixed precision (see
// the `fmt != 'f'` guard in its ftoa), so every such call runs the
// multi-precision decimal path. A float's exact value is mant*2^exp, so for
// exponents in range the printed digits are exactly
// round-half-to-even(mant*10^6 * 2^exp), which a 128-bit multiply and a shift
// compute directly; utils.AppendUint then lays down the digits. Inf, NaN and
// magnitudes outside the fast range fall back to strconv.
func formatFloatFixed(f float64, bitSize int) string {
	if bitSize == 32 {
		// strconv formats the float32 value, which is what the decomposition
		// below must see.
		f = float64(float32(f))
	}

	u := math.Float64bits(f)
	neg := u>>63 != 0
	biasedExp := int(u>>52) & 0x7FF
	mant := u & (1<<52 - 1)
	if biasedExp == 0x7FF {
		// Inf or NaN.
		return strconv.FormatFloat(f, 'f', encFloatPrec, bitSize)
	}
	if biasedExp == 0 {
		biasedExp = 1 // subnormal: no implicit leading bit
	} else {
		mant |= 1 << 52
	}
	// f == mant * 2^exp exactly.
	exp := biasedExp - 1023 - 52
	// exp >= 0 means |f| >= 2^52, where f*10^6 overflows a uint64; the
	// magnitude check bounds the scaled value for everything else.
	if exp >= 0 || f > fastFixedLimit || f < -fastFixedLimit {
		return strconv.FormatFloat(f, 'f', encFloatPrec, bitSize)
	}

	// n = round-half-to-even(mant*10^6 / 2^k), computed on the exact
	// 128-bit product so the tie test sees the value's true remainder.
	k := uint(-exp)
	hi, lo := bits.Mul64(mant, scaleFixed)
	var n uint64
	switch {
	case k >= 74:
		// mant*10^6 < 2^73 <= 2^(k-1), i.e. below half a unit: rounds to 0.
		n = 0
	case k < 64:
		n = lo>>k | hi<<(64-k)
		half := uint64(1) << (k - 1)
		if r := lo & (half<<1 - 1); r > half || (r == half && n&1 != 0) {
			n++
		}
	case k == 64:
		n = hi
		if lo > 1<<63 || (lo == 1<<63 && n&1 != 0) {
			n++
		}
	default:
		s := k - 64
		n = hi >> s
		half := uint64(1) << (s - 1)
		// The discarded bits are hi's low s bits followed by all of lo, so
		// a tie needs the half bit set and every lower bit clear.
		if r := hi & (half<<1 - 1); r > half || (r == half && (lo != 0 || n&1 != 0)) {
			n++
		}
	}

	// Sign, integer digits (at most 13), '.', then six fractional digits.
	var buf [24]byte
	dst := buf[:0]
	if neg {
		dst = append(dst, '-')
	}
	dst = utils.AppendUint(dst, n/scaleFixed)
	dst = append(dst, '.')
	// Formatting 10^6+frac and dropping its leading '1' zero-pads the
	// fraction without a per-digit loop.
	var frac [8]byte
	dst = append(dst, utils.AppendUint(frac[:0], scaleFixed+n%scaleFixed)[1:]...)
	return string(dst)
}
