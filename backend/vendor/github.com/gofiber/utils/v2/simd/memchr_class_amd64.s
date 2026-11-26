// Adapted from the coregex project, https://github.com/coregx/coregex
// (package simd), Copyright (c) 2025 Andrey Kolkov and contributors,
// MIT License. See the LICENSE file in this directory for the full text.
// Modified: legacy-SSE MOVD/MOVQ register moves replaced with VEX-encoded
// VMOVD/VMOVQ to avoid AVX-SSE transition penalties on Intel CPUs, the
// scalar tail loops replaced with one overlapping vector at the buffer end
// for the MinLen+ inputs the Go dispatch guarantees, the per-range
// clamp-and-compare classifier replaced with a nibble table lookup, the
// 32-byte main loops replaced with 4x unrolled 128-byte blocks, and the
// block test reduced to one compare over the four raw class values.

//go:build amd64

#include "textflag.h"
#include "block_align_amd64.h"

// Classifying \w = [A-Za-z0-9_] with unsigned clamp-and-compare costs three
// VPMINUB/VPMAXUB/VPCMPEQB triples plus a VPCMPEQB for '_' and three VPORs:
// thirteen vector ops per 32 bytes, ten of them on the two general vector
// ALU ports. A nibble table does the same job in six (five in the block
// loops, which defer the compare to the reduced block value), and the two
// VPSHUFBs it adds run on the otherwise idle shuffle port.
//
// The table splits each byte into its high and low nibble. Every high
// nibble that contains word characters gets one bit:
//
//	0x3_ -> 0x01   0x4_ -> 0x02   0x5_ -> 0x04   0x6_ -> 0x08   0x7_ -> 0x10
//
// and each low-nibble entry carries the set of high-nibble rows in which
// that low nibble is a word character, so a byte is a word character
// exactly when lowTable[lo] & highTable[hi] != 0:
//
//	lo 0x0        -> 0x15   ('0', 'P', 'p')
//	lo 0x1..0x9   -> 0x1F   (digits and letters in every row)
//	lo 0xA        -> 0x1E   ('J', 'Z', 'j', 'z' but not ':')
//	lo 0xB..0xE   -> 0x0A   ('K'..'N', 'k'..'n' only)
//	lo 0xF        -> 0x0E   ('O', '_', 'o')
//
// Bytes >= 0x80 need no separate guard: VPSHUFB zeroes any lane whose
// index byte has bit 7 set, so the low-nibble lookup already returns 0 for
// them and they classify as non-word, matching the SWAR fallback.
//
// Both kernels compute the same per-lane class value: the AND of the two
// lookups, which is nonzero exactly in the word-character lanes. The
// 32-byte paths compare it against zero at once, giving the "not a word
// character" mask (0xFF per non-word lane) that VPMOVMSKB can extract;
// memchrNotWord scans those bits directly and memchrWord inverts them
// with a single NOTL. The block loops instead reduce the four raw class
// values first and compare only the reduction: memchrNotWord takes the
// per-lane minimum (zero exactly where some vector has a non-word lane),
// memchrWord the per-lane OR (nonzero exactly where some vector has a
// word lane). That is three reductions and one compare per block instead
// of four compares and three reductions, and the four raw values stay
// live so the (rare) hit path can still compare them one at a time.

// WORDBITS leaves in dst the class value of every lane of src: nonzero
// for a word character, zero otherwise. t0 is scratch; dst may alias src.
#define WORDBITS(src, dst, t0) \
	VPSRLW   $4, src, t0;  \
	VPAND    Y2, t0, t0;   \
	VPSHUFB  t0, Y1, t0;   \
	VPSHUFB  src, Y0, dst; \
	VPAND    t0, dst, dst

// NOTWORD leaves in dst a mask with 0xFF in every lane of src that is NOT
// a word character. t0 is scratch; dst may alias src.
#define NOTWORD(src, dst, t0) \
	WORDBITS(src, dst, t0); \
	VPCMPEQB Y3, dst, dst

