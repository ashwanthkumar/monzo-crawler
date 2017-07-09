package main

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/ashwanthkumar/golang-utils/sets"
)

// ExtractAllOutgoingUrls parses the given html document and returns all urls from the page
func ExtractAllOutgoingUrls(doc *goquery.Document, pageurl, targetHost string, crawledSoFar sets.Set) []string {
	CrawledLock.RLock()
	defer CrawledLock.RUnlock()
	allUrls := doc.Find("a").FilterFunction(func(i int, s *goquery.Selection) bool {
		value, exists := s.Attr("href")
		if exists && IsSameHostName(value, targetHost) {
			resolvedUrl := ResolveUrl(value, pageurl)
			return !crawledSoFar.Contains(resolvedUrl)
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
