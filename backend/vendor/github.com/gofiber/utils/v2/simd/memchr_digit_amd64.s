// Adapted from the coregex project, https://github.com/coregx/coregex
// (package simd), Copyright (c) 2025 Andrey Kolkov and contributors,
// MIT License. See the LICENSE file in this directory for the full text.
// Modified: legacy-SSE MOVD/MOVQ register moves replaced with VEX-encoded
// VMOVD/VMOVQ to avoid AVX-SSE transition penalties on Intel CPUs, the
// scalar tail loop replaced with one overlapping vector at the buffer end
// for the MinLen+ inputs the Go dispatch guarantees, and the 32-byte main
// loop replaced with a 4x unrolled 128-byte block.

//go:build amd64

#include "textflag.h"
#include "block_align_amd64.h"

// DIGITS leaves in dst a mask with 0xFF in every lane of src holding an
// ASCII digit. t0 is scratch; dst may alias src.
//
//	Y0 = [0x46 x 32] (0x7F - '9')   Y1 = [0x75 x 32] ('0' + 0x46 - 1)
//
// A range test needs only one compare when the range is first slid up
// against the signed maximum: adding 0x46 maps '0'..'9' onto 0x76..0x7F,
// the ten largest positive int8 values, so a signed "> 0x75" is exact for
// every byte value. Everything above '9' lands in 0x80..0xFF and reads as
// negative, everything below '0' stays at or under 0x75, and the byte
// values from 0xBA up wrap around to 0x00..0x45, still under the
// threshold. The earlier two-compare form (> 0x2F, then not > 0x39, AND-ed
// with VPANDN) cost three vector ops and two scratch registers per vector;
// this is two ops and one.
#define DIGITS(src, dst, t0) \
	VPADDB   Y0, src, t0; \
	VPCMPGTB Y1, t0, dst

// func memchrDigitAVX2(haystack []byte) int
//
// AVX2 implementation of digit search that finds the first ASCII digit [0-9].
//
// Algorithm for range check [0-9] (bytes 0x30-0x39):
//  1. Broadcast the 0x46 bias and the 0x75 threshold to YMM registers
//  2. Load 32 bytes from haystack
//  3. VPADDB: slide the digit range up to 0x76..0x7F
//  4. VPCMPGTB: signed compare against 0x75 (see DIGITS above)
//  5. VPMOVMSKB: extract 32-bit mask
//  6. BSFL: find first set bit
//
// The block loop runs steps 2-4 over four vectors, ORs the four masks and
// extracts once; the (rare) hit path re-extracts the per-vector masks in
// address order to locate the first digit.
//
// Parameters (FP offsets):
//   haystack_base+0(FP)  - pointer to haystack data (8 bytes)
//   haystack_len+8(FP)   - haystack length (8 bytes)
//   haystack_cap+16(FP)  - haystack capacity (8 bytes, unused but part of slice)
//   ret+24(FP)           - return value: index or -1 (8 bytes)
//
// Total argument frame size: 32 bytes (8+8+8+8)
TEXT ·memchrDigitAVX2(SB), NOSPLIT, $0-32
	// Load parameters
	MOVQ    haystack_base+0(FP), SI     // SI = haystack pointer
	MOVQ    haystack_len+8(FP), DX      // DX = haystack length

	// Empty haystack check
	TESTQ   DX, DX
	JZ      not_found

	// Prepare constants for range check ['0'-'9'] = [0x30-0x39]:
	// the bias that slides the range up to 0x76..0x7F and the signed
	// threshold just below it.
	MOVQ    $0x4646464646464646, AX
	VMOVQ   AX, X0
	VPBROADCASTQ X0, Y0                  // Y0 = [0x46 x 32]

	MOVQ    $0x7575757575757575, AX
	VMOVQ   AX, X1
	VPBROADCASTQ X1, Y1                  // Y1 = [0x75 x 32]

	// Save start pointer for offset calculation, compute the end pointer
	MOVQ    SI, DI                       // DI = haystack start (preserved)
	LEAQ    (SI)(DX*1), R8               // R8 = SI + length (end pointer)

	CMPQ    DX, $32
	JB      tail_loop                    // sub-32 direct call: scalar

	LEAQ    -32(R8), R12                 // last valid 32-byte window base
	LEAQ    -128(R8), R11                // last valid 128-byte block base
	CMPQ    SI, R11
	JA      loop32_entry

	BLOCKALIGN