// SETUPTABLES loads the nibble tables, the low-nibble mask and a zero
// vector into Y0-Y3. The tables are built from immediates rather than a
// RODATA section so the kernels stay position-independent and free of
// relocations.
#define SETUPTABLES \
	MOVQ        $0x1F1F1F1F1F1F1F15, AX; \
	VMOVQ       AX, X0;                  \
	MOVQ        $0x0E0A0A0A0A1E1F1F, AX; \
	VPINSRQ     $1, AX, X0, X0;          \
	VINSERTI128 $1, X0, Y0, Y0;          \
	MOVQ        $0x1008040201000000, AX; \
	VMOVQ       AX, X1;                  \
	VINSERTI128 $1, X1, Y1, Y1;          \
	MOVQ        $0x0F0F0F0F0F0F0F0F, AX; \
	VMOVQ       AX, X2;                  \
	VPBROADCASTQ X2, Y2;                 \
	VPXOR       Y3, Y3, Y3

// func memchrWordAVX2(haystack []byte) int
//
// AVX2 implementation to find the first word character [A-Za-z0-9_].
//
// Parameters:
//   haystack_base+0(FP)  - pointer to haystack (8 bytes)
//   haystack_len+8(FP)   - length (8 bytes)
//   haystack_cap+16(FP)  - capacity (8 bytes, unused)
//   ret+24(FP)           - return value (8 bytes)
//
TEXT ·memchrWordAVX2(SB), NOSPLIT, $0-32
	// Load parameters
	MOVQ    haystack_base+0(FP), SI    // SI = haystack pointer
	MOVQ    haystack_len+8(FP), DX     // DX = length

	// Empty check
	TESTQ   DX, DX
	JZ      word_not_found

	// Save start for offset calculation and compute the end pointer
	MOVQ    SI, DI                     // DI = start pointer
	LEAQ    (SI)(DX*1), R8             // R8 = end pointer

	CMPQ    DX, $32
	JB      word_tail_loop             // sub-32 direct call: scalar

	SETUPTABLES

	LEAQ    -32(R8), R12               // last valid 32-byte window base
	LEAQ    -128(R8), R11              // last valid 128-byte block base
	CMPQ    SI, R11
	JA      word_loop32_entry

	BLOCKALIGN_LONG
word_loop128:
	VMOVDQU (SI), Y4
	VMOVDQU 32(SI), Y5
	VMOVDQU 64(SI), Y6
	VMOVDQU 96(SI), Y7
	WORDBITS(Y4, Y4, Y8)
	WORDBITS(Y5, Y5, Y8)
	WORDBITS(Y6, Y6, Y8)
	WORDBITS(Y7, Y7, Y8)

	// A word character exists in the block iff some lane's class value is
	// nonzero in some vector, so OR the four values and look for a lane
	// that is not zero.
	VPOR    Y4, Y5, Y8
	VPOR    Y6, Y7, Y9
	VPOR    Y8, Y9, Y8
	VPCMPEQB Y3, Y8, Y8
	VPMOVMSKB Y8, CX
	NOTL    CX
	TESTL   CX, CX
	JNZ     word_found_in_block

	ADDQ    $128, SI
	CMPQ    SI, R11
	JBE     word_loop128

word_loop32_entry:
	CMPQ    SI, R12
	JA      word_last32

word_loop32:
	VMOVDQU (SI), Y4
	NOTWORD(Y4, Y4, Y8)
	VPMOVMSKB Y4, CX
	NOTL    CX
	TESTL   CX, CX
	JNZ     word_found_in_vector

	ADDQ    $32, SI
	CMPQ    SI, R12
	JBE     word_loop32

