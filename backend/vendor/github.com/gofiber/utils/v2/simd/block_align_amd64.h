// BLOCKALIGN and BLOCKALIGN_LONG place the 128-byte block loop that
// follows them at byte 18, respectively 27, of a 32-byte code window:
// PCALIGN pads to the window, then two (or three) nine-byte NOPs — the
// canonical 66 0F 1F 84 00 00 00 00 00 form — shift the loop head. The
// NOPs sit before the label, so they run once per call that enters the
// block loop, never per iteration.
//
// The offset matters because of the Intel JCC erratum: on Skylake-derived
// cores (Skylake, Cascade Lake, Coffee Lake and their siblings) the
// microcode mitigation refuses to keep a jump — or a macro-fused
// compare-and-jump pair — in the decoded-uop cache when it crosses a
// 32-byte boundary or ends exactly on one. A block loop whose TEST/JNZ,
// DEC/JZ or CMP/JBE lands on such a boundary then runs from the legacy
// decoder every iteration, which measured ~10% slower on the pair kernel
// at 4KiB, and whether that happened used to depend on where the linker
// placed each kernel. Which offsets keep every pair of a loop inside one
// window depends on the loop's exact byte layout: 18 does for the six
// loops that use BLOCKALIGN, 27 for the three that use BLOCKALIGN_LONG.
// block_align_amd64_test.go checks the property against the linked test
// binary and logs each pair's position relative to its loop head, so a
// kernel edit that moves a branch onto a boundary fails the test (and
// the log shows which offset to pick) instead of silently costing
// throughput.
#define BLOCKALIGN \
	PCALIGN $32; \
	BYTE $0x66; BYTE $0x0F; BYTE $0x1F; BYTE $0x84; BYTE $0x00; BYTE $0x00; BYTE $0x00; BYTE $0x00; BYTE $0x00; \
	BYTE $0x66; BYTE $0x0F; BYTE $0x1F; BYTE $0x84; BYTE $0x00; BYTE $0x00; BYTE $0x00; BYTE $0x00; BYTE $0x00

#define BLOCKALIGN_LONG \
	BLOCKALIGN; \
	BYTE $0x66; BYTE $0x0F; BYTE $0x1F; BYTE $0x84; BYTE $0x00; BYTE $0x00; BYTE $0x00; BYTE $0x00; BYTE $0x00
