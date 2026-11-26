// Package swar provides SWAR (SIMD within a register) primitives that process
// eight ASCII bytes per uint64 word. They are exported so downstream packages
// can fuse their own scans.
//
// Results are per lane (exact for Match*, approximate for ZeroLanes) and
// identical on little- and big-endian platforms. Load8 and Store8 always use
// little-endian lane order (lane 0, the least significant byte, is the lowest
// byte index) and panic unless the eight bytes lie within len. No unsafe is used.
package swar

import (
	"math/bits"
)

const (
	// WordLen is the SWAR word width in bytes. Loop while i <= len(s)-WordLen,
	// not i+WordLen <= len(s), which can wrap near the largest int on 32-bit.
	WordLen = 8

	// Ones has 0x01 in every byte lane.
	Ones = 0x0101010101010101
	// HighBits has 0x80 in every byte lane; a Match* mask equals it iff every
	// lane matched.
	HighBits = 0x8080808080808080
	// LowSeven has 0x7F in every byte lane.
	LowSeven = 0x7f7f7f7f7f7f7f7f
)

// Load8 assembles s[i:i+8] into a little-endian uint64: s[i] becomes the
// least significant byte (lane 0). It panics unless i+8 <= len(s), also for a
// []byte with spare capacity. In word loops, pin each group with a three-index
// reslice; load at constant offsets.
func Load8[S ~string | ~[]byte](s S, i int) uint64 {
	_ = s[i+7] // len, not cap
	w := s[i : i+8]
	return uint64(w[0]) |
		uint64(w[1])<<8 |
		uint64(w[2])<<16 |
		uint64(w[3])<<24 |
		uint64(w[4])<<32 |
		uint64(w[5])<<40 |
		uint64(w[6])<<48 |
		uint64(w[7])<<56
}

// Store8 writes w into b[i:i+8] in Load8's lane order (lane 0 at b[i]). It
// panics unless i+8 <= len(b), also with spare capacity.
func Store8(b []byte, i int, w uint64) {
	d := b[:len(b):len(b)][i : i+8] //nolint:modernize // slices.Clip puts Store8 at cost 81 of the inlining budget 80
	d[0] = byte(w)
	d[1] = byte(w >> 8)
	d[2] = byte(w >> 16)
	d[3] = byte(w >> 24)
	d[4] = byte(w >> 32)
	d[5] = byte(w >> 40)
	d[6] = byte(w >> 48)
	d[7] = byte(w >> 56)
}

// Broadcast returns a word with c in every byte lane.
func Broadcast(c byte) uint64 {
	return uint64(c) * Ones
}

// ToLowerWord lower-cases every 'A'..'Z' lane in w. All other bytes,
// including those >= 0x80, pass through unchanged.
func ToLowerWord(w uint64) uint64 {
	return w | MatchRangeMask(w, 'A', 'Z')>>2 // 0x80>>2 == 0x20: set bit 5
}

// ToUpperWord upper-cases every 'a'..'z' lane in w. All other bytes,
// including those >= 0x80, pass through unchanged.
func ToUpperWord(w uint64) uint64 {
	return w &^ (MatchRangeMask(w, 'a', 'z') >> 2)
}

// ZeroLanes returns a mask flagging the zero-byte lanes of x; two ops cheaper
// than MatchByteMask for first-match scans: ZeroLanes(w ^ Broadcast(c)). The
// mask is zero iff x has no zero byte, every true zero lane is flagged (so OR
// or AND combinations of masks never lose a match), and the lowest set lane is
// the lowest true zero lane. Lanes above it may be false positives (borrows
// propagate upward): re-verify them, or use MatchByteMask for exact lanes.
func ZeroLanes(x uint64) uint64 {
	return (x - Ones) &^ x & HighBits
}

// MatchByteMask returns a word with 0x80 set in exactly the lanes of w whose
// byte equals c, and zero elsewhere. Unlike ZeroLanes it has no false
// positives, so its result is safe for LastLane or arbitrary lane masking.
func MatchByteMask(w uint64, c byte) uint64 {
	x := w ^ (uint64(c) * Ones)
	// x|HighBits keeps lanes >= 0x80: no borrow crosses lanes.
	return ^(((x | HighBits) - Ones) | x) & HighBits
}

// MatchRangeMask returns a word with 0x80 set in exactly the lanes of w whose
// byte is within [lo, hi]. It requires lo <= hi <= 0x7F, and bytes >= 0x80
// never match.
func MatchRangeMask(w uint64, lo, hi byte) uint64 {
	b := w & LowSeven
	// Masked lanes are <= 0x7F, so the biased sums cannot carry across lanes.
	ge := b + (0x80-uint64(lo))*Ones
	gt := b + (0x80-uint64(hi)-1)*Ones
	return ge &^ gt &^ w & HighBits
}

// FirstLane returns the index (0..7) of the lowest-addressed set lane in mask,
// or 8 if mask is zero. mask must contain only 0x80 lane markers.
func FirstLane(mask uint64) int {
	return bits.TrailingZeros64(mask) >> 3
}

// LastLane returns the index (0..7) of the highest-addressed set lane in mask,
// or -1 if mask is zero. mask must contain only 0x80 lane markers.
func LastLane(mask uint64) int {
	return (63 - bits.LeadingZeros64(mask)) >> 3
}
