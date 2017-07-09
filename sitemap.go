package main

import (
	"log"

	"github.com/ashwanthkumar/golang-utils/sets"
	"github.com/ashwanthkumar/golang-utils/sync"
)

// Item represents an object of the fetchedUrl, all outgoing urls we followed through from the page
// and all the assets (img, link & script) assets we found on the page
type UrlInfo struct {
	FetchedUrl   string
	OutgoingUrls sets.Set
	Assets       sets.Set
}

func NewUrlInfo(url string, links, assets []string) *UrlInfo {
	return &UrlInfo{
		FetchedUrl:   url,
		OutgoingUrls: sets.FromSlice(links),
		Assets:       sets.FromSlice(assets),
	}
}

type Listener chan *UrlInfo

// SitemapManager represents a mapping of each url to it's corresponding outgoing
// urls and assets for the page
type SitemapManager struct {
	links    map[string]sets.Set
	assets   map[string]sets.Set
	listener Listener

	// fields used for sync
	stop chan bool
	wait sync.CountWG
}

// NewSitemapManager returns a new instance of SitemapManager
// that's listening on it's listener to add items
func NewSitemapManager(listenerBuffer int) *SitemapManager {
	s := SitemapManager{
		links:    make(map[string]sets.Set),
		assets:   make(map[string]sets.Set),
		stop:     make(chan bool),
		listener: make(Listener, listenerBuffer),
	}
	go s.start()
	return &s
}

// Stop stops the listener channel and the update loop. You can't push any more items
// to this sitemap listener post this
func (s *SitemapManager) Stop() {
	s.wait.Wait()
	close(s.stop)
	close(s.listener)
}

// InfoFor returns UrlInfo associated with the url if it's found
// else it would still return a valid UrlInfo but with assets and links as empty
func (s *SitemapManager) InfoFor(url string) *UrlInfo {
	links := s.getOrEmpty(url, s.links)
	assets := s.getOrEmpty(url, s.assets)
	return &UrlInfo{
		FetchedUrl:   url,
		OutgoingUrls: links,
		Assets:       assets,
	}
}

// AddInfo adds a particular UrlInfo to the sitemap
func (s *SitemapManager) AddInfo(urlInfo *UrlInfo) {
	s.wait.Add(1)
	s.listener <- urlInfo
}

func (s *SitemapManager) start() {
	running := true
	for running {
		select {
		case item := <-s.listener:
			if nil != item {
				log.Printf("[DEBUG] Updating sitemap for url=%s, links=%d, assets=%d\n", item.FetchedUrl, item.OutgoingUrls.Size(), item.Assets.Size())
				s.addLinks(item.FetchedUrl, item.OutgoingUrls)
				s.addAssets(item.FetchedUrl, item.Assets)
				s.wait.Done()
			}
		case <-s.stop:
			running = false
		}
	}
}

func (s *SitemapManager) addLinks(url string, links sets.Set) {
	s.links[url] = s.getOrEmpty(url, s.links).Union(links)
}

func (s *SitemapManager) addAssets(url string, links sets.Set) {
	s.assets[url] = s.getOrEmpty(url, s.assets).Union(links)
}

func (s *SitemapManager) getOrEmpty(url string, collection map[string]sets.Set) sets.Set {
	value, exists := collection[url]
	if exists {
		return value
	} else {
		return sets.Empty()
	}
}
