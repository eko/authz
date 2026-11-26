package utils

import (
	"math/bits"

	"github.com/gofiber/utils/v2/swar"
)

type byteSeq interface {
	~string | ~[]byte
}

// load4 assembles s[i:i+4] little-endian; the caller guarantees i+4 <= len(s).
func load4[S byteSeq](s S, i int) uint32 {
	w := s[i : i+4]
	return uint32(w[0]) | uint32(w[1])<<8 | uint32(w[2])<<16 | uint32(w[3])<<24
}

// EqualFold tests ascii strings or bytes for equality case-insensitively
func EqualFold[S byteSeq](b, s S) bool {
	n := len(b)
	if n != len(s) {
		return false
	}

	i := 0
	for ; i <= n-8; i += 8 {
		x := swar.Load8(b, i)
		y := swar.Load8(s, i)
		if x != y && !foldEqualWords(x, y) {
			return false
		}
	}
	if i == n {
		return true
	}
	if n >= 8 {
		x := swar.Load8(b, n-8)
		y := swar.Load8(s, n-8)
		return x == y || foldEqualWords(x, y)
	}
	if n >= 4 {
		// 4..7 bytes: two overlapping 4-byte windows packed into one word.
		x := uint64(load4(b, 0)) | uint64(load4(b, n-4))<<32
		y := uint64(load4(s, 0)) | uint64(load4(s, n-4))<<32
		return x == y || foldEqualWords(x, y)
	}

	// 1..3 bytes: the first, the middle and the last one cover them all.
	x := uint64(b[0]) | uint64(b[n>>1])<<8 | uint64(b[n-1])<<16
	y := uint64(s[0]) | uint64(s[n>>1])<<8 | uint64(s[n-1])<<16
	return x == y || foldEqualWords(x, y)
}

// HashFold returns a 64-bit hash of s in which ASCII letters are folded to one
// case: strings that EqualFold reports equal hash equal. As there, only
// 'A'..'Z' and 'a'..'z' fold. It is not cryptographic and takes no seed, so do
// not use it for keys an attacker picks. Its values may change between releases.
func HashFold[S byteSeq](s S) uint64 {
	n := len(s)
	if n < 8 {
		// Keys of different lengths can pack alike ("a", "aa"), so n goes into the multiplier.
		var w uint64
		switch {
		case n >= 4:
			w = uint64(load4(s, 0)) | uint64(load4(s, n-4))<<32
		case n > 0:
			w = uint64(s[0]) | uint64(s[n>>1])<<8 | uint64(s[n-1])<<16
		}
		return hashFoldMix(swar.ToLowerWord(w)^hashFoldK0, hashFoldK1^uint64(n)<<56)
	}
	// Longer keys follow wyhash: 16 bytes per mix, the last 16 overlapping, a final mix with n.
	h := hashFoldK2
	for i := 0; n-i > 16; i += 16 {
		w := s[i : i+16]
		h = hashFoldMix(swar.ToLowerWord(swar.Load8(w, 0))^hashFoldK1, swar.ToLowerWord(swar.Load8(w, 8))^h)
	}
	hi, lo := bits.Mul64(
		swar.ToLowerWord(swar.Load8(s, max(n-16, 0)))^hashFoldK1,
		swar.ToLowerWord(swar.Load8(s, n-8))^h,
	)
	return hashFoldMix(lo^hashFoldK0^uint64(n), hi^hashFoldK1)
}

const (
	hashFoldK0 uint64 = 0xa0761d6478bd642f
	hashFoldK1 uint64 = 0xe7037ed1a0b428db
	hashFoldK2 uint64 = 0x8ebc6af09c88c6e3
)

func hashFoldMix(a, b uint64) uint64 {
	hi, lo := bits.Mul64(a, b)
	return hi ^ lo
}

// TrimLeft removes all leading occurrences of the byte cutset from s. Unlike strings.TrimLeft,
// cutset is a single byte, not a set. The result is a subslice of s, nil only if s is, where
// bytes.TrimLeft returns nil when nothing is left.
func TrimLeft[S byteSeq](s S, cutset byte) S {
	lenStr, start := len(s), 0
	for start < lenStr && s[start] == cutset {
		start++
	}
	return s[start:]
}

// Trim removes all leading and trailing occurrences of the byte cutset from s. Unlike
// strings.Trim, cutset is a single byte, not a set. The result is a subslice of s, nil only if
// s is, where bytes.Trim returns nil when nothing is left.
func Trim[S byteSeq](s S, cutset byte) S {
	i, j := 0, len(s)-1
	for ; i <= j; i++ {
		if s[i] != cutset {
			break
		}
	}
	for ; i < j; j-- {
		if s[j] != cutset {
			break
		}
	}

	return s[i : j+1]
}

// TrimRight removes all trailing occurrences of the byte cutset from s. Unlike strings.TrimRight,
// cutset is a single byte, not a set. The result is a subslice of s.
func TrimRight[S byteSeq](s S, cutset byte) S {
	lenStr := len(s)
	for lenStr > 0 && s[lenStr-1] == cutset {
		lenStr--
	}
	return s[:lenStr]
}

// TrimSpace removes leading and trailing ASCII whitespace (space, tab, LF, CR,
// VT, FF) from s; unlike the stdlib, it leaves Unicode spaces alone. The result
// is a subslice of s, nil only if s is.
func TrimSpace[S byteSeq](s S) S {
	n := len(s)
	if n == 0 {
		return s
	}

	i, j := 0, n-1
	if !whitespaceTable[s[i]] && !whitespaceTable[s[j]] {
		return s
	}

	for ; i <= j && whitespaceTable[s[i]]; i++ { //nolint:revive // we want to check for multiple whitespace chars
	}

	for ; i < j && whitespaceTable[s[j]]; j-- { //nolint:revive // we want to check for multiple whitespace chars
	}

	return s[i : j+1]
}
