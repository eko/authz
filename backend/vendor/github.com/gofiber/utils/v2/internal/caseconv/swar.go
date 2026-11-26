package caseconv

import (
	"encoding/binary"

	"github.com/gofiber/utils/v2/swar"
)

// Words move through encoding/binary, which measures ~9% faster than swar.Load8/Store8 in the
// in-place loops (benchstat -count=10, Go 1.25, Apple M2 Pro); the lane order is identical.

// WordLen is the SWAR word width in bytes, and the minimum length worth a call.
const WordLen = swar.WordLen

// ToLowerFrom returns a lower-cased copy of src. first must be FirstUpperIndex's
// own result (0 <= first < len(src); a larger one silently mis-converts) and
// len(src) >= WordLen. src is only read, so it may alias a string.
func ToLowerFrom(src []byte, first int) []byte {
	n := len(src)
	dst := make([]byte, n)

	i := first &^ (WordLen - 1)
	copy(dst, src[:i])

	for ; i <= n-4*WordLen; i += 4 * WordLen {
		sw := src[i : i+4*WordLen : i+4*WordLen]
		d := dst[i : i+4*WordLen : i+4*WordLen]
		binary.LittleEndian.PutUint64(d[0:WordLen], swar.ToLowerWord(binary.LittleEndian.Uint64(sw[0:WordLen])))
		binary.LittleEndian.PutUint64(d[WordLen:2*WordLen], swar.ToLowerWord(binary.LittleEndian.Uint64(sw[WordLen:2*WordLen])))
		binary.LittleEndian.PutUint64(d[2*WordLen:3*WordLen], swar.ToLowerWord(binary.LittleEndian.Uint64(sw[2*WordLen:3*WordLen])))
		binary.LittleEndian.PutUint64(d[3*WordLen:4*WordLen], swar.ToLowerWord(binary.LittleEndian.Uint64(sw[3*WordLen:4*WordLen])))
	}
	for ; i <= n-WordLen; i += WordLen {
		binary.LittleEndian.PutUint64(dst[i:i+WordLen:i+WordLen], swar.ToLowerWord(binary.LittleEndian.Uint64(src[i:i+WordLen:i+WordLen])))
	}
	if i == n {
		return dst
	}
	// The overlapping tail word may re-cover bytes below first; harmless only
	// because the fold leaves every byte under first unchanged.
	binary.LittleEndian.PutUint64(dst[n-WordLen:n:n], swar.ToLowerWord(binary.LittleEndian.Uint64(src[n-WordLen:n:n])))
	return dst
}

// ToUpperFrom mirrors ToLowerFrom, with first from FirstLowerIndex.
func ToUpperFrom(src []byte, first int) []byte {
	n := len(src)
	dst := make([]byte, n)

	i := first &^ (WordLen - 1)
	copy(dst, src[:i])

	for ; i <= n-4*WordLen; i += 4 * WordLen {
		sw := src[i : i+4*WordLen : i+4*WordLen]
		d := dst[i : i+4*WordLen : i+4*WordLen]
		binary.LittleEndian.PutUint64(d[0:WordLen], swar.ToUpperWord(binary.LittleEndian.Uint64(sw[0:WordLen])))
		binary.LittleEndian.PutUint64(d[WordLen:2*WordLen], swar.ToUpperWord(binary.LittleEndian.Uint64(sw[WordLen:2*WordLen])))
		binary.LittleEndian.PutUint64(d[2*WordLen:3*WordLen], swar.ToUpperWord(binary.LittleEndian.Uint64(sw[2*WordLen:3*WordLen])))
		binary.LittleEndian.PutUint64(d[3*WordLen:4*WordLen], swar.ToUpperWord(binary.LittleEndian.Uint64(sw[3*WordLen:4*WordLen])))
	}
	for ; i <= n-WordLen; i += WordLen {
		binary.LittleEndian.PutUint64(dst[i:i+WordLen:i+WordLen], swar.ToUpperWord(binary.LittleEndian.Uint64(src[i:i+WordLen:i+WordLen])))
	}
	if i == n {
		return dst
	}
	// See ToLowerFrom for why the overlapping store is harmless.
	binary.LittleEndian.PutUint64(dst[n-WordLen:n:n], swar.ToUpperWord(binary.LittleEndian.Uint64(src[n-WordLen:n:n])))
	return dst
}

