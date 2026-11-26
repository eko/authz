package utils

import (
	"time"
)

// httpDateLayout doubles as AppendHTTPDate's template: its separators stay in place.
const (
	httpDateLayout = "Mon, 02 Jan 2006 15:04:05 GMT"
	httpDateLen    = len(httpDateLayout)
)

var (
	httpWeekdays = [7]string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}
	httpMonths   = [12]string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}

	httpDaysInMonth = [13]int{0, 31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
)

// Hinnant's civil_from_days/days_from_civil; origin -0400-03-01 keeps values non-negative.
const (
	httpDateMinUnix        = -62167219200 // 0000-01-01T00:00:00Z
	httpDateMaxUnix        = 253402300799 // 9999-12-31T23:59:59Z
	secondsPerDay          = 86400
	daysPerEra             = 146097
	eraDaysToMarch         = daysPerEra - 60     // 0000-01-01 to -0400-03-01
	daysFromMarchEraToUnix = 719468 + daysPerEra // -0400-03-01 to 1970-01-01
	weekdayOfYearZero      = 6                   // 0000-01-01 was a Saturday
)

// civilFromDays converts days since 0000-01-01 (at most 3652424) to a proleptic Gregorian date.
func civilFromDays(days uint64) (year, month, day uint64) { //nolint:nonamedreturns // the three calendar fields are only readable named
	z := days + eraDaysToMarch
	era := z / daysPerEra
	doe := z - era*daysPerEra
	yoe := (doe - doe/1460 + doe/36524 - doe/146096) / 365
	doy := doe - (365*yoe + yoe/4 - yoe/100)
	mp := (5*doy + 2) / 153
	day = doy - (153*mp+2)/5 + 1
	month = mp + 3
	if mp >= 10 {
		month = mp - 9
	}
	year = yoe + era*400 - 400
	if month <= 2 {
		year++
	}
	return year, month, day
}

// unixFromCivil returns the Unix time of valid UTC fields with the year in 0..9999.
func unixFromCivil(year, month, day, hour, minute, sec int) int64 {
	y := year + 400 // one era back keeps January/February of year 0 non-negative
	mp := month + 9
	if month > 2 {
		mp = month - 3
	} else {
		y--
	}
	era := y / 400
	yoe := y - era*400
	doy := (153*mp+2)/5 + day - 1
	doe := yoe*365 + yoe/4 - yoe/100 + doy
	days := era*daysPerEra + doe - daysFromMarchEraToUnix
	return int64(days)*secondsPerDay + int64(hour*3600+minute*60+sec)
}

func putPair(b []byte, i int, n uint64) {
	_ = b[i+1]
	b[i] = decimalPairs[2*n]
	b[i+1] = decimalPairs[2*n+1]
}

// AppendHTTPDate appends t to dst in the RFC 9110 HTTP date format, byte-identical to
// t.UTC().AppendFormat(dst, http.TimeFormat) and always 29 bytes for the years 0..9999.
func AppendHTTPDate(dst []byte, t time.Time) []byte {
	sec := t.Unix()
	if sec < httpDateMinUnix || sec > httpDateMaxUnix {
		return t.UTC().AppendFormat(dst, httpDateLayout)
	}
	u := uint64(sec - httpDateMinUnix)
	days := u / secondsPerDay
	sod := u - days*secondsPerDay
	hour := sod / 3600
	sod -= hour * 3600
	minute := sod / 60
	second := sod - minute*60
	year, month, day := civilFromDays(days)

	var b [httpDateLen]byte
	copy(b[:], httpDateLayout)
	weekday := httpWeekdays[(days+weekdayOfYearZero)%7]
	b[0], b[1], b[2] = weekday[0], weekday[1], weekday[2]
	putPair(b[:], 5, day)
	monthName := httpMonths[month-1]
	b[8], b[9], b[10] = monthName[0], monthName[1], monthName[2]
	century := year / 100
	putPair(b[:], 12, century)
	putPair(b[:], 14, year-century*100)
	putPair(b[:], 17, hour)
	putPair(b[:], 20, minute)
	putPair(b[:], 23, second)
	return append(dst, b[:]...)
}

// FormatHTTPDate returns t.UTC().Format(http.TimeFormat) with a single allocation.
func FormatHTTPDate(t time.Time) string {
	if sec := t.Unix(); sec < httpDateMinUnix || sec > httpDateMaxUnix {
		// A year outside 0..9999 outgrows the stack buffer below: two allocations.
		return t.UTC().Format(httpDateLayout)
	}
	var buf [httpDateLen]byte
	return string(AppendHTTPDate(buf[:0], t))
}

