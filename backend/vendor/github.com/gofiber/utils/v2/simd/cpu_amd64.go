//go:build amd64

package simd

import (
	"golang.org/x/sys/cpu"
)

// hasAVX2 gates every AVX2 kernel in this package. golang.org/x/sys/cpu
// verifies both the CPU capability and OS support for saving the YMM
// register state, and the value is read once at init so the per-call
// dispatch is a single predictable branch.
var hasAVX2 = cpu.X86.HasAVX2

// hasPOPCNT gates countNonASCIIAVX2, the one kernel that reaches outside
// the AVX/AVX2 instruction set: its remainder loop turns a vector's
// high-bit mask into a lane count with POPCNT, whose CPUID bit
// (01H:ECX[23]) is independent of AVX2's. Every shipping AVX2 part has
// POPCNT, but a hypervisor can mask the two separately, and an ungated
// POPCNT would be SIGILL rather than a fallback — so the count dispatch
// checks both and drops to the SWAR path if either is missing.
var hasPOPCNT = cpu.X86.HasPOPCNT
