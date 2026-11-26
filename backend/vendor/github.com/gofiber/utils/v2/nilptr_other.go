//go:build !gc

package utils

func mayBeNilPointer(v any) bool {
	return nilPointerReflect(v)
}