word_last32:
	// Fewer than 32 bytes remain past SI; rescan the final window. Lanes
	// before SI were already resolved as non-matching, so BSF stays exact.
	CMPQ    SI, R8
	JAE     word_not_found
	MOVQ    R12, SI
	VMOVDQU (SI), Y4
	NOTWORD(Y4, Y4, Y8)
	VPMOVMSKB Y4, CX
	NOTL    CX
	TESTL   CX, CX
	JNZ     word_found_in_vector
	JMP     word_not_found

word_tail_loop:
	MOVBLZX (SI), AX                   // AX = byte

	// Check [A-Z]: 65 <= b <= 90
	CMPB    AL, $65
	JB      word_check_lower
	CMPB    AL, $90
	JBE     word_found_scalar

word_check_lower:
	// Check [a-z]: 97 <= b <= 122
	CMPB    AL, $97
	JB      word_check_digit
	CMPB    AL, $122
	JBE     word_found_scalar

word_check_digit:
	// Check [0-9]: 48 <= b <= 57
	CMPB    AL, $48
	JB      word_check_underscore
	CMPB    AL, $57
	JBE     word_found_scalar

word_check_underscore:
	// Check '_': 95
	CMPB    AL, $95
	JE      word_found_scalar

	// Not a word char, continue
	INCQ    SI
	CMPQ    SI, R8
	JB      word_tail_loop

word_not_found:
	MOVQ    $-1, AX
	MOVQ    AX, ret+24(FP)
	VZEROUPPER
	RET

word_found_in_block:
	// Locate the first hit by comparing and extracting the per-vector
	// class values in address order; only one of the four can hold the
	// lowest match.
	VPCMPEQB Y3, Y4, Y4
	VPMOVMSKB Y4, CX
	NOTL    CX
	TESTL   CX, CX
	JNZ     word_found_in_vector
	ADDQ    $32, SI
	VPCMPEQB Y3, Y5, Y5
	VPMOVMSKB Y5, CX
	NOTL    CX
	TESTL   CX, CX
	JNZ     word_found_in_vector
	ADDQ    $32, SI
	VPCMPEQB Y3, Y6, Y6
	VPMOVMSKB Y6, CX
	NOTL    CX
	TESTL   CX, CX
	JNZ     word_found_in_vector
	ADDQ    $32, SI
	VPCMPEQB Y3, Y7, Y7
	VPMOVMSKB Y7, CX
	NOTL    CX

word_found_in_vector:
	BSFL    CX, BX
	SUBQ    DI, SI                     // SI = offset to chunk start
	ADDQ    SI, BX                     // BX = absolute position
	MOVQ    BX, ret+24(FP)
	VZEROUPPER
	RET

word_found_scalar:
	SUBQ    DI, SI
	MOVQ    SI, ret+24(FP)
	VZEROUPPER
	RET

// func memchrNotWordAVX2(haystack []byte) int
//
// AVX2 implementation to find the first non-word character, the complement
// of memchrWordAVX2. The nibble classifier already produces the non-word
// mask directly, so this kernel skips the NOTL memchrWordAVX2 needs.
//
TEXT ·memchrNotWordAVX2(SB), NOSPLIT, $0-32
	// Load parameters
	MOVQ    haystack_base+0(FP), SI
	MOVQ    haystack_len+8(FP), DX

	// Empty check
	TESTQ   DX, DX
	JZ      notword_not_found

	// Save start for offset calculation and compute the end pointer
	MOVQ    SI, DI
	LEAQ    (SI)(DX*1), R8

	CMPQ    DX, $32
	JB      notword_tail_loop          // sub-32 direct call: scalar

	SETUPTABLES

	LEAQ    -32(R8), R12
	LEAQ    -128(R8), R11
	CMPQ    SI, R11
	JA      notword_loop32_entry

	BLOCKALIGN_LONG
