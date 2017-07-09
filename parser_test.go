package main

import (
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/ashwanthkumar/golang-utils/io"
	"github.com/ashwanthkumar/golang-utils/sets"
	"github.com/stretchr/testify/assert"
)

func TestExtractAllOutgoingUrls(t *testing.T) {
	doc := parseDocument(t, "tests/fixtures/tomblomfield_home.html")
	outgoingUrls := ExtractAllOutgoingUrls(doc, "http://tomblomfield.com/", "tomblomfield.com", sets.Empty())
	assert.Len(t, outgoingUrls, 22)

	doc = parseDocument(t, "tests/fixtures/a9_echo.html")
	outgoingUrls = ExtractAllOutgoingUrls(doc, "https://www.amazon.com/dp/B01DFKC2SO?psc=1", "www.amazon.com", sets.Empty())
	assert.Len(t, outgoingUrls, 100)
}

func TestExtractAssetsOnPage(t *testing.T) {
	doc := parseDocument(t, "tests/fixtures/tomblomfield_home.html")
	assetUrls := ExtractAllAssetsOnPage(doc, "http://tomblomfield.com/")
	assert.Len(t, assetUrls, 15)

	doc = parseDocument(t, "tests/fixtures/a9_echo.html")
	assetUrls = ExtractAllAssetsOnPage(doc, "https://www.amazon.com/dp/B01DFKC2SO?psc=1")
	assert.Len(t, assetUrls, 108)
}

func parseDocument(t *testing.T, path string) *goquery.Document {
	body, err := io.ReadFullyFromFile(path)
	assert.NoError(t, err)

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	assert.NoError(t, err)

	return doc
}
