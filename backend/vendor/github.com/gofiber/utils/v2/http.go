// ⚡️ Fiber is an Express inspired web framework written in Go with ☕️
// 🤖 Github Repository: https://github.com/gofiber/fiber
// 📌 API Documentation: https://docs.gofiber.io

package utils

import (
	"mime"
	"strings"

	casestrings "github.com/gofiber/utils/v2/strings"
	"github.com/gofiber/utils/v2/swar"
)

const MIMEOctetStream = "application/octet-stream"
const (
	contentTypeApplicationJSON            = "application/json"
	contentTypeApplicationXML             = "application/xml"
	contentTypeApplicationFormURLEncoded  = "application/x-www-form-urlencoded"
	contentTypePrefixApplicationWithSlash = "application/"
)

// Extensions index an open-addressed table by the lower-cased extension packed into one word.
const (
	mimeTableBits       = 9
	mimeTableMask       = 1<<mimeTableBits - 1
	mimeTableMaxEntries = 1 << (mimeTableBits - 1) // at most half full, so probes end at a free slot
	mimeHashMul         = 0x9E3779B97F4A7C15       // Fibonacci hashing multiplier
	mimeKeyMaxLen       = swar.WordLen - 1
)

type mimeEntry struct {
	key      uint64
	mimeType string
}

var mimeTable = buildMIMETable(mimeExtensions)

func buildMIMETable(exts map[string]string) [1 << mimeTableBits]mimeEntry {
	if len(exts) > mimeTableMaxEntries {
		panic("utils: the MIME extension table holds at most " + FormatInt(mimeTableMaxEntries) + " entries; raise mimeTableBits")
	}
	var t [1 << mimeTableBits]mimeEntry
	for ext, mimeType := range exts {
		if ext == "" || len(ext) > mimeKeyMaxLen {
			panic("utils: MIME extension \"" + ext + "\" must be 1.." + FormatInt(mimeKeyMaxLen) + " bytes to fit a packed key")
		}
		key := packExtension(ext)
		h := mimeHash(key)
		for t[h].mimeType != "" {
			h = (h + 1) & mimeTableMask
		}
		t[h] = mimeEntry{key: key, mimeType: mimeType}
	}
	return t
}

// packExtension packs 1..7 bytes into a word with the length in the top lane, so "html\x00" and "html" differ.
func packExtension(ext string) uint64 {
	return foldNeedle(ext) | uint64(len(ext))<<(8*mimeKeyMaxLen)
}

func mimeHash(key uint64) int {
	return int(key * mimeHashMul >> (64 - mimeTableBits))
}

// GetMIME returns the media type of a file extension, without parameters, with or without the
// dot and in any case. The built-in table wins over mime.TypeByExtension, which serves other
// extensions from the system's MIME databases and mime.AddExtensionType. Unknown gives
// MIMEOctetStream, empty gives "".
func GetMIME(extension string) string {
	if len(extension) == 0 {
		return ""
	}

	ext := extension
	if ext[0] == '.' {
		ext = ext[1:]
	}
	if len(ext) > 0 && len(ext) <= mimeKeyMaxLen {
		key := packExtension(ext)
		for h := mimeHash(key); mimeTable[h].mimeType != ""; h = (h + 1) & mimeTableMask {
			if mimeTable[h].key == key {
				return mimeTable[h].mimeType
			}
		}
	}

	if extension[0] != '.' {
		extension = "." + extension
	}
	if foundMime := mime.TypeByExtension(extension); foundMime != "" {
		// mime adds parameters ("text/csv; charset=utf-8"); the table has none.
		if i := strings.IndexByte(foundMime, ';'); i != -1 {
			foundMime = TrimSpace(foundMime[:i])
		}
		return foundMime
	}

	return MIMEOctetStream
}