// ParseHTTPDate parses an HTTP date the way net/http.ParseTime does: the RFC
// 9110 preferred format or the obsolete RFC 850 and ANSI C forms. Unlike
// ParseTime, it ignores space, tab, CR and LF around the date. Canonical
// preferred-format input yields a time in UTC; other forms take time.Parse's
// location and errors.
func ParseHTTPDate[S byteSeq](s S) (time.Time, error) {
	if t, ok := parseRFC1123(s); ok {
		return t, nil
	}
	return parseHTTPDateSlow(string(s))
}

// parseRFC1123 is the strict fast path for canonical input; anything else,
// out-of-range fields included, reports ok=false so time.Parse gives the errors.
func parseRFC1123[S byteSeq](s S) (time.Time, bool) {
	if len(s) != httpDateLen {
		return time.Time{}, false
	}
	if s[3] != ',' || s[4] != ' ' || s[7] != ' ' || s[11] != ' ' || s[16] != ' ' ||
		s[19] != ':' || s[22] != ':' || s[25] != ' ' ||
		s[26] != 'G' || s[27] != 'M' || s[28] != 'T' {
		return time.Time{}, false
	}
	if !isHTTPWeekday(pack3(s[0], s[1], s[2])) {
		return time.Time{}, false
	}
	month := httpMonthNum(pack3(s[8], s[9], s[10]))
	if month == 0 {
		return time.Time{}, false
	}
	day, okDay := twoDigits(s[5], s[6])
	yh, okYh := twoDigits(s[12], s[13])
	yl, okYl := twoDigits(s[14], s[15])
	hour, okHour := twoDigits(s[17], s[18])
	minute, okMin := twoDigits(s[20], s[21])
	sec, okSec := twoDigits(s[23], s[24])
	if !okDay || !okYh || !okYl || !okHour || !okMin || !okSec {
		return time.Time{}, false
	}
	year := yh*100 + yl
	if hour > 23 || minute > 59 || sec > 59 ||
		day < 1 || day > daysInHTTPMonth(month, year) {
		return time.Time{}, false
	}
	return time.Unix(unixFromCivil(year, month, day, hour, minute, sec), 0).UTC(), true
}

func parseHTTPDateSlow(s string) (time.Time, error) {
	start, end := 0, len(s)
	for start < end && isHTTPDateSpace(s[start]) {
		start++
	}
	for start < end && isHTTPDateSpace(s[end-1]) {
		end--
	}
	s = s[start:end]
	t, err := time.Parse(httpDateLayout, s)
	if err == nil {
		return t, nil
	}
	if t, err850 := time.Parse(time.RFC850, s); err850 == nil {
		return t, nil
	}
	t, errANSIC := time.Parse(time.ANSIC, s)
	if errANSIC == nil {
		return t, nil
	}
	return time.Time{}, errANSIC
}

// isHTTPDateSpace matches textproto.TrimString's notion of whitespace.
func isHTTPDateSpace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\r' || c == '\n'
}

func pack3(a, b, c byte) uint32 {
	return uint32(a)<<16 | uint32(b)<<8 | uint32(c)
}

// isHTTPWeekday checks the name only, not the date, like time.Parse.
func isHTTPWeekday(w uint32) bool {
	switch w {
	case 'S'<<16 | 'u'<<8 | 'n', 'M'<<16 | 'o'<<8 | 'n', 'T'<<16 | 'u'<<8 | 'e',
		'W'<<16 | 'e'<<8 | 'd', 'T'<<16 | 'h'<<8 | 'u', 'F'<<16 | 'r'<<8 | 'i',
		'S'<<16 | 'a'<<8 | 't':
		return true
	}
	return false
}

func httpMonthNum(w uint32) int {
	switch w {
	case 'J'<<16 | 'a'<<8 | 'n':
		return 1
	case 'F'<<16 | 'e'<<8 | 'b':
		return 2
	case 'M'<<16 | 'a'<<8 | 'r':
		return 3
	case 'A'<<16 | 'p'<<8 | 'r':
		return 4
	case 'M'<<16 | 'a'<<8 | 'y':
		return 5
	case 'J'<<16 | 'u'<<8 | 'n':
		return 6
	case 'J'<<16 | 'u'<<8 | 'l':
		return 7
	case 'A'<<16 | 'u'<<8 | 'g':
		return 8
	case 'S'<<16 | 'e'<<8 | 'p':
		return 9
	case 'O'<<16 | 'c'<<8 | 't':
		return 10
	case 'N'<<16 | 'o'<<8 | 'v':
		return 11
	case 'D'<<16 | 'e'<<8 | 'c':
		return 12
	}
	return 0
}

func twoDigits(a, b byte) (int, bool) {
	a -= '0'
	b -= '0'
	if a > 9 || b > 9 {
		return 0, false
	}
	return int(a)*10 + int(b), true
}

func daysInHTTPMonth(month, year int) int {
	if month == 2 && year%4 == 0 && (year%100 != 0 || year%400 == 0) {
		return 29
	}
	return httpDaysInMonth[month]
}
