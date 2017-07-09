package main

import (
	"log"
	"sync"

	"github.com/ashwanthkumar/golang-utils/sets"
)

// Item represents an object of the fetchedUrl, all outgoing urls we followed through from the page
// and all the assets (img, link & script) assets we found on the page
type UrlInfo struct {
	FetchedUrl   string
	OutgoingUrls sets.Set
	Assets       sets.Set
}

func NewUrlInfo(url string, links, assets []string) UrlInfo {
	return UrlInfo{
		FetchedUrl:   url,
		OutgoingUrls: sets.FromSlice(links),
		Assets:       sets.FromSlice(assets),
	}
}

type Listener chan UrlInfo

// SitemapManager represents a mapping of each url to it's corresponding outgoing
// urls and assets for the page
type SitemapManager struct {
	links    map[string]sets.Set
	assets   map[string]sets.Set
	listener Listener
	_running bool
	wait     sync.WaitGroup
}

// NewSitemapManager returns a new instance of SitemapManager
// that's listening on it's listener to add items
func NewSitemapManager(listenerBuffer int) *SitemapManager {
	s := SitemapManager{
		links:    make(map[string]sets.Set),
		assets:   make(map[string]sets.Set),
		listener: make(Listener, listenerBuffer),
		_running: true,
	}
	go s.start()
	return &s
}

// Listener returns a channel to which we can push Item structs to add it to Sitemap
func (s *SitemapManager) Listener() Listener {
	return s.listener
}

// Stop stops the listener channel and the update loop. You can't push any more items
// to this sitemap listener post this
func (s *SitemapManager) Stop() {
	s.wait.Wait()
	s._running = false
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
func (s *SitemapManager) AddInfo(urlInfo UrlInfo) {
	s.wait.Add(1)
	s.listener <- urlInfo
}

func (s *SitemapManager) start() {
	for s._running {
		select {
		case item := <-s.listener:
			s.wait.Done()
			log.Printf("[DEBUG] Updating sitemap for url=%s, links=%d, assets=%d\n", item.FetchedUrl, item.OutgoingUrls.Size(), item.Assets.Size())
			s.addLinks(item.FetchedUrl, item.OutgoingUrls)
			s.addAssets(item.FetchedUrl, item.Assets)
		}
	}
}

func (s *SitemapManager) addLinks(url string, links sets.Set) {
	s.links[url] = s.getOrEmptyLinks(url).Union(links)
}

func (s *SitemapManager) addAssets(url string, links sets.Set) {
	s.assets[url] = s.getOrEmptyLinks(url).Union(links)
}

func (s *SitemapManager) getOrEmptyLinks(url string) sets.Set {
	return s.getOrEmpty(url, s.links)
}

func (s *SitemapManager) getOrEmptyAssets(url string) sets.Set {
	return s.getOrEmpty(url, s.assets)
}

func (s *SitemapManager) getOrEmpty(url string, collection map[string]sets.Set) sets.Set {
	value, exists := collection[url]
	if exists {
		return value
	} else {
		return sets.Empty()
	}
}
