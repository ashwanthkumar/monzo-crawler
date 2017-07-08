package main

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/ashwanthkumar/golang-utils/sets"
	"github.com/ashwanthkumar/golang-utils/worker"
	"github.com/hashicorp/go-multierror"
	"github.com/parnurzeal/gorequest"
)

type Queue chan string

// ToCrawl is a queue of all pending urls that we have discovered and yet to crawl
var ToCrawl Queue

// Crawled contains list of all urls we've crawled so far.
// Might not be very efficient to store them all in memory, but a persistent store
// for this problem is an overkill at this point.
var Crawled sets.Set

// TargetHost contains the target host that we need to crawl
var TargetHost string

var workerPool worker.Pool

func main() {
	if len(os.Args) < 2 {
		log.Println("USAGE: ./monzo-crawler [tomblomfield.com]")
		os.Exit(1)
	}
	TargetHost := os.Args[1]
	ToCrawl = make(Queue, 512) // TODO - make this 512 configurable
	Crawled = sets.Empty()
	workerPool = worker.Pool{
		MaxWorkers: 2, // TODO - make this configurable
		Op:         Crawl,
	}
	workerPool.Initialize()

	log.Printf("Starting to crawl %s\n", TargetHost)
	ToCrawl <- "http://" + TargetHost
	running := true
	for running {
		select {
		case url := <-ToCrawl:
			log.Printf("Found url to add to work - %v\n", url)
			workerPool.AddWork(url)
		case <-time.Tick(30 * time.Second):
			if workerPool.Count() == 0 {
				running = false
			}
		}
	}
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

	Crawled.Add(url)

	log.Printf("Fetched the page - %v\n", resp.StatusCode)
	urls, err := ExtractAllOutgoingUrls(body, url)
	for _, u := range urls {
		ToCrawl <- u
	}

	return err
}

// ExtractAllOutgoingUrls parses the given html document and returns all urls from the page
func ExtractAllOutgoingUrls(body, pageurl string) ([]string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	if err != nil {
		return nil, err
	}

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

	return allUrls, nil
}

func combineErrors(errs []error) error {
	var err error
	for _, e := range errs {
		err = multierror.Append(err, e)
	}
	return err
}
