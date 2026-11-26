package utils

import (
	"github.com/gofiber/utils/v2/internal/unsafeconv"
	"github.com/gofiber/utils/v2/simd"
	"github.com/gofiber/utils/v2/swar"
)

const (
	quoteLanes     = uint64('"') * swar.Ones
	backslashLanes = uint64('\\') * swar.Ones
	tabLanes       = uint64('\t') * swar.Ones
)

// Dispatch below gates on simd.MinLen, the single source of truth for the AVX2 crossover.

// IsASCII reports whether s contains only ASCII bytes (no byte >= 0x80).
// It mirrors simd.isASCIIGeneric; retune both together.
func IsASCII[S byteSeq](s S) bool {
	n := len(s)
	if n >= simd.MinLen && simd.Accelerated() {
		return simd.IsASCII(unsafeconv.Bytes(s))
	}
	if n >= 8 && n < 16 {
		return (swar.Load8(s, 0)|swar.Load8(s, n-8))&swar.HighBits == 0
	}
	i := 0
	for ; i <= n-32; i += 32 {
		w := s[i : i+32]
		acc := swar.Load8(w, 0) | swar.Load8(w, 8) | swar.Load8(w, 16) | swar.Load8(w, 24)
		if acc&swar.HighBits != 0 {
			return false
		}
	}
	for ; i <= n-16; i += 16 {
		w := s[i : i+16]
		if (swar.Load8(w, 0)|swar.Load8(w, 8))&swar.HighBits != 0 {
			return false
		}
	}
	for ; i <= n-8; i += 8 {
		if swar.Load8(s, i)&swar.HighBits != 0 {
			return false
		}
	}
	if i == n {
		return true
	}
	if n >= 8 {
		return swar.Load8(s, n-8)&swar.HighBits == 0
	}
	for ; i < n; i++ {
		if s[i] >= 0x80 {
			return false
		}
	}
	return true
}

// IndexNonQuotable returns the index of the first byte of s that cannot appear
// verbatim inside an RFC 9110 quoted-string, or -1 if every byte is quotable.
// Those are backslash, double quote, DEL, and controls below 0x20 except HTAB;
// bytes >= 0x80 are quotable.
func IndexNonQuotable[S byteSeq](s S) int {
	n := len(s)
	i := 0
	for ; i <= n-16; i += 16 {
		// Pinned window: constant-offset loads skip bounds checks (see swar.Load8).
		w := s[i : i+16]
		m0 := nonQuotableMask(swar.Load8(w, 0))
		m1 := nonQuotableMask(swar.Load8(w, 8))
		if m0|m1 != 0 {
			if m0 != 0 {
				return i + swar.FirstLane(m0)
			}
			return i + 8 + swar.FirstLane(m1)
		}
	}
	for ; i <= n-8; i += 8 {
		if m := nonQuotableMask(swar.Load8(s, i)); m != 0 {
			return i + swar.FirstLane(m)
		}
	}
	if i == n {
		return -1
	}
	if n >= 8 {
		if m := nonQuotableMask(swar.Load8(s, n-8)); m != 0 {
			return n - 8 + swar.FirstLane(m)
		}
		return -1
	}
	for ; i < n; i++ {
		if c := s[i]; c == '\\' || c == '"' || (c < 0x20 && c != '\t') || c == 0x7f {
			return i
		}
	}
	return -1
}

// nonQuotableMask is, like swar.ZeroLanes, exact in and below its first marked
// lane, which is all IndexNonQuotable consumes.
func nonQuotableMask(w uint64) uint64 {
	return controlLanes(w)&otherLanes(w, tabLanes) | swar.ZeroLanes(w^quoteLanes) | swar.ZeroLanes(w^backslashLanes)
}
