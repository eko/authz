package utils

import (
	"encoding/binary"

	"github.com/gofiber/utils/v2/internal/caseconv"
	"github.com/gofiber/utils/v2/internal/unsafeconv"
	"github.com/gofiber/utils/v2/swar"
)

// headerTokenTable flags the bytes valid in an HTTP header field name — the
// RFC 9110 token alphabet, matching net/http's validHeaderFieldByte.
var headerTokenTable = buildHeaderTokenTable()

func buildHeaderTokenTable() [256]bool {
	var t [256]bool
	for c := byte('0'); c <= '9'; c++ {
		t[c] = true
	}
	for c := byte('a'); c <= 'z'; c++ {
		t[c] = true
	}
	for c := byte('A'); c <= 'Z'; c++ {
		t[c] = true
	}
	for _, c := range []byte("!#$%&'*+-.^_`|~") {
		t[c] = true
	}
	return t
}

// CanonicalHeaderKey returns the canonical form of the HTTP header key:
// the first letter and any letter following a hyphen upper-cased, the rest
// lower-cased. The output matches net/http.CanonicalHeaderKey byte for
// byte, including its guard: a key containing any byte outside the header
// token alphabet is returned unchanged. Unlike the stdlib it accepts byte
// slices as well as strings and allocates only when the key actually needs
// rewriting — already-canonical keys, the overwhelmingly common case on
// receive paths, are validated a word at a time and returned as-is
// without the stdlib's canonical-key cache lookup.
func CanonicalHeaderKey[S byteSeq](key S) S {
	n := len(key)
	if n < 4 {
		return canonicalHeaderKeyTable(key)
	}
	// Keys of letters, digits and hyphens, as nearly all are, are checked a word at a time; any
	// other byte hands the key to the table loop.
	if n < 8 {
		// One word, the lanes past the key zero: two 4-byte windows, alike where they overlap.
		w := uint64(load4(key, 0)) | uint64(load4(key, n-4))<<(8*(n-4))
		letter := swar.MatchRangeMask(w|caseBitLanes, 'a', 'z')
		hyphen := swar.MatchByteMask(w, '-')
		if letter|hyphen|swar.MatchRangeMask(w, '0', '9')|swar.HighBits<<(8*n) != swar.HighBits {
			return canonicalHeaderKeyTable(key)
		}
		wrong := headerKeyWrong(w, letter, hyphen, headerKeyStart)
		if wrong == 0 {
			return key
		}
		w ^= wrong >> 2
		buf := make([]byte, n)
		for i := range buf {
			buf[i] = byte(w >> (8 * i))
		}
		return unsafeconv.Seq[S](buf) // sound for a string: buf is never written again
	}
	var wrong uint64
	for i := 0; ; i += 8 {
		p := min(i, n-8) // the last word overlaps the one before it
		w := swar.Load8(key, p)
		letter := swar.MatchRangeMask(w|caseBitLanes, 'a', 'z')
		hyphen := swar.MatchByteMask(w, '-')
		if letter|hyphen|swar.MatchRangeMask(w, '0', '9') != swar.HighBits {
			return canonicalHeaderKeyTable(key)
		}
		wrong |= headerKeyWrong(w, letter, hyphen, headerKeyAfter(key, p))
		if p == n-8 { // not i < n: i+8 can wrap near the largest int on 32-bit
			break
		}
	}
	if wrong == 0 {
		return key
	}
	return rewriteHeaderKey(key)
}

// headerKeyStart marks lane 0 as following a hyphen, which the start of a key counts as.
const headerKeyStart = 0x80

// headerKeyAfter is headerKeyStart if the word of key at p starts the key or follows a hyphen.
func headerKeyAfter[S byteSeq](key S, p int) uint64 {
	if p > 0 && key[p-1] != '-' {
		return 0
	}
	return headerKeyStart
}

// headerKeyWrong flags the lanes of w that hold a letter in the wrong case, given its letter
// and hyphen lanes and headerKeyAfter of it.
func headerKeyWrong(w, letter, hyphen, after uint64) uint64 {
	// w<<2 moves the case bit to the top of its lane: set for a lower-case letter, which is
	// wrong exactly where the lane follows a hyphen.
	return letter &^ (w<<2 ^ (hyphen<<8 | after))
}

// rewriteHeaderKey returns the canonical copy of a key of 8 bytes or more, all letters, digits
// and hyphens: flipping the case bit of the wrong lanes is the whole rewrite.
func rewriteHeaderKey[S byteSeq](key S) S {
	n := len(key)
	buf := make([]byte, n)
	for i := 0; ; i += 8 {
		p := min(i, n-8)
		w := swar.Load8(key, p)
		wrong := headerKeyWrong(w, swar.MatchRangeMask(w|caseBitLanes, 'a', 'z'), swar.MatchByteMask(w, '-'), headerKeyAfter(key, p))
		binary.LittleEndian.PutUint64(buf[p:], w^wrong>>2)
		if p == n-8 {
			break
		}
	}
	// The zero-copy view is sound: buf is uniquely owned here and is never
	// written again, which is what the string instantiation requires.
	return unsafeconv.Seq[S](buf)
}

// canonicalHeaderKeyTable is CanonicalHeaderKey a byte at a time, for keys with a token byte
// other than a letter, a digit or a hyphen, or with a byte that is no token at all.
func canonicalHeaderKeyTable[S byteSeq](key S) S {
	upper := true
	fix := -1
	for i := 0; i < len(key); i++ {
		c := key[i]
		if !headerTokenTable[c] {
			return key
		}
		if fix < 0 && (upper && 'a' <= c && c <= 'z' || !upper && 'A' <= c && c <= 'Z') {
			fix = i
		}
		upper = c == '-'
	}
	if fix < 0 {
		return key
	}

	buf := make([]byte, len(key))
	copy(buf, unsafeconv.Bytes(key))
	upper = fix == 0 || buf[fix-1] == '-'
	for i := fix; i < len(buf); i++ {
		c := buf[i]
		if upper {
			c = caseconv.ToUpperTable[c]
		} else {
			c = caseconv.ToLowerTable[c]
		}
		buf[i] = c
		upper = c == '-'
	}
	// The zero-copy view is sound: buf is uniquely owned here and is never
	// written again, which is what the string instantiation requires.
	return unsafeconv.Seq[S](buf)
}