notword_loop128:
	VMOVDQU (SI), Y4
	VMOVDQU 32(SI), Y5
	VMOVDQU 64(SI), Y6
	VMOVDQU 96(SI), Y7
	WORDBITS(Y4, Y4, Y8)
	WORDBITS(Y5, Y5, Y8)
	WORDBITS(Y6, Y6, Y8)
	WORDBITS(Y7, Y7, Y8)

	// A non-word character exists in the block iff some lane's class
	// value is zero in some vector, so take the per-lane minimum of the
	// four values and look for a zero lane.
	VPMINUB Y4, Y5, Y8
	VPMINUB Y6, Y7, Y9
	VPMINUB Y8, Y9, Y8
	VPCMPEQB Y3, Y8, Y8
	VPMOVMSKB Y8, CX
	TESTL   CX, CX
	JNZ     notword_found_in_block

	ADDQ    $128, SI
	CMPQ    SI, R11
	JBE     notword_loop128

notword_loop32_entry:
	CMPQ    SI, R12
	JA      notword_last32

notword_loop32:
	VMOVDQU (SI), Y4
	NOTWORD(Y4, Y4, Y8)
	VPMOVMSKB Y4, CX
	TESTL   CX, CX
	JNZ     notword_found_in_vector

	ADDQ    $32, SI
	CMPQ    SI, R12
	JBE     notword_loop32

notword_last32:
	CMPQ    SI, R8
	JAE     notword_not_found
	MOVQ    R12, SI
	VMOVDQU (SI), Y4
	NOTWORD(Y4, Y4, Y8)
	VPMOVMSKB Y4, CX
	TESTL   CX, CX
	JNZ     notword_found_in_vector
	JMP     notword_not_found

notword_tail_loop:
	MOVBLZX (SI), AX

	// Check if NOT a word char (inverse of word_tail_loop logic)
	// Word chars: [A-Z], [a-z], [0-9], '_'
	// If it's a word char, skip; if not, found

	// Check [A-Z]: 65 <= b <= 90
	CMPB    AL, $65
	JB      notword_check_lower_gap
	CMPB    AL, $90
	JBE     notword_is_word           // It's [A-Z], skip

notword_check_lower_gap:
	// Check [a-z]: 97 <= b <= 122
	CMPB    AL, $97
	JB      notword_check_digit_gap
	CMPB    AL, $122
	JBE     notword_is_word           // It's [a-z], skip

notword_check_digit_gap:
	// Check [0-9]: 48 <= b <= 57
	CMPB    AL, $48
	JB      notword_check_underscore_gap
	CMPB    AL, $57
	JBE     notword_is_word           // It's [0-9], skip

notword_check_underscore_gap:
	// Check '_': 95
	CMPB    AL, $95
	JE      notword_is_word           // It's '_', skip

	// Not a word char - found!
	JMP     notword_found_scalar

notword_is_word:
	// It's a word char, continue searching
	INCQ    SI
	CMPQ    SI, R8
	JB      notword_tail_loop

notword_not_found:
	MOVQ    $-1, AX
	MOVQ    AX, ret+24(FP)
	VZEROUPPER
	RET

notword_found_in_block:
	VPCMPEQB Y3, Y4, Y4
	VPMOVMSKB Y4, CX
	TESTL   CX, CX
	JNZ     notword_found_in_vector
	ADDQ    $32, SI
	VPCMPEQB Y3, Y5, Y5
	VPMOVMSKB Y5, CX
	TESTL   CX, CX
	JNZ     notword_found_in_vector
	ADDQ    $32, SI
	VPCMPEQB Y3, Y6, Y6
	VPMOVMSKB Y6, CX
	TESTL   CX, CX
	JNZ     notword_found_in_vector
	ADDQ    $32, SI
	VPCMPEQB Y3, Y7, Y7
	VPMOVMSKB Y7, CX

notword_found_in_vector:
	BSFL    CX, BX
	SUBQ    DI, SI
	ADDQ    SI, BX
	MOVQ    BX, ret+24(FP)
	VZEROUPPER
	RET

notword_found_scalar:
	SUBQ    DI, SI
	MOVQ    SI, ret+24(FP)
	VZEROUPPER
	RET
