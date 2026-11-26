package utils

import (
	"net"
)

// hexFieldMaxLen is the number of hex digits an IPv6 field may hold.
const hexFieldMaxLen = 4

// IsIPv4 reports whether s is a dotted-decimal IPv4 address. It accepts
// exactly what net.ParseIP accepts for the IPv4 case, without the net.IP slice,
// so it makes no allocations. It is the validate-only spelling of ParseIPv4,
// which accepts the same strings and also returns the address.
func IsIPv4(s string) bool {
	_, ok := ParseIPv4(s)
	return ok
}

// IsIPv6 reports whether s is an IPv6 address. It accepts exactly what
// net.ParseIP accepts for the IPv6 case, without the net.IP slice, so it makes
// no allocations. That excludes a "%zone" suffix, because net.ParseIP rejects
// zones.
//
// ParseIPv6 follows netip.ParseAddr instead and does accept zones, so the two
// disagree on zoned input by design. Callers relying on IsIPv6 to reject zones
// must screen '%' themselves before substituting ParseIPv6.
func IsIPv6(s string) bool {
	ellipsis := -1 // position of ellipsis in ip

	// Might have leading ellipsis
	if len(s) >= 2 && s[0] == ':' && s[1] == ':' {
		ellipsis = 0
		s = s[2:]
		// Might be only ellipsis
		if len(s) == 0 {
			return true
		}
	}

	// Loop, parsing hex numbers followed by colon.
	i := 0
	for i < net.IPv6len {
		// Hex field. Bounding the digit count is what rejects "00001":
		// leading zeros hold the value below 0xFFFF however many of them
		// there are, so a value-only bound lets an over-long field through.
		// Four digits cannot exceed 0xFFFF, so no value check is needed.
		ci := 0
		for ; ci < len(s); ci++ {
			c := s[ci]
			if (c < '0' || c > '9') && (c < 'a' || c > 'f') && (c < 'A' || c > 'F') {
				break
			}
			if ci == hexFieldMaxLen {
				return false
			}
		}
		if ci == 0 {
			return false
		}

		if ci < len(s) && s[ci] == '.' {
			if ellipsis < 0 && i != net.IPv6len-net.IPv4len {
				return false
			}
			if i+net.IPv4len > net.IPv6len {
				return false
			}

			if !IsIPv4(s) {
				return false
			}

			s = ""
			i += net.IPv4len
			break
		}

		// Save this 16-bit chunk.
		i += 2

		// Stop at end of string.
		s = s[ci:]
		if len(s) == 0 {
			break
		}

		// Otherwise must be followed by colon and more.
		if s[0] != ':' || len(s) == 1 {
			return false
		}
		s = s[1:]

		// Look for ellipsis.
		if s[0] == ':' {
			if ellipsis >= 0 { // already have one
				return false
			}
			ellipsis = i
			s = s[1:]
			if len(s) == 0 { // can be at end
				break
			}
		}
	}

	// Must have used entire string.
	if len(s) != 0 {
		return false
	}

	// If didn't parse enough, expand ellipsis.
	if i < net.IPv6len {
		if ellipsis < 0 {
			return false
		}
	} else if ellipsis >= 0 {
		// Ellipsis must represent at least one 0 group.
		return false
	}
	return true
}
