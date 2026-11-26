package utils

import (
	"bytes"
	"encoding/binary"
	"net/netip"

	"github.com/gofiber/utils/v2/internal/unsafeconv"
)

// ParseIPv4 parses a dotted-decimal IPv4 address into a netip.Addr. It
// accepts exactly the IPv4 strings netip.ParseAddr does — four octets 0-255,
// no leading zeros, nothing before or after — and reports ok=false for
// everything else, including IPv6 forms. It never allocates, so []byte
// callers skip both the string conversion and netip's error construction.
//
// IsIPv4 is the validate-only counterpart; it accepts the same strings.
func ParseIPv4[S byteSeq](s S) (netip.Addr, bool) {
	var b [4]byte
	n := len(s)
	i := 0
	for oct := range 4 {
		if oct > 0 {
			if i >= n || s[i] != '.' {
				return netip.Addr{}, false
			}
			i++
		}
		if i >= n {
			return netip.Addr{}, false
		}
		c := s[i] - '0'
		if c > 9 {
			return netip.Addr{}, false
		}
		v := uint32(c)
		i++
		digits := 1
		for i < n && digits < 3 {
			c = s[i] - '0'
			if c > 9 {
				break
			}
			v = v*10 + uint32(c)
			i++
			digits++
		}
		// Reject values above 255 and, like netip, leading zeros.
		if v > 255 || (digits > 1 && s[i-digits] == '0') {
			return netip.Addr{}, false
		}
		b[oct] = byte(v)
	}
	if i != n {
		return netip.Addr{}, false
	}
	return netip.AddrFrom4(b), true
}

// ParseIPv6 parses an IPv6 address into a netip.Addr. It accepts exactly
// the IPv6 strings netip.ParseAddr does — 16-bit hex fields, one optional
// "::", an optional embedded dotted-decimal IPv4 tail, and an optional
// non-empty "%zone" suffix — and reports ok=false for everything else,
// including plain IPv4 forms (use ParseIPv4 for those). The parse itself
// never allocates; only a present zone is materialized as a string because
// netip.Addr stores zones as strings.
//
// IsIPv6 is the validate-only counterpart, but it follows net.ParseIP, which
// rejects zones. The two therefore disagree on zoned input by design: prefer
// ParseIPv6 when a zone is acceptable, IsIPv6 when it is not.
func ParseIPv6[S byteSeq](s S) (netip.Addr, bool) {
	n := len(s)
	zone := bytes.IndexByte(unsafeconv.Bytes(s), '%')
	end := n
	if zone >= 0 {
		if zone == n-1 {
			// Empty zone.
			return netip.Addr{}, false
		}
		end = zone
	}

	var ip [16]byte
	off := 0
	ellipsis := -1
	i := 0
	if end >= 2 && s[0] == ':' && s[1] == ':' {
		ellipsis = 0
		i = 2
	} else if end > 0 && s[0] == ':' {
		// A single leading colon starts no field.
		return netip.Addr{}, false
	}

	for off < 16 && i < end {
		// Parse up to four hex digits; a fifth is an error.
		fieldStart := i
		val := uint32(0)
		lim := min(i+4, end)
		for i < lim {
			d := unhexTable[s[i]]
			if d == 0xFF {
				break
			}
			val = val<<4 | uint32(d)
			i++
		}
		if i == lim && i < end && unhexTable[s[i]] != 0xFF {
			return netip.Addr{}, false
		}
		if i < end && s[i] == '.' {
			// Embedded IPv4 tail filling the final four bytes.
			if off > 12 {
				return netip.Addr{}, false
			}
			v4, ok := ParseIPv4(s[fieldStart:end])
			if !ok {
				return netip.Addr{}, false
			}
			b4 := v4.As4()
			copy(ip[off:], b4[:])
			off += 4
			i = end
			break
		}
		if i == fieldStart {
			// Every colon-separated field needs at least one digit.
			return netip.Addr{}, false
		}
		// off is always even and < 16 here, so the two-byte store fits.
		binary.BigEndian.PutUint16(ip[off:], uint16(val))
		off += 2
		if i == end {
			break
		}
		if s[i] != ':' {
			return netip.Addr{}, false
		}
		i++
		if i < end && s[i] == ':' {
			if ellipsis >= 0 {
				// At most one "::".
				return netip.Addr{}, false
			}
			ellipsis = off
			i++
			continue
		}
		if i == end {
			// A trailing single colon is invalid.
			return netip.Addr{}, false
		}
	}
	if i != end {
		return netip.Addr{}, false
	}
	if off < 16 {
		if ellipsis < 0 {
			return netip.Addr{}, false
		}
		// Expand "::": shift the fields after it to the end, zero the gap.
		gap := 16 - off
		for j := off - 1; j >= ellipsis; j-- {
			ip[j+gap] = ip[j]
			ip[j] = 0
		}
	} else if ellipsis >= 0 {
		// "::" must stand for at least one field of zeros.
		return netip.Addr{}, false
	}

	addr := netip.AddrFrom16(ip)
	if zone >= 0 {
		addr = addr.WithZone(string(s[zone+1:]))
	}
	return addr, true
}