// FirstUpperIndex returns b's first ASCII uppercase index, or -1; b is only read.
func FirstUpperIndex(b []byte) int {
	n := len(b)
	i := 0
	if n >= 32 {
		// First word alone: mixed-case names (Content-Type) hit at the start.
		if m := swar.MatchRangeMask(binary.LittleEndian.Uint64(b[0:8]), 'A', 'Z'); m != 0 {
			return swar.FirstLane(m)
		}
		// Four masks per branch; MatchRangeMask stays inline so its range folds.
		for i = 8; i <= n-32; i += 32 {
			w := b[i : i+32 : i+32]
			m0 := swar.MatchRangeMask(binary.LittleEndian.Uint64(w[0:8]), 'A', 'Z')
			m1 := swar.MatchRangeMask(binary.LittleEndian.Uint64(w[8:16]), 'A', 'Z')
			m2 := swar.MatchRangeMask(binary.LittleEndian.Uint64(w[16:24]), 'A', 'Z')
			m3 := swar.MatchRangeMask(binary.LittleEndian.Uint64(w[24:32]), 'A', 'Z')
			if m0|m1|m2|m3 != 0 {
				return firstSetLane32(i, m0, m1, m2, m3)
			}
		}
		if i == n {
			return -1
		}
		// Overlapping tail: MatchRangeMask is exact, so scanned lanes stay zero.
		i = n - 32
		w := b[i : i+32 : i+32]
		m0 := swar.MatchRangeMask(binary.LittleEndian.Uint64(w[0:8]), 'A', 'Z')
		m1 := swar.MatchRangeMask(binary.LittleEndian.Uint64(w[8:16]), 'A', 'Z')
		m2 := swar.MatchRangeMask(binary.LittleEndian.Uint64(w[16:24]), 'A', 'Z')
		m3 := swar.MatchRangeMask(binary.LittleEndian.Uint64(w[24:32]), 'A', 'Z')
		if m0|m1|m2|m3 != 0 {
			return firstSetLane32(i, m0, m1, m2, m3)
		}
		return -1
	}
	for ; i <= n-8; i += 8 {
		if m := swar.MatchRangeMask(binary.LittleEndian.Uint64(b[i:i+8]), 'A', 'Z'); m != 0 {
			return i + swar.FirstLane(m)
		}
	}
	if i == n {
		return -1
	}
	if n >= 8 {
		if m := swar.MatchRangeMask(binary.LittleEndian.Uint64(b[n-8:n]), 'A', 'Z'); m != 0 {
			return n - 8 + swar.FirstLane(m)
		}
		return -1
	}
	for ; i < n; i++ {
		if c := b[i]; c >= 'A' && c <= 'Z' {
			return i
		}
	}
	return -1
}

// firstSetLane32 returns the index of the first marked byte in the 32-byte block
// starting at base, given its four word masks, at least one of them non-zero.
func firstSetLane32(base int, m0, m1, m2, m3 uint64) int {
	if m0 != 0 {
		return base + swar.FirstLane(m0)
	}
	if m1 != 0 {
		return base + 8 + swar.FirstLane(m1)
	}
	if m2 != 0 {
		return base + 16 + swar.FirstLane(m2)
	}
	return base + 24 + swar.FirstLane(m3)
}

