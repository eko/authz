// Adapted from the coregex project, https://github.com/coregx/coregex
// (package simd), Copyright (c) 2025 Andrey Kolkov and contributors,
// MIT License. See the LICENSE file in this directory for the full text.

//go:build amd64

package simd

// The ASCII kernels are implemented in ascii_amd64.s. All three read the
// same signal (a byte is non-ASCII exactly when its high bit is set, which
// VPMOVMSKB gathers 32 at a time) and process four vectors per iteration.
// Like the rest of the package's kernels they are self-bounding: each one
// clamps its own vector loops and finishes short remainders itself, so
// none of them carries a length precondition for the caller to uphold.
//
//go:noescape
func isASCIIAVX2(data []byte) bool

//go:noescape
func firstNonASCIIAVX2(data []byte) int

//go:noescape
func countNonASCIIAVX2(data []byte) int

// isASCIIImpl dispatches to the AVX2 kernel for inputs large enough to amortize
// the vector setup; smaller inputs and pre-AVX2 CPUs take the SWAR path.
func isASCIIImpl(data []byte) bool {
	if hasAVX2 && len(data) >= MinLen {
		return isASCIIAVX2(data)
	}
	return isASCIIGeneric(data)
}

// firstNonASCIIImpl dispatches like isASCIIImpl.
func firstNonASCIIImpl(data []byte) int {
	if hasAVX2 && len(data) >= MinLen {
		return firstNonASCIIAVX2(data)
	}
	return firstNonASCIIGeneric(data)
}

// countNonASCIIImpl dispatches like isASCIIImpl, with the extra POPCNT
// check its kernel needs. The kernel counts its own sub-vector remainder,
// so there is no prefix/suffix split here. An earlier version did split,
// handing the remainder to countNonASCIIGeneric, and paid a second
// non-inlinable call even when the remainder was empty — enough to put
// this scan behind its own SWAR fallback at MinLen (11.3ns vs 8.8ns at
// 32B). With the split gone the two are level at 32B and the kernel pulls
// ahead from 64B up, so MinLen stays the shared threshold.
func countNonASCIIImpl(data []byte) int {
	if hasAVX2 && hasPOPCNT && len(data) >= MinLen {
		return countNonASCIIAVX2(data)
	}
	return countNonASCIIGeneric(data)
}
