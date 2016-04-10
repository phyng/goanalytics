package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"testing"
)

var testServer = "http://localhost:8001"

func errorHandle(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func testHeader(header map[string]string, url string, debug string, except string) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	errorHandle(err)

	req.Header.Set("Referer", "http://www.example.com")
	for key, value := range header {
		req.Header.Set(key, value)
	}

	resp, err := client.Do(req)
	errorHandle(err)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	errorHandle(err)
	if string(body) != except {
		log.Fatalf(`ERR %s: "%s" is not "%s"`, debug, body, except)
	} else {
		log.Printf("OK %s %s", debug, except)
	}
}

func testUserAgent(UserAgent string, debug string, except string) {
	header := map[string]string{
		"User-Agent": UserAgent,
	}
	url := testServer + "/?debug=" + debug
	testHeader(header, url, debug, except)
}

func testXForwardedFor(XForwardedFor string, debug string, except string) {
	header := map[string]string{
		"X-Forwarded-For": XForwardedFor,
	}
	url := testServer + "/?debug=" + debug
	testHeader(header, url, debug, except)
}

func testSource(referer string, debug string, except string) {
	header := map[string]string{
		"Host":    "www.localhost.com",
		"Referer": "http://www.localhost.com",
	}
	url := testServer + "/?referer=" + referer + "&debug=" + debug
	testHeader(header, url, debug, except)
}

// TestOld old tests
func TestOld(*testing.T) {
	testUserAgent("mobile", "mobile", "true")
	testUserAgent("desktop", "mobile", "false")
	testUserAgent("MicroMessenger", "wechat", "true")
	testUserAgent("Line", "wechat", "false")

	testUserAgent("Mozilla/5.0 (compatible; MSIE 10.0; Windows NT 6.2; Trident/6.0; Xbox; Xbox One)", "platform", "Windows")
	testUserAgent("Mozilla/5.0 (Windows NT 10.0; WOW64; Trident/7.0; rv:11.0) like Gecko", "platform", "Windows")
	testUserAgent("Mozilla/5.0 (Linux; Android 4.0.4; Galaxy Nexus Build/IMM76B) AppleWebKit/535.19 (KHTML, like Gecko) Chrome/18.0.1025.133 Mobile Safari/535.19", "platform", "Android")
	testUserAgent("Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2693.2 Safari/537.36", "platform", "Linux")
	testUserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_9_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/35.0.1916.114 Safari/537.36", "platform", "Mac")
	testUserAgent("Mozilla/5.0 (iPhone; CPU iPhone OS 9_1 like Mac OS X) AppleWebKit/601.1.46 (KHTML, like Gecko) Version/9.0 Mobile/13B143 Safari/601.1", "platform", "iPhone")

	testXForwardedFor("", "ip", "127.0.0.1")
	testXForwardedFor("8.8.8.8", "ip", "8.8.8.8")
	testXForwardedFor("8.8.8.8, 114.114.114.114", "ip", "114.114.114.114")

	testSource("https://www.google.com", "source", "google")
	testSource("https://www.baidu.com", "source", "baidu")
	testSource("https://www.bing.com", "source", "bing")
	testSource("https://www.sogou.com", "source", "sogou")
	testSource("http://www.localhost.com", "source", "inner")
	testSource("http://www.other.com", "source", "referral")
}
