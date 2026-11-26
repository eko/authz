// ⚡️ Fiber is an Express inspired web framework written in Go with ☕️
// 🤖 Github Repository: https://github.com/gofiber/fiber
// 📌 API Documentation: https://docs.gofiber.io

package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"math"
	"math/bits"
	"net"
	"os"
	"reflect"
	"runtime"
	"slices"
	"strconv"

	"github.com/google/uuid"
)

const (
	defaultSecureTokenLength  = 32
	maxFastTokenEncodedLength = 43 // base64.RawURLEncoding.EncodedLen(32)
)

var whitespaceTable = [256]bool{
	'\t': true,
	'\n': true,
	'\v': true,
	'\f': true,
	'\r': true,
	' ':  true,
}

// UUIDv4 returns a random (version 4) UUID from crypto/rand; it panics if the random source fails.
func UUIDv4() string {
	token, err := uuid.NewRandom()
	if err != nil {
		panic(fmt.Errorf("utils: failed to generate secure UUID: %w", err))
	}
	return token.String()
}

// UUID is an alias for UUIDv4, kept for backward compatibility.
func UUID() string {
	return UUIDv4()
}

// GenerateSecureToken returns a cryptographically secure random token of length bytes, encoded with
// base64.RawURLEncoding, so it is longer; length <= 0 means 32 bytes. The length is not capped, so
// bound one that comes from outside. If the random source fails, the program ends.
func GenerateSecureToken(length int) string {
	if length <= 0 {
		length = defaultSecureTokenLength
	}

	if length == defaultSecureTokenLength {
		var randomBuf [defaultSecureTokenLength]byte
		src := randomBuf[:]
		rand.Read(src) //nolint:errcheck // never fails: crypto/rand ends the program instead (Go 1.24+)

		var encoded [maxFastTokenEncodedLength]byte
		encodedLen := base64.RawURLEncoding.EncodedLen(length)
		base64.RawURLEncoding.Encode(encoded[:encodedLen], src)
		return string(encoded[:encodedLen])
	}

	bytes := make([]byte, length)
	rand.Read(bytes) //nolint:errcheck // never fails: crypto/rand ends the program instead (Go 1.24+)
	return base64.RawURLEncoding.EncodeToString(bytes)
}

// SecureToken generates a secure token with 32 bytes of entropy; see GenerateSecureToken.
func SecureToken() string {
	return GenerateSecureToken(defaultSecureTokenLength)
}

// FunctionName returns function name
func FunctionName(fn any) string {
	if fn == nil {
		return ""
	}
	v := reflect.ValueOf(fn)
	if v.Kind() == reflect.Func {
		if v.IsNil() {
			return ""
		}
		return runtime.FuncForPC(v.Pointer()).Name() // "" for an address that belongs to no function
	}
	return v.Type().String()
}

// GetArgument reports whether arg is among the program arguments; os.Args[0] does not count.
func GetArgument(arg string) bool {
	if len(os.Args) < 2 {
		return false
	}
	return slices.Contains(os.Args[1:], arg)
}

