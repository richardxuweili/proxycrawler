package source

import (
	"bufio"
	"io"
	"regexp"

	"github.com/xlaurent/proxycrawler/proxy"
)

type Source func() ([]*proxy.Proxy, error)

var proxyRegex = regexp.MustCompile(`(^(\w*?)://|^)(.+?):(\d+)`)

func Reader(rd io.Reader) func() ([]*proxy.Proxy, error) {
	return func() ([]*proxy.Proxy, error) {
		scanner := bufio.NewScanner(rd)
		var proxys []*proxy.Proxy
		for scanner.Scan() {
			str := scanner.Text()
			match := proxyRegex.FindStringSubmatch(str)
			if len(match) != 5 {
				continue
			}
			proxy := &proxy.Proxy{IP: match[3], Port: match[4]}

			if proxy.Scheme = match[2]; proxy.Scheme == "" {
				proxy.Scheme = "http"
			}
			proxys = append(proxys, proxy)
		}
		return proxys, nil
	}
}
