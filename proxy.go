package proxy

import (
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type Proxy struct {
	Scheme   string
	IP       string
	Port     string
	ConnTime time.Duration
}

func (p *Proxy) Test(client *http.Client, URL string, check func(resp *http.Response) error) error {
	transport, err := p.Transport()
	if err != nil {
		return err
	}
	client.Transport = transport
	before := time.Now()
	resp, err := client.Get(URL)
	connTime := time.Now().Sub(before)
	if err != nil {
		return err
	}
	p.ConnTime = connTime
	if check == nil {
		return nil
	}
	err = check(resp)
	if err != nil {
		return err
	}
	return nil
}

func (p *Proxy) Transport() (*http.Transport, error) {
	URL, err := url.Parse(p.String())
	if err != nil {
		return nil, fmt.Errorf("can't parse proxy url %s,%v", p.String(), err)
	}
	return &http.Transport{Proxy: http.ProxyURL(URL)}, nil
}

func (p Proxy) String() string {
	return fmt.Sprintf("%s://%s:%s", p.Scheme, p.IP, p.Port)
}
