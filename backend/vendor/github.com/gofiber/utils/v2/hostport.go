package utils

import (
	"bytes"
	"math/bits"

	"github.com/gofiber/utils/v2/internal/unsafeconv"
	"github.com/gofiber/utils/v2/swar"
)

// SplitHostPort splits "host:port", "[host]:port", and their %zone forms
// with net.SplitHostPort's exact rules, reporting failure as ok == false
// instead of a *net.AddrError; it accepts byte slices as well as strings,
// and the parts alias hostport.
func SplitHostPort[S byteSeq](hostport S) (host, port S, ok bool) { //nolint:nonamedreturns // the two same-typed parts are only readable named
	b := unsafeconv.Bytes(hostport)
	if n := len(b); n >= 8 && n <= hostPortWordMax && b[0] != '[' {
		// Without brackets the last colon must be the only ':', '[' or ']' there is: one pass
		// over the words counts them.
		count, last := 0, 0
		i := 0
		for ; i <= n-8; i += 8 {
			if m := hostPortLanes(swar.Load8(b, i)); m != 0 {
				count += bits.OnesCount64(m)
				last = i + swar.LastLane(m)
			}
		}
		if i < n {
			// The last word overlaps the one before it, so only its new lanes count.
			old := uint(8 - (n - i))
			if m := hostPortLanes(swar.Load8(b, n-8)) &^ (1<<(8*old) - 1); m != 0 {
				count += bits.OnesCount64(m)
				last = n - 8 + swar.LastLane(m)
			}
		}
		if count != 1 || b[last] != ':' {
			return host, port, false
		}
		return hostport[:last], hostport[last+1:], true
	}
	// The port starts after the last colon.
	i := bytes.LastIndexByte(b, ':')
	if i < 0 {
		return host, port, false
	}
	j, k, hostEnd := 0, 0, i
	if hostport[0] == '[' {
		// Expect the first ']' just before the last ':'.
		end := bytes.IndexByte(b, ']')
		if end < 0 || end+1 != i {
			return host, port, false
		}
		j, k, hostEnd = 1, end+1, end // there can't be a '[' resp. ']' before these positions
	} else if bytes.IndexByte(b[:i], ':') >= 0 {
		return host, port, false
	}
	if bytes.IndexByte(b[j:], '[') >= 0 || bytes.IndexByte(b[k:], ']') >= 0 {
		return host, port, false
	}
	return hostport[j:hostEnd], hostport[i+1:], true
}

// hostPortWordMax is the longest input the word pass takes. The IndexByte scans that take the
// longer ones are vectorized and overtake it on a valid host:port: about level at 40 and at 48
// bytes (-4% and +8%, neither significant), it costs 18% at 56, 123% at 128 and 214% at 255
// (benchstat, ten layouts, Go 1.26, Xeon 2.1 GHz).
const hostPortWordMax = 40

// hostPortLanes flags the lanes of w that hold a ':', a '[' or a ']', exactly per lane. At cost
// 78 of the inlining budget 80 it has no room for more.
func hostPortLanes(w uint64) uint64 {
	return swar.MatchByteMask(w, ':') | swar.MatchByteMask(w, '[') | swar.MatchByteMask(w, ']')
}
