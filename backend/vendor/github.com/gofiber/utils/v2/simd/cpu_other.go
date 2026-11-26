//go:build !amd64

package simd

// hasAVX2 is constant false off amd64, letting the compiler prune the
// accelerated branches entirely.
const hasAVX2 = false