// FirstLowerIndex is FirstUpperIndex for ASCII lowercase bytes.
func FirstLowerIndex(b []byte) int {
	n := len(b)
	i := 0
	if n >= 32 {
		// Structure and rationale mirror FirstUpperIndex.
		if m := swar.MatchRangeMask(binary.LittleEndian.Uint64(b[0:8]), 'a', 'z'); m != 0 {
			return swar.FirstLane(m)
		}
		for i = 8; i <= n-32; i += 32 {
			w := b[i : i+32 : i+32]
			m0 := swar.MatchRangeMask(binary.LittleEndian.Uint64(w[0:8]), 'a', 'z')
			m1 := swar.MatchRangeMask(binary.LittleEndian.Uint64(w[8:16]), 'a', 'z')
			m2 := swar.MatchRangeMask(binary.LittleEndian.Uint64(w[16:24]), 'a', 'z')
			m3 := swar.MatchRangeMask(binary.LittleEndian.Uint64(w[24:32]), 'a', 'z')
			if m0|m1|m2|m3 != 0 {
				return firstSetLane32(i, m0, m1, m2, m3)
			}
		}
		if i == n {
			return -1
		}
		i = n - 32
		w := b[i : i+32 : i+32]
		m0 := swar.MatchRangeMask(binary.LittleEndian.Uint64(w[0:8]), 'a', 'z')
		m1 := swar.MatchRangeMask(binary.LittleEndian.Uint64(w[8:16]), 'a', 'z')
		m2 := swar.MatchRangeMask(binary.LittleEndian.Uint64(w[16:24]), 'a', 'z')
		m3 := swar.MatchRangeMask(binary.LittleEndian.Uint64(w[24:32]), 'a', 'z')
		if m0|m1|m2|m3 != 0 {
			return firstSetLane32(i, m0, m1, m2, m3)
		}
		return -1
	}
	for ; i <= n-8; i += 8 {
		if m := swar.MatchRangeMask(binary.LittleEndian.Uint64(b[i:i+8]), 'a', 'z'); m != 0 {
			return i + swar.FirstLane(m)
		}
	}
	if i == n {
		return -1
	}
	if n >= 8 {
		if m := swar.MatchRangeMask(binary.LittleEndian.Uint64(b[n-8:n]), 'a', 'z'); m != 0 {
			return n - 8 + swar.FirstLane(m)
		}
		return -1
	}
	for ; i < n; i++ {
		if c := b[i]; c >= 'a' && c <= 'z' {
			return i
		}
	}
	return -1
}

// ToLowerInPlace lower-cases every ASCII uppercase byte of b in place.
func ToLowerInPlace(b []byte) {
	n := len(b)
	i := 0
	for ; i <= n-32; i += 32 {
		w := b[i : i+32 : i+32]
		binary.LittleEndian.PutUint64(w[0:8], swar.ToLowerWord(binary.LittleEndian.Uint64(w[0:8])))
		binary.LittleEndian.PutUint64(w[8:16], swar.ToLowerWord(binary.LittleEndian.Uint64(w[8:16])))
		binary.LittleEndian.PutUint64(w[16:24], swar.ToLowerWord(binary.LittleEndian.Uint64(w[16:24])))
		binary.LittleEndian.PutUint64(w[24:32], swar.ToLowerWord(binary.LittleEndian.Uint64(w[24:32])))
	}
	for ; i <= n-8; i += 8 {
		w := b[i : i+8 : i+8]
		binary.LittleEndian.PutUint64(w, swar.ToLowerWord(binary.LittleEndian.Uint64(w)))
	}
	// Finish byte-wise: an overlapping word stalls on store-to-load forwarding.
	for ; i < n; i++ {
		b[i] = ToLowerTable[b[i]]
	}
}

// ToUpperInPlace upper-cases every ASCII lowercase byte of b in place.
func ToUpperInPlace(b []byte) {
	n := len(b)
	i := 0
	for ; i <= n-32; i += 32 {
		w := b[i : i+32 : i+32]
		binary.LittleEndian.PutUint64(w[0:8], swar.ToUpperWord(binary.LittleEndian.Uint64(w[0:8])))
		binary.LittleEndian.PutUint64(w[8:16], swar.ToUpperWord(binary.LittleEndian.Uint64(w[8:16])))
		binary.LittleEndian.PutUint64(w[16:24], swar.ToUpperWord(binary.LittleEndian.Uint64(w[16:24])))
		binary.LittleEndian.PutUint64(w[24:32], swar.ToUpperWord(binary.LittleEndian.Uint64(w[24:32])))
	}
	for ; i <= n-8; i += 8 {
		w := b[i : i+8 : i+8]
		binary.LittleEndian.PutUint64(w, swar.ToUpperWord(binary.LittleEndian.Uint64(w)))
	}
	for ; i < n; i++ {
		b[i] = ToUpperTable[b[i]]
	}
}
