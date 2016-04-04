package main

import (
	"io"
	"net/http"
	"os"
	"regexp"
)

// 预编译正则表达式
var (
	patternMobile         = regexp.MustCompile(`(?i)mobile`)
	patternWechat         = regexp.MustCompile(`(?i)MicroMessenger`)
	patternWindows        = regexp.MustCompile(`(?i)windows nt`)
	patternMac            = regexp.MustCompile(`(?i)macintosh`)
	patternLinux          = regexp.MustCompile(`(?i)Linux`)
	patternAndroid        = regexp.MustCompile(`(?i)Android`)
	patternIphone         = regexp.MustCompile(`(?i)iPhone`)
	patternWindowsVersion = regexp.MustCompile(`(?i)Windows([a-zA-Z0-9.]+)`)
)

// 解析 PC/Mobile/Wechat
func parseMobile(userAgent []byte) (bool, bool, bool) {
	var isMobile = patternMobile.Match(userAgent)
	var isWechat = patternWechat.Match(userAgent)
	return !isMobile, isMobile, isWechat
}

// 解析平台类型
func parsePlatform(userAgent []byte) string {
	var platform string
	switch true {
	case patternMac.Match(userAgent):
		platform = "Mac"
	case patternWindows.Match(userAgent):
		platform = "Windows"
	case patternAndroid.Match(userAgent):
		platform = "Android"
	case patternIphone.Match(userAgent):
		platform = "iPhone"
	case patternLinux.Match(userAgent):
		platform = "Linux"
	default:
		platform = "other"
	}
	return platform
}

func booToString(boolValue bool) string {
	if boolValue {
		return "true"
	}
	return "false"
}

func handle(w http.ResponseWriter, r *http.Request) {
	var (
		host      = r.Host
		userAgent = []byte(r.Header.Get("User-Agent"))
	)

	debug := r.URL.Query().Get("debug")

	var isPc, isMobile, isWechat = parseMobile(userAgent)
	var platform = parsePlatform(userAgent)

	println(host, userAgent)
	println(isPc, isMobile, isWechat, platform)

	switch debug {
	case "mobile":
		io.WriteString(w, booToString(isMobile))
	case "wechat":
		io.WriteString(w, booToString(isWechat))
	case "platform":
		io.WriteString(w, platform)
	default:
		io.WriteString(w, "Hello world!")
	}

}

func main() {
	http.HandleFunc("/", handle)
	http.ListenAndServe(":"+os.Args[1], nil)
}