loop128:
	VMOVDQU (SI), Y2
	VMOVDQU 32(SI), Y3
	VMOVDQU 64(SI), Y4
	VMOVDQU 96(SI), Y5
	DIGITS(Y2, Y2, Y6)
	DIGITS(Y3, Y3, Y6)
	DIGITS(Y4, Y4, Y6)
	DIGITS(Y5, Y5, Y6)

	VPOR    Y2, Y3, Y6
	VPOR    Y4, Y5, Y7
	VPOR    Y6, Y7, Y6
	VPMOVMSKB Y6, CX
	TESTL   CX, CX
	JNZ     found_in_block

	ADDQ    $128, SI
	CMPQ    SI, R11
	JBE     loop128

loop32_entry:
	CMPQ    SI, R12
	JA      last32

loop32:
	VMOVDQU (SI), Y2
	DIGITS(Y2, Y2, Y6)
	VPMOVMSKB Y2, CX
	TESTL   CX, CX
	JNZ     found_in_vector

	ADDQ    $32, SI
	CMPQ    SI, R12
	JBE     loop32

last32:
	// Fewer than 32 bytes remain past SI; rescan the final window. Lanes
	// before SI were already resolved as non-digits, so BSF stays exact.
	CMPQ    SI, R8
	JAE     not_found
	MOVQ    R12, SI
	VMOVDQU (SI), Y2
	DIGITS(Y2, Y2, Y6)
	VPMOVMSKB Y2, CX
	TESTL   CX, CX
	JNZ     found_in_vector
	JMP     not_found

tail_loop:
	// Load one byte
	MOVBLZX (SI), BX                     // BX = haystack[SI] (zero-extended)

	// Check if byte is in range ['0'-'9'] (0x30-0x39)
	CMPB    BL, $0x30                    // Compare with '0'
	JB      tail_next                    // If < '0', not a digit
	CMPB    BL, $0x39                    // Compare with '9'
	JA      tail_next                    // If > '9', not a digit

	// Found digit!
	SUBQ    DI, SI                       // SI = offset from start
	MOVQ    SI, ret+24(FP)               // Return index
	VZEROUPPER                           // Clear upper YMM bits (CRITICAL!)
	RET

tail_next:
	// Advance to next byte
	INCQ    SI
	CMPQ    SI, R8                       // Check if reached end
	JB      tail_loop                    // Continue if SI < end

not_found:
	// No digit found in entire haystack
	MOVQ    $-1, AX                      // Return -1
	MOVQ    AX, ret+24(FP)
	VZEROUPPER                           // Clear upper YMM bits (CRITICAL!)
	RET

found_in_block:
	// Locate the first digit by re-extracting the per-vector masks in
	// address order; only one of the four can hold the lowest match.
	VPMOVMSKB Y2, CX
	TESTL   CX, CX
	JNZ     found_in_vector
	ADDQ    $32, SI
	VPMOVMSKB Y3, CX
	TESTL   CX, CX
	JNZ     found_in_vector
	ADDQ    $32, SI
	VPMOVMSKB Y4, CX
	TESTL   CX, CX
	JNZ     found_in_vector
	ADDQ    $32, SI
	VPMOVMSKB Y5, CX

found_in_vector:
	// Digit found in vector! CX contains a 32-bit mask with set bits at
	// digit positions; BSFL finds the index of the first one.
	BSFL    CX, BX                       // BX = index of first set bit (0-31)
	SUBQ    DI, SI                       // SI = offset from start to current chunk
	ADDQ    SI, BX                       // BX = absolute index
	MOVQ    BX, ret+24(FP)               // Return index
	VZEROUPPER                           // Clear upper YMM bits
	RET
