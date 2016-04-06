package main

import (
	"fmt"
	"io"
	"net/http"
	urllib "net/url"
	"os"
	"regexp"
	"strings"

	tldlib "github.com/jpillora/go-tld"
	"github.com/wangtuanjie/ip17mon"
)

// ViewLog
type ViewLog struct {
	url             string
	domain          string
	userAgent       string
	browser         string
	browserVersion  string
	platform        string
	platformVersion string
	isMobile        bool
	isWechat        bool
	referer         string
	cookieid        string
	width           string
	height          string
	color           string
	language        string
	title           string
	ip              string
	country         string
	province        string
	city            string
	operators       string
	source          string
	sourceKey       string
}

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
	patternAndroidVersion = regexp.MustCompile(`(?i)Android ([0-9.]+)`)
	patternMacVersion     = regexp.MustCompile(`(?i)Mac OS X ([0-9_]+)`)
	patternIphoneVersion  = regexp.MustCompile(`(?i)iPhone OS ([0-9_]+)`)
	patternBrowserChrome  = regexp.MustCompile(`(?i)chrome\/(\d+)`)
	patternBrowserIEOld   = regexp.MustCompile(`(?i)MSIE\s(\d+)`)
	patternBrowserIE      = regexp.MustCompile(`(?i)Trident\/\d+\.\d+;.*[rv:]+(\d+)`)
	patternBrowserFirefox = regexp.MustCompile(`(?i)firefox\/([\d]+)`)
	patternBrowserIOS     = regexp.MustCompile(`(?i)iphone os ([\d]+)`)
	patternBrowserAndroid = regexp.MustCompile(`(?i)android (\d\.\d)`)
	patternIPv4           = regexp.MustCompile(`\d+[.]\d+[.]\d+[.]\d+`)
)

var windowsVersionMap = map[string]string{
	"Windows NT 10.0": "Windows 10",
	"Windows NT 6.3":  "Windows 8.1",
	"Windows NT 6.2":  "Windows 8",
	"Windows NT 6.1":  "Windows 7",
	"Windows NT 6.0":  "Windows Vista",
	"Windows NT 5.2":  "Windows 2003",
	"Windows NT 5.1":  "Windows XP",
	"Windows NT 5.0":  "Windows 2000",
}

var refererMap = map[string]string{
	"google.com": "google",
	"bing.com":   "bing",
	"360.cn":     "360",
	"so.com":     "360",
	"haosou.com": "360",
	"baidu.com":  "baidu",
	"sogou.com":  "sogou",
}

var refererFromMap = map[string]string{
	"timeline":      "timeline",
	"groupmessage":  "groupmessage",
	"singlemessage": "singlemessage",
}

// 解析 PC/Mobile/Wechat
func parseMobile(userAgent []byte) (bool, bool) {
	var isMobile = patternMobile.Match(userAgent)
	var isWechat = patternWechat.Match(userAgent)
	return isMobile, isWechat
}

// 解析浏览器类型
func parseBrowser(userAgent []byte) (string, string) {
	var browser string
	var browserVersion string
	switch true {
	case patternBrowserChrome.Match(userAgent):
		match := patternBrowserChrome.FindSubmatch(userAgent)
		if len(match) == 2 {
			browser = "Chrome"
			browserVersion = browser + " " + string(match[1])
		}
	case patternBrowserIEOld.Match(userAgent):
		match := patternBrowserIEOld.FindSubmatch(userAgent)
		if len(match) == 2 {
			browser = "IE"
			browserVersion = browser + " " + string(match[1])
		}
	case patternBrowserIE.Match(userAgent):
		match := patternBrowserIE.FindSubmatch(userAgent)
		if len(match) == 2 {
			browser = "IE"
			browserVersion = browser + " " + string(match[1])
		}
	case patternBrowserFirefox.Match(userAgent):
		match := patternBrowserFirefox.FindSubmatch(userAgent)
		if len(match) == 2 {
			browser = "Firefox"
			browserVersion = browser + " " + string(match[1])
		}
	case patternBrowserIOS.Match(userAgent):
		match := patternBrowserIOS.FindSubmatch(userAgent)
		if len(match) == 2 {
			browser = "iOS"
			browserVersion = browser + " " + string(match[1])
		}
	case patternBrowserAndroid.Match(userAgent):
		match := patternBrowserIOS.FindSubmatch(userAgent)
		if len(match) == 2 {
			browser = "Android"
			browserVersion = string(match[0])
		}
	}
	return browser, browserVersion
}

