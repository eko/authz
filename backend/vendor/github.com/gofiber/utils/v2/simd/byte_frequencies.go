// Adapted from the coregex project, https://github.com/coregx/coregex
// (package simd), Copyright (c) 2025 Andrey Kolkov and contributors,
// MIT License. See the LICENSE file in this directory for the full text.

package simd

// byteFrequencies holds empirical byte frequency ranks derived from English
// text, source code, and binary data samples. Lower rank means rarer byte,
// which makes it a better anchor for a search prefilter. The approach
// mirrors Rust's memchr crate (https://github.com/BurntSushi/memchr).
// The table is unexported (read it through ByteRank) so callers cannot
// mutate the ranking Memmem's prefilter selection depends on.
var byteFrequencies = [256]byte{
	// 0x00-0x0F: control characters (generally rare)
	0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 1, 0, 0,
	// 0x10-0x1F: more control characters
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	// 0x20-0x2F: space and punctuation
	255, 60, 140, 50, 40, 35, 30, 160, 130, 130, 80, 55, 200, 140, 210, 100,
	// 0x30-0x3F: digits and more punctuation
	180, 190, 170, 150, 140, 140, 130, 120, 120, 120, 150, 100, 70, 160, 70, 50,
	// 0x40-0x4F: '@' and uppercase A-O
	25, 120, 80, 90, 85, 130, 75, 70, 80, 115, 30, 35, 90, 85, 100, 105,
	// 0x50-0x5F: uppercase P-Z and brackets
	80, 15, 100, 110, 115, 70, 45, 55, 20, 50, 10, 90, 60, 90, 20, 110,
	// 0x60-0x6F: backtick and lowercase a-o
	30, 225, 140, 170, 165, 245, 135, 130, 150, 200, 25, 65, 175, 155, 195, 205,
	// 0x70-0x7F: lowercase p-z and braces
	145, 15, 195, 200, 215, 150, 75, 95, 45, 120, 20, 85, 40, 85, 15, 0,
	// 0x80-0xFF: UTF-8 continuation/lead bytes (rare in typical text)
	5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5,
	5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5,
	5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5,
	5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5,
	5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5,
	5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5,
	5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5,
	5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5,
}

// ByteRank returns the empirical frequency rank of b. Lower values indicate
// rarer bytes, which are better anchors for search prefilters; SelectRareBytes
// and Memmem rank needle bytes with it.
func ByteRank(b byte) byte {
	return byteFrequencies[b]
}

// RareByteInfo describes the two rarest bytes of a needle, as selected by
// SelectRareBytes. When the needle contains at least two distinct byte
// values, Byte2 always differs from Byte1; only for empty, single-byte, or
// all-identical needles do Byte2/Index2 equal Byte1/Index1.
type RareByteInfo struct {
	// Byte1 is the rarest byte found in the needle.
	Byte1 byte
	// Byte2 is the rarest byte with a value different from Byte1, when the
	// needle has one; otherwise it equals Byte1.
	Byte2 byte
	// Index1 is the position of Byte1's first occurrence in the needle.
	Index1 int
	// Index2 is the position of Byte2's first occurrence in the needle.
	Index2 int
}

// SelectRareBytes returns the two rarest bytes of needle according to
// ByteRank, preferring distinct byte values: whenever the needle contains
// two distinct values, Byte2 is the rarest byte different from Byte1, so
// duplicated occurrences of the rarest byte never crowd out a usable
// second anchor. Memmem feeds the result to MemchrPair to build a highly
// selective substring prefilter; it is exported so downstream code can
// build its own prefilters the same way.
func SelectRareBytes(needle []byte) RareByteInfo {
	n := len(needle)
	if n == 0 {
		return RareByteInfo{}
	}

	byte1, idx1 := needle[0], 0
	byte2, idx2 := needle[0], 0
	have2 := false

	for i := 1; i < n; i++ {
		b := needle[i]
		switch {
		case byteFrequencies[b] < byteFrequencies[byte1]:
			// b is strictly rarer than byte1, so it also has a different
			// value; the demoted byte1 is at least as rare as any previous
			// byte2 and becomes the second anchor.
			byte2, idx2 = byte1, idx1
			byte1, idx1 = b, i
			have2 = true
		case b != byte1 && (!have2 || byteFrequencies[b] < byteFrequencies[byte2]):
			byte2, idx2 = b, i
			have2 = true
		}
	}
	if !have2 {
		// All bytes identical: no distinct second anchor exists.
		byte2, idx2 = byte1, idx1
	}

	return RareByteInfo{Byte1: byte1, Byte2: byte2, Index1: idx1, Index2: idx2}
}
