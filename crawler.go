package proxy

import (
	"log"
	"net/http"
	"sync"
	"time"
)

type Source func() ([]*Proxy, error)

type Crawler struct {
	clients []*http.Client
	sources []Source
}

func NewCrawler(sources []Source, timeout time.Duration, concurrency int) *Crawler {
	clients := make([]*http.Client, len(sources)*concurrency)
	for i := range clients {
		clients[i] = &http.Client{Timeout: timeout}
	}
	return &Crawler{
		clients: clients,
		sources: sources,
	}
}

func (p *Crawler) FetchProxys(
	recvCh chan<- *Proxy,
	testURL string,
	check func(resp *http.Response) error,
	cancelCh <-chan struct{},
) {
	waiter := &sync.WaitGroup{}
	for i := range p.sources {
		waiter.Add(1)
		go func(i int) {
			defer waiter.Done()
			proxys, err := p.sources[i]()
			if err != nil {
				log.Printf("fetch proxys function %d :%v", i, err)
				return
			}
			shuffleSlice(proxys)
			concurrency := len(p.clients) / len(p.sources)
			eachLen := len(proxys) / concurrency
			if eachLen == 0 {
				eachLen = len(proxys)
				concurrency = 1
			}
			for j := 0; j < concurrency; j++ {
				end := (j + 1) * eachLen
				if j == concurrency-1 {
					end = len(proxys)
				}
				waiter.Add(1)
				go func(subProxys []*Proxy, client *http.Client) {
					defer waiter.Done()
					for _, proxy := range subProxys {
						if err := proxy.Test(client, testURL, check); err != nil {
							continue
						}
						select {
						case recvCh <- proxy:
						case <-cancelCh:
							return
						default:
							log.Printf("because the channel is full,stop fetching source %d", i)
							return
						}
					}
				}(proxys[j*eachLen:end], p.clients[i*concurrency+j])
			}
		}(i)
	}
	waiter.Wait()
	close(recvCh)
}