// IncrementIPRange adds one to ip in place, wrapping the last address to all zeros. A 16-byte IPv4
// address from net.ParseIP carries out of 255.255.255.255 into the ::ffff: prefix; call To4 first.
func IncrementIPRange(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

// ConvertToBytes returns the size in bytes of a human-readable string, ex. 42kb, 42M, 1.5e3K.
// The first letter after the number picks the unit (k, m, g, t, p in either case, and e before a
// b or i, as in EB, since a bare E is an exponent; any other letter means bytes); the rest is
// ignored. Units are powers of 1000 (42k = 42000), or of 1024 with an 'i' right after the letter,
// as in Ki and KiB but not in kilobytes (42Ki = 43008), so the KB that ByteSize writes for 1024
// reads back as 1000. The result is rounded down and capped at math.MaxInt; it is 0 for a
// missing, negative or malformed number.
func ConvertToBytes(humanReadableString string) int {
	strLen := len(humanReadableString)
	if strLen == 0 {
		return 0
	}

	var sizeFast uint64
	maxInt := uint64(math.MaxInt)
	i := 0
	for ; i < strLen; i++ {
		c := humanReadableString[i]
		if c < '0' || c > '9' {
			break
		}
		d := uint64(c - '0')
		if sizeFast > maxInt/10 || (sizeFast == maxInt/10 && d > maxInt%10) {
			sizeFast = maxInt
		} else if sizeFast < maxInt {
			sizeFast = sizeFast*10 + d
		}
	}
	if i > 0 {
		if i == strLen {
			return int(sizeFast)
		}
		if i+1 == strLen {
			last := humanReadableString[i]
			if last == 'b' || last == 'B' {
				return int(sizeFast)
			}
		}
	}

	lastNumberPos := -1
	unitPrefixPos := 0
	for j := strLen - 1; j >= 0; j-- {
		c := humanReadableString[j]
		if c >= '0' && c <= '9' {
			lastNumberPos = j
			break
		}
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') {
			unitPrefixPos = j
		}
	}

	if lastNumberPos < 0 {
		return 0
	}

	numPart := humanReadableString[:lastNumberPos+1]

	mult := uint64(1)
	if unitPrefixPos > 0 {
		binary := false
		if next := unitPrefixPos + 1; next < strLen && (humanReadableString[next] == 'i' || humanReadableString[next] == 'I') {
			binary = next+1 == strLen || !continuesUnitWord(humanReadableString[next+1])
		}
		var decMul, binMul uint64
		switch humanReadableString[unitPrefixPos] {
		case 'k', 'K':
			decMul, binMul = 1e3, 1<<10
		case 'm', 'M':
			decMul, binMul = 1e6, 1<<20
		case 'g', 'G':
			decMul, binMul = 1e9, 1<<30
		case 't', 'T':
			decMul, binMul = 1e12, 1<<40
		case 'p', 'P':
			decMul, binMul = 1e15, 1<<50
		case 'e', 'E':
			decMul, binMul = 1, 1
			if next := unitPrefixPos + 1; next < strLen && (humanReadableString[next]|0x20 == 'b' || humanReadableString[next]|0x20 == 'i') {
				decMul, binMul = 1e18, 1<<60
			}
		default:
			decMul, binMul = 1, 1
		}
		mult = decMul
		if binary {
			mult = binMul
		}
	}

	if lastNumberPos == i-1 {
		hi, lo := bits.Mul64(sizeFast, mult)
		if hi != 0 || lo > maxInt {
			return math.MaxInt
		}
		return int(lo)
	}

	// Integers, not floats: 1.001 * 1000 is 1000.9999999999999.
	return scaleDecimal(numPart, mult)
}

// continuesUnitWord reports whether c makes the "i" of "Ki" part of a word like "kilobytes".
func continuesUnitWord(c byte) bool {
	l := c | 0x20
	return l >= 'a' && l <= 'z' && l != 'b'
}

// ParseByteSize is ConvertToBytes without guessing: a number, blanks and at most one unit (b; k to p,
// with i and b: KiB; e, with i or b: EB), else 0 and a *strconv.NumError, ErrSyntax or, for
// math.MaxInt or more, ErrRange. The only sign is '+'.
func ParseByteSize(s string) (int, error) {
	if !validSize(s) {
		return 0, &strconv.NumError{Func: "ParseByteSize", Num: s, Err: strconv.ErrSyntax}
	}
	n := ConvertToBytes(s)
	if n == math.MaxInt {
		return 0, &strconv.NumError{Func: "ParseByteSize", Num: s, Err: strconv.ErrRange}
	}
	return n, nil
}

func validSize(s string) bool {
	n := len(s)
	i := 0
	if i < n && s[i] == '+' {
		i++
	}
	j := skipDigits(s, i)
	digits := j - i
	i = j
	if i < n && s[i] == '.' {
		j = skipDigits(s, i+1)
		digits += j - i - 1
		i = j
	}
	if digits == 0 {
		return false
	}
	// An e without digits after it starts the unit instead.
	if i < n && s[i]|0x20 == 'e' {
		j = i + 1
		if j < n && (s[j] == '+' || s[j] == '-') {
			j++
		}
		if k := skipDigits(s, j); k > j {
			i = k
		}
	}
	for i < n && (s[i] == ' ' || s[i] == '\t') {
		i++
	}
	return i == n || validSizeUnit(s[i:])
}

func skipDigits(s string, i int) int {
	for i < len(s) && s[i]-'0' <= 9 {
		i++
	}
	return i
}

