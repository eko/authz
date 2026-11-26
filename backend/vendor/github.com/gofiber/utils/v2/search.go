package utils

import (
	"github.com/gofiber/utils/v2/internal/caseconv"
	"github.com/gofiber/utils/v2/internal/unsafeconv"
	"github.com/gofiber/utils/v2/simd"
	"github.com/gofiber/utils/v2/swar"
)

// IndexAny2/IndexAny3 finish with an overlapping word at n-8; already-scanned
// lanes never match, so the first set lane lies in the new bytes.

// IndexAny2 returns the index of the first occurrence in s of a or b, or -1.
func IndexAny2[S byteSeq](s S, a, b byte) int {
	n := len(s)
	if n >= simd.MinLen && simd.Accelerated() {
		return simd.Memchr2(unsafeconv.Bytes(s), a, b)
	}
	if n >= 8 {
		bcA := swar.Broadcast(a)
		bcB := swar.Broadcast(b)
		i := 0
		// Two words per branch: if only the second matched, its first set lane is exact.
		for ; i <= n-16; i += 16 {
			w0 := swar.Load8(s, i)
			w1 := swar.Load8(s, i+8)
			m0 := swar.ZeroLanes(w0^bcA) | swar.ZeroLanes(w0^bcB)
			m1 := swar.ZeroLanes(w1^bcA) | swar.ZeroLanes(w1^bcB)
			if m0|m1 != 0 {
				if m0 != 0 {
					return i + swar.FirstLane(m0)
				}
				return i + 8 + swar.FirstLane(m1)
			}
		}
		for ; i <= n-8; i += 8 {
			w := swar.Load8(s, i)
			if m := swar.ZeroLanes(w^bcA) | swar.ZeroLanes(w^bcB); m != 0 {
				return i + swar.FirstLane(m)
			}
		}
		if i == n {
			return -1
		}
		w := swar.Load8(s, n-8)
		if m := swar.ZeroLanes(w^bcA) | swar.ZeroLanes(w^bcB); m != 0 {
			return n - 8 + swar.FirstLane(m)
		}
		return -1
	}
	for i := range n {
		if s[i] == a || s[i] == b {
			return i
		}
	}
	return -1
}

// IndexAny3 returns the index of the first occurrence in s of a, b, or c, or -1.
func IndexAny3[S byteSeq](s S, a, b, c byte) int {
	n := len(s)
	if n >= simd.MinLen && simd.Accelerated() {
		return simd.Memchr3(unsafeconv.Bytes(s), a, b, c)
	}
	if n >= 8 {
		bcA := swar.Broadcast(a)
		bcB := swar.Broadcast(b)
		bcC := swar.Broadcast(c)
		i := 0
		for ; i <= n-16; i += 16 {
			w0 := swar.Load8(s, i)
			w1 := swar.Load8(s, i+8)
			m0 := swar.ZeroLanes(w0^bcA) | swar.ZeroLanes(w0^bcB) | swar.ZeroLanes(w0^bcC)
			m1 := swar.ZeroLanes(w1^bcA) | swar.ZeroLanes(w1^bcB) | swar.ZeroLanes(w1^bcC)
			if m0|m1 != 0 {
				if m0 != 0 {
					return i + swar.FirstLane(m0)
				}
				return i + 8 + swar.FirstLane(m1)
			}
		}
		for ; i <= n-8; i += 8 {
			w := swar.Load8(s, i)
			if m := swar.ZeroLanes(w^bcA) | swar.ZeroLanes(w^bcB) | swar.ZeroLanes(w^bcC); m != 0 {
				return i + swar.FirstLane(m)
			}
		}
		if i == n {
			return -1
		}
		w := swar.Load8(s, n-8)
		if m := swar.ZeroLanes(w^bcA) | swar.ZeroLanes(w^bcB) | swar.ZeroLanes(w^bcC); m != 0 {
			return n - 8 + swar.FirstLane(m)
		}
		return -1
	}
	for i := range n {
		if s[i] == a || s[i] == b || s[i] == c {
			return i
		}
	}
	return -1
}

