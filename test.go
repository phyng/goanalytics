package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func errorHandle(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func testUserAgent(UserAgent string, debug string, except string) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://localhost:"+os.Args[1]+"/?debug="+debug, nil)
	errorHandle(err)
	req.Header.Set("User-Agent", UserAgent)
	resp, err := client.Do(req)
	errorHandle(err)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	errorHandle(err)
	if string(body) != except {
		log.Fatalln(except)
	} else {
		log.Printf("OK %s %s %s", debug, except, UserAgent)
	}
}

func main() {
	testUserAgent("mobile", "mobile", "true")
	testUserAgent("desktop", "mobile", "false")

	testUserAgent("MicroMessenger", "wechat", "true")
	testUserAgent("Line", "wechat", "false")

	testUserAgent("Mozilla/5.0 (Windows NT 6.3; Trident/7.0; rv:11.0) like Gecko", "platform", "Windows")
	testUserAgent("Mozilla/5.0 (Linux; U; Android 4.0.1; ja-jp; Galaxy Nexus……", "platform", "Android")
	testUserAgent("Mozilla/5.0 (X11; Linux x86_64) ...", "platform", "Linux")
	testUserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_9_3) ...", "platform", "Mac")
}
