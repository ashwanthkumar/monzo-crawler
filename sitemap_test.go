package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSitemapManager(t *testing.T) {
	s := NewSitemapManager(1)
	s.AddInfo(NewUrlInfo("url", []string{"link1", "link2"}, []string{"asset1", "asset2"}))
	s.Stop()

	info := s.InfoFor("url")
	assert.Equal(t, "url", info.FetchedUrl)
	assert.Equal(t, []string{"link1", "link2"}, info.OutgoingUrls.Values())
	assert.EqualValues(t, []string{"asset1", "asset2"}, info.Assets.Values())
}