// 解析平台类型
func parsePlatform(userAgent []byte) (string, string) {
	var platform string
	var platformVersion string
	switch true {
	case patternWindows.Match(userAgent):
		platform = "Windows"
		match := patternWindowsVersion.FindSubmatch(userAgent)
		if len(match) == 2 {
			platformVersion = string(match[0])
			if value, ok := windowsVersionMap[platformVersion]; ok {
				platformVersion = value
			}
		}
	case patternAndroid.Match(userAgent):
		platform = "Android"
		match := patternAndroidVersion.FindSubmatch(userAgent)
		if len(match) == 2 {
			platformVersion = string(match[0])
		}
	case patternIphone.Match(userAgent):
		platform = "iPhone"
		match := patternIphoneVersion.FindSubmatch(userAgent)
		if len(match) == 2 {
			platformVersion = strings.Replace(string(match[0]), "_", ".", -1)
		}
	case patternMac.Match(userAgent):
		platform = "Mac"
		match := patternMacVersion.FindSubmatch(userAgent)
		if len(match) == 2 {
			platformVersion = strings.Replace(string(match[0]), "_", ".", -1)
		}
	case patternLinux.Match(userAgent):
		platform = "Linux"
		platformVersion = "Linux"
	default:
		platform = "other"
	}
	return platform, platformVersion
}

func parseIPAddress(IP string) (string, string, string, string) {
	loc, err := ip17mon.Find(IP)
	if err == nil {
		return loc.Country, loc.Region, loc.City, loc.Isp
	}
	return "", "", "", ""
}

// parseIP
func parseIP(r *http.Request) string {
	XForwardedFor := r.Header.Get("X-Forwarded-For")
	IP := patternIPv4.FindAllString(XForwardedFor, -1)
	if len(IP) > 0 {
		return IP[len(IP)-1]
	}
	return strings.Split(r.RemoteAddr, ":")[0]
}

func parseSource(r *http.Request) (string, string) {
	var source string
	var sourceKey string
	refererURL := r.Header.Get("Referer")
	rootDomain := getRootDomain(refererURL)

	if value, ok := refererMap[rootDomain]; ok {
		source = value
		return source, sourceKey
	}

	url, err := urllib.Parse(refererURL)
	if err != nil {
		return "", ""
	}
	queryFrom := url.Query().Get("from")
	if value, ok := refererFromMap[queryFrom]; ok {
		source = value
	}
	return source, sourceKey
}

func boolToString(boolValue bool) string {
	if boolValue {
		return "true"
	}
	return "false"
}

func getAbsURI(r *http.Request) string {
	scheme := "http"
	if r.URL.Scheme != "" {
		scheme = r.URL.Scheme
	}
	return scheme + "://" + r.Host + r.RequestURI
}

func getRootDomain(url string) string {
	result, _ := tldlib.Parse(url)
	return result.Domain + "." + result.TLD
}

func handle(w http.ResponseWriter, r *http.Request) {

	url := r.Header.Get("Referer")
	domain := r.Host
	userAgent := []byte(r.Header.Get("User-Agent"))
	query := r.URL.Query()
	header := r.Header
	debug := query.Get("debug")

	viewlog := ViewLog{}
	viewlog.url = url
	viewlog.referer = query.Get("referer")
	viewlog.cookieid = query.Get("cookieid")
	viewlog.width = query.Get("width")
	viewlog.height = query.Get("height")
	viewlog.color = query.Get("color")
	viewlog.language = query.Get("language")
	viewlog.title = query.Get("title")
	viewlog.domain = domain
	viewlog.userAgent = header.Get("User-Agent")
	viewlog.isMobile, viewlog.isWechat = parseMobile(userAgent)
	viewlog.platform, viewlog.platformVersion = parsePlatform(userAgent)
	viewlog.browser, viewlog.browserVersion = parseBrowser(userAgent)
	viewlog.ip = parseIP(r)
	viewlog.country, viewlog.province, viewlog.city, viewlog.operators = parseIPAddress(viewlog.ip)
	viewlog.source, viewlog.sourceKey = parseSource(r)

	fmt.Println(viewlog)

	switch debug {
	case "mobile":
		io.WriteString(w, boolToString(viewlog.isMobile))
	case "wechat":
		io.WriteString(w, boolToString(viewlog.isWechat))
	case "platform":
		io.WriteString(w, viewlog.platform)
	case "ip":
		io.WriteString(w, viewlog.ip)
	case "source":
		io.WriteString(w, viewlog.source)
	default:
		io.WriteString(w, "Hello world!")
	}

}

func init() {
	if err := ip17mon.Init("17monipdb.dat"); err != nil {
		panic(err)
	}
}

func main() {
	http.HandleFunc("/", handle)
	http.ListenAndServe(":"+os.Args[1], nil)
}