// ParseVendorSpecificContentType reduces a content type with a "+suffix" in its subtype to its
// type and suffix, without parameters: "application/vnd.api+json; version=1" is "application/json".
// Spaces and tabs before the ';' or at the end are ignored, and a '/' after the '+' stays in the
// result. A '+' that first appears in the parameters, as in a multipart boundary, only strips
// them, and the blanks before the ';' stay. A content type without a '+', or with no '/' before
// it, is returned as it is. Matching is case-sensitive; with caseInsensitive[0] set it ignores
// case and builds lower-case results.
func ParseVendorSpecificContentType(cType string, caseInsensitive ...bool) string {
	useLower := len(caseInsensitive) > 0 && caseInsensitive[0]

	working := cType
	if useLower {
		working = casestrings.ToLower(cType)
	}
	plusIndex := strings.IndexByte(working, '+')

	if plusIndex == -1 {
		return cType
	}

	var parsableType string
	if semiColonIndex := strings.IndexByte(working, ';'); semiColonIndex == -1 {
		parsableType = working[plusIndex+1:]
	} else if plusIndex < semiColonIndex {
		parsableType = working[plusIndex+1 : semiColonIndex]
	} else {
		return cType[:semiColonIndex]
	}

	// A '/' after the '+' belongs to the suffix or parameters; splicing there would invent a type.
	slashIndex := strings.IndexByte(working, '/')

	if slashIndex == -1 || slashIndex > plusIndex {
		return cType
	}

	for n := len(parsableType); n > 0 && (parsableType[n-1] == ' ' || parsableType[n-1] == '\t'); n-- {
		parsableType = parsableType[:n-1]
	}

	if strings.HasPrefix(working, contentTypePrefixApplicationWithSlash) {
		switch parsableType {
		case "json":
			return contentTypeApplicationJSON
		case "xml":
			return contentTypeApplicationXML
		case "x-www-form-urlencoded":
			return contentTypeApplicationFormURLEncoded
		}
	}

	return working[:slashIndex+1] + parsableType
}

const (
	statusMessageMin = 100
	statusMessageMax = 511
)

// StatusMessage returns the correct message for the provided HTTP statuscode
func StatusMessage(status int) string {
	if status < statusMessageMin || status > statusMessageMax {
		return ""
	}
	return statusMessage[status]
}

// NOTE: Keep this in sync with fiber's status code list
var statusMessage = []string{
	100: "Continue",
	101: "Switching Protocols",
	102: "Processing",
	103: "Early Hints",

	200: "OK",
	201: "Created",
	202: "Accepted",
	203: "Non-Authoritative Information",
	204: "No Content",
	205: "Reset Content",
	206: "Partial Content",
	207: "Multi-Status",
	208: "Already Reported",
	226: "IM Used",

	300: "Multiple Choices",
	301: "Moved Permanently",
	302: "Found",
	303: "See Other",
	304: "Not Modified",
	305: "Use Proxy",
	306: "Switch Proxy",
	307: "Temporary Redirect",
	308: "Permanent Redirect",

	400: "Bad Request",
	401: "Unauthorized",
	402: "Payment Required",
	403: "Forbidden",
	404: "Not Found",
	405: "Method Not Allowed",
	406: "Not Acceptable",
	407: "Proxy Authentication Required",
	408: "Request Timeout",
	409: "Conflict",
	410: "Gone",
	411: "Length Required",
	412: "Precondition Failed",
	413: "Request Entity Too Large",
	414: "Request URI Too Long",
	415: "Unsupported Media Type",
	416: "Requested Range Not Satisfiable",
	417: "Expectation Failed",
	418: "I'm a teapot",
	421: "Misdirected Request",
	422: "Unprocessable Entity",
	423: "Locked",
	424: "Failed Dependency",
	425: "Too Early",
	426: "Upgrade Required",
	428: "Precondition Required",
	429: "Too Many Requests",
	431: "Request Header Fields Too Large",
	451: "Unavailable For Legal Reasons",

	500: "Internal Server Error",
	501: "Not Implemented",
	502: "Bad Gateway",
	503: "Service Unavailable",
	504: "Gateway Timeout",
	505: "HTTP Version Not Supported",
	506: "Variant Also Negotiates",
	507: "Insufficient Storage",
	508: "Loop Detected",
	510: "Not Extended",
	511: "Network Authentication Required",
}

