package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/xlaurent/proxy"
)

var timeout = flag.Duration("t", 2*time.Second, "set timeout")
var cap = flag.Int("m", 50, "the max number of proxys can be fetched")
var testURL = flag.String("url", "http://www.douban.com", "test URL")

func main() {
	crawler := proxy.NewCrawler([]proxy.Source{proxy.CyberSource}, *timeout, 3)
	recvCh := make(chan *proxy.Proxy, *cap)
	go crawler.FetchProxys(recvCh, *testURL, proxy.DefaultCheck, nil)
	for p := range recvCh {
		fmt.Println(p.String(), " ", p.ConnTime)
	}
}
