package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/xlaurent/proxy"
)

var config = struct {
	ConfigFile string
	TestURL    string
	Interval   time.Duration
	Cap        int
}{}

func init() {
	SetFlag()
	ParseFlag()
}

func main() {
	cache := NewProxysCache()
	c := crawler{
		Crawler:   proxy.NewCrawler([]proxy.Source{proxy.CyberSource}, 2*time.Second, 3),
		Observers: []Observer{cache},
		TestURL:   config.TestURL,
		Cap:       config.Cap,
		Check:     proxy.DefaultCheck,
	}

	go func() {
		ticker := time.Tick(config.Interval)
		for {
			cancelCh := make(chan struct{})
			go c.Crawl(cancelCh)
			<-ticker
			close(cancelCh)
		}
	}()

	http.HandleFunc("/proxys", CacheHandler(cache))

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func SetFlag() {
	flag.StringVar(&config.ConfigFile, "c", "", "set config path")
	flag.StringVar(&config.TestURL, "url", "http://www.douban.com", "use it to test proxy")
	flag.DurationVar(&config.Interval, "t", 1*time.Minute, "flush cache")
	flag.IntVar(&config.Cap, "max", 100, "max number of proxy")
}

func ParseFlag() {
	flag.Parse()
	if config.ConfigFile == "" {
		return
	}
	f, err := os.Open(config.ConfigFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer f.Close()
	if err := json.NewDecoder(f).Decode(&config); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return
}

func CacheHandler(cache *proxysCache) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(cache.Read())
	}
}
