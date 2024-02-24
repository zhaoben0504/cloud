package tool

import (
	"crypto/md5"
	"encoding/hex"
	"golang.org/x/text/language"
	"io"
	"net/http"
	"strings"
	"time"
)

// UnixSecond unix time second
func UnixSecond() int64 {
	return time.Now().Unix()
}

// UnixMillisecond unix time millisecond
func UnixMillisecond() int64 {
	time.Now().UnixMilli()
	return time.Now().UnixNano() / 1e6
}

// MD5 md5 32 lowercase
func MD5(txt string) string {
	md5Hash := md5.New()
	_, _ = io.WriteString(md5Hash, txt)
	md5Bytes := md5Hash.Sum(nil)
	return strings.ToLower(hex.EncodeToString(md5Bytes))
}

// GetHeaderLanguage 获取头里的语言
func GetHeaderLanguage(header http.Header) string {
	defaultLang := language.Chinese.String()
	if header == nil {
		return defaultLang
	}
	lang := header.Get("Accept-Language")
	if lang == "" {
		return defaultLang
	}
	return lang
}

// GetContentType 获取Content-Type
func GetContentType(filename string) string {
	str := strings.ToUpper(filename)
	if strings.HasSuffix(str, "PNG") {
		return "image/png"
	}
	if strings.HasSuffix(str, "JPG") ||
		strings.HasSuffix(str, "JPE") ||
		strings.HasSuffix(str, "JPEG") {
		return "image/jpeg"
	}
	if strings.HasSuffix(str, "GIF") {
		return "image/gif"
	}
	if strings.HasSuffix(str, "SVG") {
		return "text/xml"
	}
	if strings.HasSuffix(str, "HTM") || strings.HasSuffix(str, "HTML") {
		return "text/html"
	}
	if strings.HasSuffix(str, "PDF") {
		return "application/pdf"
	}
	if strings.HasSuffix(str, "JSON") {
		return "application/json"
	}

	if strings.HasSuffix(str, "JS") {
		return "application/x-javascript"
	}
	if strings.HasSuffix(str, "CSS") {
		return "text/css"
	}
	if strings.HasSuffix(str, "TXT") {
		return "text/plain"
	}
	return "application/octet-stream"
}
