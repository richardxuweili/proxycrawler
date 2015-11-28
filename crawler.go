package proxycrawler

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/xlaurent/proxycrawler/proxy"
	"github.com/xlaurent/proxycrawler/source"
)

type Crawler struct {
	clients []*http.Client
	sources []source.Source
}

func New(sources []source.Source, timeout time.Duration) *Crawler {
	clients := make([]*http.Client, len(sources))
	for i := range clients {
		clients[i] = &http.Client{Timeout: timeout}
	}
	return &Crawler{
		clients: clients,
		sources: sources,
	}
}

func (p *Crawler) FetchProxys(
	testURL string,
	max int,
	check func(resp *http.Response) error,
	cancelCh <-chan struct{},
) chan *proxy.Proxy {
	recvCh := make(chan *proxy.Proxy, max)
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
			for _, proxy := range proxys {
				select {
				case <-cancelCh:
					return
				default:
				}
				if err := proxy.Test(p.clients[i], testURL, check); err != nil {
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
		}(i)
	}
	go func() {
		waiter.Wait()
		close(recvCh)
	}()
	return recvCh
}
