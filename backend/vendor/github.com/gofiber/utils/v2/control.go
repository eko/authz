package utils

import (
	"github.com/gofiber/utils/v2/swar"
)

const noControlExemption = 0x80

// IndexControl returns the index of the first control byte (below 0x20 or DEL) in s, or -1.
func IndexControl[S byteSeq](s S) int {
	return scanControl(s, noControlExemption)
}

// IndexControlExceptTab is IndexControl with HTAB allowed.
func IndexControlExceptTab[S byteSeq](s S) int {
	return scanControl(s, '\t')
}

func scanControl[S byteSeq](s S, exempt byte) int {
	n := len(s)
	exemptLanes := swar.Broadcast(exempt)
	i := 0
	for ; i <= n-16; i += 16 {
		w := s[i : i+16]
		w0, w1 := swar.Load8(w, 0), swar.Load8(w, 8)
		m0, m1 := controlLanes(w0), controlLanes(w1)
		if m0|m1 != 0 {
			if m0 &= otherLanes(w0, exemptLanes); m0 != 0 {
				return i + swar.FirstLane(m0)
			}
			if m1 &= otherLanes(w1, exemptLanes); m1 != 0 {
				return i + 8 + swar.FirstLane(m1)
			}
		}
	}
	for ; i <= n-8; i += 8 {
		w := swar.Load8(s, i)
		if m := controlLanes(w); m != 0 {
			if m &= otherLanes(w, exemptLanes); m != 0 {
				return i + swar.FirstLane(m)
			}
		}
	}
	if i == n {
		return -1
	}
	if n >= 8 {
		w := swar.Load8(s, n-8)
		if m := controlLanes(w); m != 0 {
			if m &= otherLanes(w, exemptLanes); m != 0 {
				return n - 8 + swar.FirstLane(m)
			}
		}
		return -1
	}
	for ; i < n; i++ {
		if c := s[i]; (c < 0x20 || c == 0x7f) && c != exempt {
			return i
		}
	}
	return -1
}

// controlLanes flags the lanes of w below 0x20 or DEL, exactly per lane.
func controlLanes(w uint64) uint64 {
	b := w & swar.LowSeven
	return (^(b + (0x80-0x20)*swar.Ones) | (b + swar.Ones)) &^ w & swar.HighBits
}

// otherLanes flags the lanes of w that differ from lanes, exactly per lane.
func otherLanes(w, lanes uint64) uint64 {
	x := w ^ lanes
	return ((x | swar.HighBits) - swar.Ones) | x
}
