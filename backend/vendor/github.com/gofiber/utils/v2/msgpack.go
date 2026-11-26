package utils

// MsgPackMarshal returns the MsgPack encoding of v.
type MsgPackMarshal func(v any) ([]byte, error)

// MsgPackUnmarshal parses the MsgPack-encoded data and stores the result
// in the value pointed to by v. If v is nil or not a pointer,
// Unmarshal returns an error.
type MsgPackUnmarshal func(data []byte, v any) error
