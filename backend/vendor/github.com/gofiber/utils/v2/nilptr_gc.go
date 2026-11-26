//go:build gc

package utils

import (
	"unsafe"
)

// mayBeNilPointer reports whether v's data word is nil, as a nil pointer's is; a nil map, func, chan
// or one-pointer struct matches too (nilReceiverString tells them apart). It relies on the gc
// interface layout: nilPointerReflect is exact but +2.4 to +4.2 ns (benchstat -count=10, Xeon 2.3 GHz).
func mayBeNilPointer(v any) bool {
	// #nosec G103
	return (*[2]unsafe.Pointer)(unsafe.Pointer(&v))[1] == nil
}
