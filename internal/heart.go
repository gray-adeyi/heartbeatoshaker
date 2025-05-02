package internal

import (
	"errors"
	"net/http"
	"time"
)

const HttpTimeoutBeforeInterval = 200 * time.Millisecond

type Config struct {
	Targets []UrlTarget `yaml:"targets"`
}

type UrlTarget struct {
	Url      string        `yaml:"url"`
	Interval time.Duration `yaml:"interval"`
}

type Heart struct {
	targets []UrlTarget
	tickers map[string]*time.Ticker
	Done    chan struct{}
}

func NewHeart(targets []UrlTarget) *Heart {
	return &Heart{
		targets: targets,
		Done:    make(chan struct{}),
	}
}

func (p *Heart) Beat() error {
	if p.tickers != nil {
		return errors.New("pinger is already pinging some urls. stop the pinger first")
	}
	p.tickers = make(map[string]*time.Ticker, len(p.targets))
	for _, target := range p.targets {
		ticker := time.NewTicker(target.Interval)
		p.tickers[target.Url] = ticker
		go p.ping(target.Url, target.Interval)
	}
	return nil
}

func (p *Heart) Pause() {
	for _, ticker := range p.tickers {
		ticker.Stop()
	}
}

func (p *Heart) Resume() {
	for url, ticker := range p.tickers {
		interval := p.getIntervalFor(url)
		ticker.Reset(interval)
	}
}

func (p *Heart) Stop() {
	p.Pause()
	p.Done <- struct{}{}
}

func (p *Heart) ping(url string, interval time.Duration) {
	ticker := p.tickers[url]
	for {
		_, ok := <-ticker.C
		if !ok {
			break
		}
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			continue
		}
		client := http.Client{
			Timeout: interval - HttpTimeoutBeforeInterval,
		}
		client.Do(req)
	}
}

func (p *Heart) getIntervalFor(url string) time.Duration {
	for _, target := range p.targets {
		if target.Url == url {
			return target.Interval
		}
	}
	return 0
}
