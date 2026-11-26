package utils

import (
	"time"
)

const durationBufLen = 32 // the longest rendering is "-2562047h47m16.854775808s"

// AppendDuration appends d formatted exactly like time.Duration.String to dst
// and returns the extended slice, without the string allocation.
func AppendDuration(dst []byte, d time.Duration) []byte {
	var buf [durationBufLen]byte
	w := durationToBuf(&buf, d)
	return append(dst, buf[w:]...)
}

// durationToBuf renders d right-aligned into buf and returns the first offset.
func durationToBuf(buf *[durationBufLen]byte, d time.Duration) int {
	w := len(buf)
	u := uint64(d)
	neg := d < 0
	if neg {
		u = -u
	}

	w--
	buf[w] = 's'
	if u < uint64(time.Second) {
		// Sub-second: ns, µs (two UTF-8 bytes), or ms. The unit splits are
		// constant divisions, which compile to multiplies.
		w--
		switch {
		case u == 0:
			buf[w] = '0'
			return w
		case u < uint64(time.Microsecond):
			buf[w] = 'n'
		case u < uint64(time.Millisecond):
			w--
			buf[w] = 0xC2 // "µ" is 0xC2 0xB5
			buf[w+1] = 0xB5
			w = durationFrac(buf, w, u%uint64(time.Microsecond), 3)
			u /= uint64(time.Microsecond)
		default:
			buf[w] = 'm'
			w = durationFrac(buf, w, u%uint64(time.Millisecond), 6)
			u /= uint64(time.Millisecond)
		}
		w = durationInt(buf, w, u)
	} else {
		w = durationFrac(buf, w, u%uint64(time.Second), 9)
		u /= uint64(time.Second)
		w = durationInt(buf, w, u%60)
		u /= 60
		if u > 0 {
			w--
			buf[w] = 'm'
			w = durationInt(buf, w, u%60)
			u /= 60
			if u > 0 {
				w--
				buf[w] = 'h'
				w = durationInt(buf, w, u)
			}
		}
	}

	if neg {
		w--
		buf[w] = '-'
	}
	return w
}

// durationFrac writes the prec-digit fraction frac ending at buf[w] without
// its trailing zeros (nothing at all when it is zero, the common case) and
// returns the new offset.
func durationFrac(buf *[durationBufLen]byte, w int, frac uint64, prec int) int {
	if frac == 0 {
		return w
	}
	for frac%1000 == 0 {
		frac /= 1000
		prec -= 3
	}
	for frac%10 == 0 {
		frac /= 10
		prec--
	}
	for ; prec >= 2; prec -= 2 {
		q := frac / 100
		w -= 2
		putPair(buf[:], w, frac-q*100)
		frac = q
	}
	if prec == 1 {
		w--
		buf[w] = byte(frac) + '0'
	}
	w--
	buf[w] = '.'
	return w
}

// durationInt writes the digits of v ending at buf[w] and returns the new offset.
func durationInt(buf *[durationBufLen]byte, w int, v uint64) int {
	for v >= 100 {
		q := v / 100
		w -= 2
		putPair(buf[:], w, v-q*100)
		v = q
	}
	if v >= 10 {
		w -= 2
		putPair(buf[:], w, v)
		return w
	}
	w--
	buf[w] = byte(v) + '0'
	return w
}
