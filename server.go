package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	urllib "net/url"
	"os"
	"regexp"
	"strings"
	"time"

	tldlib "github.com/jpillora/go-tld"
	"github.com/wangtuanjie/ip17mon"
	// "github.com/mattbaird/elastigo"
)

// ViewLog core data structure
type ViewLog struct {
	Created         string `json:"created"`
	URL             string `json:"url"`
	Domain          string `json:"domain"`
	UserAgent       string `json:"useragent"`
	Browser         string `json:"browser"`
	BrowserVersion  string `json:"browser_version"`
	Platform        string `json:"platform"`
	PlatformVersion string `json:"platform_version"`
	IsMobile        bool   `json:"is_mobile"`
	IsWechat        bool   `json:"is_wechat"`
	Referer         string `json:"referer"`
	Cookieid        string `json:"cookieid"`
	Width           string `json:"width"`
	Height          string `json:"height"`
	Color           string `json:"color"`
	Language        string `json:"language"`
	Title           string `json:"title"`
	IP              string `json:"ip"`
	Country         string `json:"country"`
	Province        string `json:"province"`
	City            string `json:"city"`
	Operators       string `json:"operators"`
	Source          string `json:"source"`
	SourceKey       string `json:"sourcekey"`
}

const buffLength = 10

// LogChannel log channel
var LogChannel = make(chan ViewLog, buffLength)
var gifData, _ = base64.StdEncoding.DecodeString("R0lGODlhAQABAID/AP///wAAACwAAAAAAQABAAACAkQBADs=")

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
	var Country, Region, City, Isp string
	loc, err := ip17mon.Find(IP)
	if err == nil {
		Country, Region, City, Isp = loc.Country, loc.Region, loc.City, loc.Isp
		const NULL = "N/A"
		if Country == NULL {
			Country = ""
		}
		if Region == NULL {
			Region = ""
		}
		if City == NULL {
			City = ""
		}
		if Isp == NULL {
			Isp = ""
		}
		return Country, Region, City, Isp
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

func parseSource(url string, referer string) (string, string) {
	var source string
	var sourceKey string

	if referer == "" {
		return "direct", ""
	}

	refererURL, err := urllib.Parse(referer)
	if err != nil {
		return "", ""
	}
	urlRootDomain := getRootDomain(url)
	refererURLRootDomain := getRootDomain(referer)
	if refererURLRootDomain != "" {
		sourceKey = "referral-" + refererURLRootDomain
	}

	// 站内点击
	if urlRootDomain == refererURLRootDomain {
		source = "inner"
	}

	// 搜索引擎
	if value, ok := refererMap[refererURLRootDomain]; ok {
		source = value
	}

	// 微信相关
	queryFrom := refererURL.Query().Get("from")
	if value, ok := refererFromMap[queryFrom]; ok {
		source = value
	}

	// 外站引流
	if refererURLRootDomain != "" && source == "" {
		source = "referral"
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

func yield(r *http.Request) {
	LogChannel <- parseRequest(r)
}

func digest() {
	length := len(LogChannel)
	fmt.Println(length)
	if length == buffLength {
		for i := 0; i < length; i++ {
			viewlog := <-LogChannel
			jsonBody, _ := json.Marshal(viewlog)
			fmt.Println(string(jsonBody))
		}
	}
}

func parseRequest(r *http.Request) ViewLog {
	url := r.Header.Get("Referer")
	domain := r.Host
	userAgent := []byte(r.Header.Get("User-Agent"))
	query := r.URL.Query()
	header := r.Header
	created, _ := time.Time.MarshalText(time.Now())

	viewlog := ViewLog{}
	viewlog.Created = string(created)
	viewlog.URL = url
	viewlog.Domain = domain
	viewlog.Referer = query.Get("referer")
	viewlog.Cookieid = query.Get("cookieid")
	viewlog.Width = query.Get("width")
	viewlog.Height = query.Get("height")
	viewlog.Color = query.Get("color")
	viewlog.Language = query.Get("language")
	viewlog.Title = query.Get("title")
	viewlog.UserAgent = header.Get("User-Agent")
	viewlog.IsMobile, viewlog.IsWechat = parseMobile(userAgent)
	viewlog.Platform, viewlog.PlatformVersion = parsePlatform(userAgent)
	viewlog.Browser, viewlog.BrowserVersion = parseBrowser(userAgent)
	viewlog.IP = parseIP(r)
	viewlog.Country, viewlog.Province, viewlog.City, viewlog.Operators = parseIPAddress(viewlog.IP)
	viewlog.Source, viewlog.SourceKey = parseSource(viewlog.URL, viewlog.Referer)

	return viewlog
}

func handle(w http.ResponseWriter, r *http.Request) {
	debug := r.URL.Query().Get("debug")
	if debug == "" {
		go yield(r)
		go digest()
		w.Header().Set("Content-Type", "image/gif")
		w.Write(gifData)
	} else {
		viewlog := parseRequest(r)
		switch debug {
		case "mobile":
			io.WriteString(w, boolToString(viewlog.IsMobile))
		case "wechat":
			io.WriteString(w, boolToString(viewlog.IsWechat))
		case "platform":
			io.WriteString(w, viewlog.Platform)
		case "ip":
			io.WriteString(w, viewlog.IP)
		case "source":
			io.WriteString(w, viewlog.Source)
		default:
			io.WriteString(w, "")
		}
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