// validSizeUnit reports whether u, not empty, is a unit ParseByteSize accepts.
func validSizeUnit(u string) bool {
	c, rest := u[0]|0x20, u[1:]
	if c == 'b' {
		return rest == ""
	}
	if c != 'k' && c != 'm' && c != 'g' && c != 't' && c != 'p' && c != 'e' {
		return false
	}
	hasI := rest != "" && rest[0]|0x20 == 'i'
	if hasI {
		rest = rest[1:]
	}
	hasB := rest != "" && rest[0]|0x20 == 'b'
	if hasB {
		rest = rest[1:]
	}
	return rest == "" && (c != 'e' || hasI || hasB)
}

func scaleDecimal(s string, mult uint64) int {
	var mantissa uint64
	var exp int64
	end := len(s)
	places, dropped := -1, 0 // places: digits after the point, -1 before it
	droppedFrac := 0
	inexact := false
scan:
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		case c-'0' <= 9:
			switch {
			case mantissa < maxUint64Cutoff:
				mantissa = mantissa*10 + uint64(c-'0')
				if places >= 0 {
					places++
				}
			case places < 0:
				dropped++
				inexact = inexact || c != '0'
			default:
				droppedFrac++
				inexact = inexact || c != '0'
			}
		case c == '.' && places < 0:
			places = 0
		case c == 'e' || c == 'E':
			var ok bool
			if exp, ok = parseExp(s, i+1); !ok {
				return 0
			}
			end = i
			break scan
		case c == '+' && i == 0:
		default:
			return 0
		}
	}
	if mantissa == 0 {
		return 0
	}

	// The mantissa is a prefix of the digits, so the number lies in [mantissa, mantissa+1) * 10^shift;
	// if both ends scale alike, so does it.
	shift := exp + int64(dropped) - int64(max(places, 0))
	n := scaleMantissa(mantissa, shift, mult)
	if inexact && n != scaleMantissa(mantissa+1, shift, mult) {
		return scaleDigits(s, end, int64(droppedFrac)-shift+int64(dropped), mult)
	}
	return n
}

func scaleMantissa(mantissa uint64, shift int64, mult uint64) int {
	if shift > 0 {
		if shift >= int64(len(pow10)) {
			return math.MaxInt
		}
		hi, lo := bits.Mul64(mantissa, pow10[shift])
		if hi != 0 {
			return math.MaxInt
		}
		mantissa = lo
	}
	hi, lo := bits.Mul64(mantissa, mult)
	if shift >= 0 {
		if hi != 0 || lo > uint64(math.MaxInt) {
			return math.MaxInt
		}
		return int(lo)
	}

	var quotient uint64
	if down := -shift; down < int64(len(pow10)) {
		divisor := pow10[down]
		switch {
		case hi == 0:
			quotient = lo / divisor
		case hi >= divisor:
			return math.MaxInt
		default:
			quotient, _ = bits.Div64(hi, lo, divisor)
		}
	} else if down < int64(2*len(pow10)-1) {
		// The product is below 2^114, so dividing by 10^19 cannot overflow.
		quotient, _ = bits.Div64(hi, lo, pow10[len(pow10)-1])
		quotient /= pow10[down-int64(len(pow10)-1)]
	} // else the product is below 10^down
	if quotient > uint64(math.MaxInt) {
		return math.MaxInt
	}
	return int(quotient)
}

// scaleDigits multiplies the digits of s[:end] by mult and keeps the product's digits from the
// down-th on: the integer part, as the point is down places from the end.
func scaleDigits(s string, end int, down int64, mult uint64) int {
	if down <= 0 {
		return math.MaxInt
	}
	var result, carry uint64
	pos := int64(0)
	for i := end - 1; i >= 0 || carry != 0; i-- {
		var d uint64
		if i >= 0 {
			c := s[i] - '0'
			if c > 9 {
				continue
			}
			d = uint64(c)
		}
		t := d*mult + carry
		carry = t / 10
		if digit := t % 10; digit != 0 && pos >= down {
			if pos-down >= int64(len(pow10)-1) { // 10^19 > math.MaxInt
				return math.MaxInt
			}
			result += digit * pow10[pos-down] // below 2^64 even after the last sum
			if result > uint64(math.MaxInt) {
				return math.MaxInt
			}
		}
		pos++
	}
	return int(result)
}