// Copied from https://github.com/nginx/nginx/blob/67d2a9541826ecd5db97d604f23460210fd3e517/conf/mime.types
// with xml as application/xml (RFC 7303), js as text/javascript (RFC 9239), msgpack as
// application/vnd.msgpack (IANA), and cbor (RFC 8949) and zst (RFC 8878) added, which nginx lacks.
var mimeExtensions = map[string]string{
	"html":    "text/html",
	"htm":     "text/html",
	"shtml":   "text/html",
	"css":     "text/css",
	"xml":     "application/xml",
	"cbor":    "application/cbor",
	"gif":     "image/gif",
	"jpeg":    "image/jpeg",
	"jpg":     "image/jpeg",
	"js":      "text/javascript",
	"atom":    "application/atom+xml",
	"rss":     "application/rss+xml",
	"mml":     "text/mathml",
	"txt":     "text/plain",
	"jad":     "text/vnd.sun.j2me.app-descriptor",
	"wml":     "text/vnd.wap.wml",
	"htc":     "text/x-component",
	"avif":    "image/avif",
	"png":     "image/png",
	"svg":     "image/svg+xml",
	"svgz":    "image/svg+xml",
	"tif":     "image/tiff",
	"tiff":    "image/tiff",
	"wbmp":    "image/vnd.wap.wbmp",
	"webp":    "image/webp",
	"ico":     "image/x-icon",
	"jng":     "image/x-jng",
	"bmp":     "image/x-ms-bmp",
	"woff":    "font/woff",
	"woff2":   "font/woff2",
	"jar":     "application/java-archive",
	"war":     "application/java-archive",
	"ear":     "application/java-archive",
	"json":    "application/json",
	"msgpack": "application/vnd.msgpack",
	"hqx":     "application/mac-binhex40",
	"doc":     "application/msword",
	"pdf":     "application/pdf",
	"ps":      "application/postscript",
	"eps":     "application/postscript",
	"ai":      "application/postscript",
	"rtf":     "application/rtf",
	"m3u8":    "application/vnd.apple.mpegurl",
	"kml":     "application/vnd.google-earth.kml+xml",
	"kmz":     "application/vnd.google-earth.kmz",
	"xls":     "application/vnd.ms-excel",
	"eot":     "application/vnd.ms-fontobject",
	"ppt":     "application/vnd.ms-powerpoint",
	"odg":     "application/vnd.oasis.opendocument.graphics",
	"odp":     "application/vnd.oasis.opendocument.presentation",
	"ods":     "application/vnd.oasis.opendocument.spreadsheet",
	"odt":     "application/vnd.oasis.opendocument.text",
	"pptx":    "application/vnd.openxmlformats-officedocument.presentationml.presentation",
	"xlsx":    "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
	"docx":    "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
	"wmlc":    "application/vnd.wap.wmlc",
	"wasm":    "application/wasm",
	"7z":      "application/x-7z-compressed",
	"cco":     "application/x-cocoa",
	"jardiff": "application/x-java-archive-diff",
	"jnlp":    "application/x-java-jnlp-file",
	"run":     "application/x-makeself",
	"pl":      "application/x-perl",
	"pm":      "application/x-perl",
	"prc":     "application/x-pilot",
	"pdb":     "application/x-pilot",
	"rar":     "application/x-rar-compressed",
	"rpm":     "application/x-redhat-package-manager",
	"sea":     "application/x-sea",
	"swf":     "application/x-shockwave-flash",
	"sit":     "application/x-stuffit",
	"tcl":     "application/x-tcl",
	"tk":      "application/x-tcl",
	"der":     "application/x-x509-ca-cert",
	"pem":     "application/x-x509-ca-cert",
	"crt":     "application/x-x509-ca-cert",
	"xpi":     "application/x-xpinstall",
	"xhtml":   "application/xhtml+xml",
	"xspf":    "application/xspf+xml",
	"zip":     "application/zip",
	"zst":     "application/zstd",
	"bin":     "application/octet-stream",
	"exe":     "application/octet-stream",
	"dll":     "application/octet-stream",
	"deb":     "application/octet-stream",
	"dmg":     "application/octet-stream",
	"iso":     "application/octet-stream",
	"img":     "application/octet-stream",
	"msi":     "application/octet-stream",
	"msp":     "application/octet-stream",
	"msm":     "application/octet-stream",
	"mid":     "audio/midi",
	"midi":    "audio/midi",
	"kar":     "audio/midi",
	"mp3":     "audio/mpeg",
	"ogg":     "audio/ogg",
	"m4a":     "audio/x-m4a",
	"ra":      "audio/x-realaudio",
	"3gpp":    "video/3gpp",
	"3gp":     "video/3gpp",
	"ts":      "video/mp2t",
	"mp4":     "video/mp4",
	"mpeg":    "video/mpeg",
	"mpg":     "video/mpeg",
	"mov":     "video/quicktime",
	"webm":    "video/webm",
	"flv":     "video/x-flv",
	"m4v":     "video/x-m4v",
	"mng":     "video/x-mng",
	"asx":     "video/x-ms-asf",
	"asf":     "video/x-ms-asf",
	"wmv":     "video/x-ms-wmv",
	"avi":     "video/x-msvideo",
}