// IndexFold returns the index of the first ASCII case-insensitive occurrence of
// needle in s, or -1; an empty needle matches at 0. Only 'A'..'Z' and 'a'..'z'
// fold, so "no\rcache" does not match "no-cache". Time is linear in len(s) and len(needle)
// short of hash collisions: after too many failed candidates a Rabin-Karp search takes over,
// as in strings.Index.
func IndexFold[S byteSeq](s S, needle string) int {
	n, k := len(s), len(needle)
	if k == 0 {
		return 0
	}
	if k > n {
		return -1
	}
	if k > 8 {
		return indexFoldLong(s, needle)
	}

	table := caseconv.ToLowerTable
	first := table[needle[0]]
	last := n - k

	if n >= 8 {
		// Mask false positives only cost a verify; candidates run low to high,
		// so pos > last ends the search.
		bc1 := swar.Broadcast(first)
		bc2 := bc1
		if first >= 'a' && first <= 'z' {
			bc2 = swar.Broadcast(first - 0x20)
		}

		// Built on the first candidate: eagerly it costs +44% on 8-byte and +13% to +17% on 32-64B
		// misses (benchstat -count=10, Go 1.25, Apple M2 Pro).
		var fn foldedNeedle
		built := false

		i := 0
		for ; i <= n-8; i += 8 {
			w := swar.Load8(s, i)
			cand := swar.ZeroLanes(w^bc1) | swar.ZeroLanes(w^bc2)
			for cand != 0 {
				pos := i + swar.FirstLane(cand)
				if pos > last {
					return -1
				}
				if !built {
					fn = foldPrep(needle)
					built = true
				}
				if foldedWindowAt(s, pos, fn.mask) == fn.word {
					return pos
				}
				cand &= cand - 1
			}
		}
		// Tail: i = n - n%8, so the remaining starts hold at most 7 bytes of s.
		for ; i <= last; i++ {
			if table[s[i]] == first {
				if !built {
					fn = foldPrep(needle)
					built = true
				}
				if foldedWindowAt(s, i, fn.mask) == fn.word {
					return i
				}
			}
		}
		return -1
	}

outer:
	for i := 0; i <= last; i++ {
		if table[s[i]] != first {
			continue
		}
		for j := 1; j < k; j++ {
			if table[s[i+j]] != table[needle[j]] {
				continue outer
			}
		}
		return i
	}
	return -1
}

// indexFoldLong is IndexFold for a needle of 9 bytes or more: candidates are verified word by
// word until the failures add up, then Rabin-Karp finishes.
func indexFoldLong[S byteSeq](s S, needle string) int {
	n, k := len(s), len(needle)
	last := n - k
	first := caseconv.ToLowerTable[needle[0]]
	bc1 := swar.Broadcast(first)
	bc2 := bc1
	if first >= 'a' && first <= 'z' {
		bc2 = swar.Broadcast(first - 0x20)
	}

	work := 0 // the cost of the failed verifications, see indexFoldWorkSlack
	for i := 0; i <= n-8; i += 8 {
		w := swar.Load8(s, i)
		cand := swar.ZeroLanes(w^bc1) | swar.ZeroLanes(w^bc2)
		for cand != 0 {
			pos := i + swar.FirstLane(cand)
			if pos > last {
				return -1
			}
			if foldEqualAt(s, pos, needle) {
				return pos
			}
			if work += 1 + k>>6; work > indexFoldWorkSlack+pos>>1 {
				return indexFoldRK(s, pos+1, needle)
			}
			cand &= cand - 1
		}
	}
	return -1
}

// foldedNeedle is the short-needle (k <= 8) verify state: folded word plus lane mask.
type foldedNeedle struct {
	word, mask uint64
}

func foldPrep(needle string) foldedNeedle {
	return foldedNeedle{foldNeedle(needle), ^uint64(0) >> ((8 - len(needle)) * 8)}
}

// foldNeedle lower-cases needle (len 1..8) into one little-endian word, lanes past it zero.
func foldNeedle(needle string) uint64 {
	k := len(needle)
	var word uint64
	switch {
	case k == 8:
		word = swar.Load8(needle, 0)
	case k >= 4:
		// Two 4-byte windows; where they overlap, both hold the same bytes.
		word = uint64(load4(needle, 0)) | uint64(load4(needle, k-4))<<(8*(k-4))
	default:
		word = uint64(needle[0]) | uint64(needle[k>>1])<<(8*(k>>1)) | uint64(needle[k-1])<<(8*(k-1))
	}
	return swar.ToLowerWord(word)
}

// foldedWindowAt requires len(s) >= 8 and the window at pos to fit in s.
func foldedWindowAt[S byteSeq](s S, pos int, lenMask uint64) uint64 {
	if pos <= len(s)-8 {
		return swar.ToLowerWord(swar.Load8(s, pos)) & lenMask
	}
	from := len(s) - 8
	w := swar.ToLowerWord(swar.Load8(s, from)) >> (8 * (pos - from))
	return w & lenMask
}

// caseBitLanes has the bit that tells 'a' from 'A' in every lane; caseBitLanes32 in four.
const (
	caseBitLanes   = 0x20 * swar.Ones
	caseBitLanes32 = 0x20202020
)

