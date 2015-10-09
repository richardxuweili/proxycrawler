package main

import (
	"net/http"
	"sync"
	"time"

	"github.com/xlaurent/proxy"
)

type Observer interface {
	Update([]*proxy.Proxy)
}

type crawler struct {
	Proxys    []*proxy.Proxy
	Observers []Observer
	*proxy.Crawler
	TestURL string
	Cap     int
	Check   func(*http.Response) error
}

type proxysCache struct {
	Total      int
	Proxys     []*proxy.Proxy
	LastUpdate time.Time
	*sync.RWMutex
}

func (c crawler) Crawl(cancelCh chan struct{}) {
	recvCh := make(chan *proxy.Proxy, c.Cap)
	go c.FetchProxys(recvCh, c.TestURL, c.Check, cancelCh)
	for proxy := range recvCh {
		c.Proxys = append(c.Proxys, proxy)
	}
	c.UpdateObservers()
}

func (c crawler) UpdateObservers() {
	for _, ob := range c.Observers {
		ob.Update(c.Proxys)
	}
}

func NewProxysCache() *proxysCache {
	return &proxysCache{
		Total:      0,
		LastUpdate: time.Now(),
		RWMutex:    &sync.RWMutex{},
	}
}

func (pc *proxysCache) Update(proxys []*proxy.Proxy) {
	pc.Lock()
	defer pc.Unlock()
	pc.Proxys = proxys
	pc.Total = len(pc.Proxys)
	pc.LastUpdate = time.Now()
}

func (pc *proxysCache) Read() *proxysCache {
	pc.RLock()
	proxys := pc
	pc.RUnlock()
	return proxys
}
