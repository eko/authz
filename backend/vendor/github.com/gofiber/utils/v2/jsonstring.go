package utils

import (
	"unicode/utf8"

	"github.com/gofiber/utils/v2/internal/unsafeconv"
	"github.com/gofiber/utils/v2/swar"
)

const (
	ltLanes  = uint64('<') * swar.Ones
	gtLanes  = uint64('>') * swar.Ones
	ampLanes = uint64('&') * swar.Ones
)

const jsonHexDigits = "0123456789abcdef"

// DEL passes, as in encoding/json; bytes >= 0x80 are false and take the UTF-8 path.
var jsonSafeTable = buildJSONSafeTable()

func buildJSONSafeTable() [256]bool {
	var t [256]bool
	for c := byte(0x20); c < utf8.RuneSelf; c++ {
		t[c] = c != '"' && c != '\\' && c != '<' && c != '>' && c != '&'
	}
	return t
}

// AppendJSONString appends s to dst as a double-quoted JSON string. '<', '>' and '&' become
// \u00XX, U+2028 and U+2029 are escaped, and each invalid UTF-8 byte becomes the six characters
// \ufffd. That is byte-identical to json.Marshal of s through Go 1.26; from Go 1.27 Marshal writes
// a raw U+FFFD for an invalid byte instead, and both decode to the same string. dst must not alias s.
func AppendJSONString[S byteSeq](dst []byte, s S) []byte {
	bs := unsafeconv.Bytes(s)
	n := len(bs)
	dst = append(dst, '"')
	start, i := 0, 0
	for i < n {
		if i <= n-8 {
			m := jsonUnsafeMask(swar.Load8(bs, i))
			if m == 0 {
				i += 8
				continue
			}
			i += swar.FirstLane(m)
		} else if c := bs[i]; jsonSafeTable[c] {
			i++
			continue
		}
		c := bs[i]
		if c < utf8.RuneSelf {
			dst = append(dst, bs[start:i]...)
			dst = append(dst, '\\')
			switch c {
			case '"', '\\':
				dst = append(dst, c)
			case '\b':
				dst = append(dst, 'b')
			case '\f':
				dst = append(dst, 'f')
			case '\n':
				dst = append(dst, 'n')
			case '\r':
				dst = append(dst, 'r')
			case '\t':
				dst = append(dst, 't')
			default:
				dst = append(dst, 'u', '0', '0', jsonHexDigits[c>>4], jsonHexDigits[c&0x0F])
			}
			i++
			start = i
			continue
		}
		r, size := utf8.DecodeRune(bs[i:])
		switch {
		case r == utf8.RuneError && size == 1:
			dst = append(dst, bs[start:i]...)
			dst = append(dst, '\\', 'u', 'f', 'f', 'f', 'd')
			i++
			start = i
		case r == '\u2028' || r == '\u2029':
			dst = append(dst, bs[start:i]...)
			dst = append(dst, '\\', 'u', '2', '0', '2', byte('8'+r-'\u2028'))
			i += size
			start = i
		default:
			i += size
		}
	}
	dst = append(dst, bs[start:]...)
	return append(dst, '"')
}

// jsonUnsafeMask flags lanes of w needing escaping or UTF-8 handling; exact only in and
// below the first set lane, which is all the first-match scan uses.
func jsonUnsafeMask(w uint64) uint64 {
	return swar.MatchRangeMask(w, 0x00, 0x1F) | w&swar.HighBits |
		swar.ZeroLanes(w^quoteLanes) | swar.ZeroLanes(w^backslashLanes) |
		swar.ZeroLanes(w^ltLanes) | swar.ZeroLanes(w^gtLanes) | swar.ZeroLanes(w^ampLanes)
}
