package proxy

import (
	"bufio"
	"io"
	"regexp"
)

var proxyRegex = regexp.MustCompile(`(^(\w*?)://|^)(.+?):(\d+)`)

func ReaderSource(rd io.Reader) ([]*Proxy, error) {
	scanner := bufio.NewScanner(rd)
	var proxys []*Proxy
	for scanner.Scan() {
		str := scanner.Text()
		match := proxyRegex.FindStringSubmatch(str)
		if len(match) != 5 {
			continue
		}
		proxy := &Proxy{IP: match[3], Port: match[4]}
		proxy.Scheme = match[2]
		if match[2] == "" {
			proxy.Scheme = "http"
		}
		proxys = append(proxys, proxy)
	}
	return proxys, nil
}
