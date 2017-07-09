package main

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSitemapManager(t *testing.T) {
	s := NewSitemapManager(1)
	s.AddInfo(NewUrlInfo("url", []string{"link1", "link2"}, []string{"asset1", "asset2"}))
	s.Stop()

	info := s.InfoFor("url")
	assert.Equal(t, "url", info.FetchedUrl)
	outgoingUrls := info.OutgoingUrls.Values()
	sort.Strings(outgoingUrls)
	assets := info.Assets.Values()
	sort.Strings(assets)
	assert.Equal(t, []string{"link1", "link2"}, outgoingUrls)
	assert.EqualValues(t, []string{"asset1", "asset2"}, assets)
}