// foldEqualWords reports whether x and y hold the same bytes lane for lane, ASCII letters
// compared without case.
func foldEqualWords(x, y uint64) bool {
	d := x ^ y
	if d&^caseBitLanes != 0 {
		return false // a lane differs in more than the case bit; only a shortcut for the test below
	}
	// The lanes that differ do so in the case bit alone, which is a match only for a letter.
	return d&^(swar.MatchRangeMask(x|caseBitLanes, 'a', 'z')>>2) == 0
}

// foldEqualAt requires len(needle) >= 8 and s[pos:pos+len(needle)] within s.
func foldEqualAt[S byteSeq](s S, pos int, needle string) bool {
	k := len(needle)
	i := 0
	for ; i <= k-8; i += 8 {
		if !foldEqualWords(swar.Load8(s, pos+i), swar.Load8(needle, i)) {
			return false
		}
	}
	if i == k {
		return true
	}
	return foldEqualWords(swar.Load8(s, pos+k-8), swar.Load8(needle, k-8))
}

const (
	// A failed long-needle verification costs 1 + len/64 units; past indexFoldWorkSlack plus half a
	// unit per byte scanned IndexFold switches, so the work stays near two words per byte.
	indexFoldWorkSlack = 64
	// primeRK is bytealg's Rabin-Karp multiplier.
	primeRK = 16777619
)

// indexFoldRK is Rabin-Karp over the lower-cased bytes of s[from:] for a needle of 9 bytes or more;
// -1 if no start remains.
func indexFoldRK[S byteSeq](s S, from int, needle string) int {
	k := len(needle)
	last := len(s) - k
	if from > last {
		return -1
	}
	table := caseconv.ToLowerTable
	var hashNeedle, hash, pow uint32 = 0, 0, 1 // pow is primeRK^k
	for j := range k {
		hashNeedle = hashNeedle*primeRK + uint32(table[needle[j]])
		hash = hash*primeRK + uint32(table[s[from+j]])
		pow *= primeRK
	}
	for i := from; ; i++ {
		if hash == hashNeedle && foldEqualAt(s, i, needle) {
			return i
		}
		if i == last {
			return -1
		}
		hash = hash*primeRK + uint32(table[s[i+k]]) - pow*uint32(table[s[i]])
	}
}

// ContainsFold reports whether s contains needle, ASCII case-insensitively.
func ContainsFold[S byteSeq](s S, needle string) bool {
	return IndexFold(s, needle) >= 0
}

// HasPrefixFold reports whether s begins with prefix, ASCII case-insensitively. An empty prefix
// matches any s.
func HasPrefixFold[S byteSeq](s S, prefix string) bool {
	k := len(prefix)
	if len(s) < k {
		return false
	}
	table := caseconv.ToLowerTable
	if k >= 8 {
		if table[s[0]] != table[prefix[0]] {
			return false
		}
		return foldEqualAt(s, 0, prefix)
	}
	if k >= 4 {
		// 4..7 bytes, as "bearer " and "basic " are: two overlapping 4-byte windows in one word,
		// the first of which settles most mismatches.
		lo, plo := load4(s, 0), load4(prefix, 0)
		if (lo^plo)&^caseBitLanes32 != 0 {
			return false
		}
		return foldEqualWords(uint64(lo)|uint64(load4(s, k-4))<<32, uint64(plo)|uint64(load4(prefix, k-4))<<32)
	}
	for j := range k {
		if table[s[j]] != table[prefix[j]] {
			return false
		}
	}
	return true
}

// HasSuffixFold is HasPrefixFold for a suffix; an empty suffix matches any s.
func HasSuffixFold[S byteSeq](s S, suffix string) bool {
	k := len(suffix)
	n := len(s)
	if n < k {
		return false
	}
	table := caseconv.ToLowerTable
	if k >= 8 {
		if table[s[n-1]] != table[suffix[k-1]] {
			return false
		}
		return foldEqualAt(s, n-k, suffix)
	}
	if k >= 4 {
		// See HasPrefixFold: "gzip", "+json" and "+xml" land here.
		lo, plo := load4(s, n-k), load4(suffix, 0)
		if (lo^plo)&^caseBitLanes32 != 0 {
			return false
		}
		return foldEqualWords(uint64(lo)|uint64(load4(s, n-4))<<32, uint64(plo)|uint64(load4(suffix, k-4))<<32)
	}
	for j := range k {
		if table[s[n-k+j]] != table[suffix[j]] {
			return false
		}
	}
	return true
}
