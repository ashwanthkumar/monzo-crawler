package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/ashwanthkumar/golang-utils/sets"
	"github.com/ashwanthkumar/golang-utils/worker"
	"github.com/hashicorp/go-multierror"
	"github.com/parnurzeal/gorequest"
)

type Queue chan string

// BUFFER_SIZE represents the buffered channel sizes we use in the crawler
const BUFFER_SIZE = 512

// MAX_FETCHERS represents the maximum number of concurrent fetchers that can run parallely.
// We do 4 fetches per core on the machine
var MAX_FETCHERS = runtime.NumCPU() * 4

// ToCrawl is a queue of all pending urls that we have discovered and yet to crawl
var ToCrawl Queue

// Crawled contains list of all urls we've crawled so far.
// Might not be very efficient to store them all in memory, but a persistent store
// for this problem is an overkill at this point.
var Crawled sets.Set
var CrawledLock sync.RWMutex

// TargetHost contains the target host that we need to crawl
var TargetHost string

var sitemapManager *SitemapManager
var workerPool worker.Pool

func main() {
	if len(os.Args) < 2 {
		log.Println("USAGE: ./monzo-crawler [tomblomfield.com]")
		os.Exit(1)
	}
	TargetHost := os.Args[1]
	ToCrawl = make(Queue, BUFFER_SIZE)
	Crawled = sets.Empty()
	workerPool = worker.Pool{
		MaxWorkers: MAX_FETCHERS,
		Op:         Crawl,
	}
	workerPool.Initialize()
	sitemapManager = NewSitemapManager(BUFFER_SIZE)

	log.Printf("Starting to crawl %s\n", TargetHost)
	homePage := "http://" + TargetHost
	ToCrawl <- homePage
	running := true
	for running {
		select {
		case url := <-ToCrawl:
			log.Printf("Enqueuing url=%v\n", url)
			workerPool.AddWork(url)
		case <-time.Tick(5 * time.Second): // TODO - Make this configurable
			log.Printf("[DEBUG] in_flight_urls=%d\n", workerPool.Count())
			log.Printf("[DEBUG] pending_urls=%d\n", len(ToCrawl))
			if len(ToCrawl)+workerPool.Count() == 0 {
				log.Println("No Active / Pending work left")
				running = false
			}
		}
	}

	workerPool.Wait()
	sitemapManager.Stop()
	PrintSitemap(homePage)
	log.Println("Good bye!")
}

// Crawl does the actual fetching of the page
func Crawl(req worker.Request) error {
	url := req.(string)
	log.Printf("Fetching %v\n", url)
	resp, body, errs := gorequest.New().Get(url).End()
	err := combineErrors(errs)

	if err != nil {
		log.Printf("%v\n", err)
		return err
	}

	CrawledLock.Lock()
	Crawled.Add(url)
	CrawledLock.Unlock()

	log.Printf("Fetched url=%s statusCode=%d\n", url, resp.StatusCode)
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	if err != nil {
		return err
	}

	links := ExtractAllOutgoingUrls(doc, url)
	for _, u := range links {
		ToCrawl <- u
	}
	assets := ExtractAllAssetsOnPage(doc, url)
	sitemapManager.AddInfo(NewUrlInfo(url, links, assets))

	return err
}

// ExtractAllOutgoingUrls parses the given html document and returns all urls from the page
func ExtractAllOutgoingUrls(doc *goquery.Document, pageurl string) []string {
	CrawledLock.RLock()
	defer CrawledLock.RUnlock()
	allUrls := doc.Find("a").FilterFunction(func(i int, s *goquery.Selection) bool {
		value, exists := s.Attr("href")
		if exists && IsSameHostName(value, TargetHost) {
			resolvedUrl := ResolveUrl(value, pageurl)
			return !Crawled.Contains(resolvedUrl)
		}
		return false
	}).Map(func(i int, s *goquery.Selection) string {
		value, _ := s.Attr("href")
		resolvedUrl := ResolveUrl(value, pageurl)
		return resolvedUrl
	})

	return allUrls
}

// ExtractAllAssetsOnPage parses the given html document and returns all urls of the
// assets on the page - img, link, script etc.
func ExtractAllAssetsOnPage(doc *goquery.Document, pageurl string) []string {
	assetUrls := doc.Find("img,link,script").FilterFunction(func(i int, s *goquery.Selection) bool {
		_, srcExists := s.Attr("src")
		_, hrefExists := s.Attr("href")
		return srcExists || hrefExists
	}).Map(func(i int, s *goquery.Selection) string {
		srcValue, exists := s.Attr("src")
		var value string
		if exists {
			value = srcValue
		} else {
			hrefValue, _ := s.Attr("href")
			value = hrefValue
		}
		resolvedUrl := ResolveUrl(value, pageurl)
		return resolvedUrl
	})

	return assetUrls
}

func PrintSitemap(url string) {
	fmt.Println(".")
	_printSitemap(url, sets.Empty(), 0)
}

func _printSitemap(url string, seenSoFar sets.Set, depth int) {
	fmt.Printf("%s└── %s\n", strings.Repeat(" ", depth), url)
	seenSoFar.Add(url)
	info := sitemapManager.InfoFor(url)
	if info.Assets.Size() > 0 {
		fmt.Printf("%s│%s└── ASSETS \n", strings.Repeat(" ", depth+2), strings.Repeat(" ", depth+2))
		for idx, asset := range info.Assets.Values() {
			marker := "├"
			if idx == info.Assets.Size()-1 {
				marker = "└"
			}
			fmt.Printf("%s│%s%s── %s\n", strings.Repeat(" ", depth+2), strings.Repeat(" ", depth+4), marker, asset)
		}
	}

	if info.OutgoingUrls.Size() > 0 {
		for _, link := range info.OutgoingUrls.Values() {
			if !seenSoFar.Contains(link) {
				seenSoFar.Add(link)
				_printSitemap(link, seenSoFar, depth+2)
			}
		}
	}
}

func combineErrors(errs []error) error {
	var err error
	for _, e := range errs {
		err = multierror.Append(err, e)
	}
	return err
}
