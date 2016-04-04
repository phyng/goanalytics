package main

import (
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
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
	patternWindowsVersion = regexp.MustCompile(`(?i)Windows NT ([a-zA-Z0-9.]+)`)
	patternMacVersion     = regexp.MustCompile(`(?i)Mac OS X ([0-9_]+)`)
)

// 解析 PC/Mobile/Wechat
func parseMobile(userAgent []byte) (bool, bool, bool) {
	var isMobile = patternMobile.Match(userAgent)
	var isWechat = patternWechat.Match(userAgent)
	return !isMobile, isMobile, isWechat
}

// 解析平台类型
func parsePlatform(userAgent []byte) (string, string) {
	var platform string
	var platformVersion string
	switch true {
	case patternMac.Match(userAgent):
		platform = "Mac"
		match := patternMacVersion.FindSubmatch(userAgent)
		if len(match) == 2 {
			platformVersion = strings.Replace(string(match[0]), "_", ".", -1)
		}
	case patternWindows.Match(userAgent):
		platform = "Windows"
		match := patternWindowsVersion.FindSubmatch(userAgent)
		if len(match) == 2 {
			platformVersion = string(match[0])
		}
	case patternAndroid.Match(userAgent):
		platform = "Android"
	case patternIphone.Match(userAgent):
		platform = "iPhone"
	case patternLinux.Match(userAgent):
		platform = "Linux"
	default:
		platform = "other"
	}
	return platform, platformVersion
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
	var platform, platformVersion = parsePlatform(userAgent)

	println(host, userAgent)
	println(isPc, isMobile, isWechat, platform, platformVersion)

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
