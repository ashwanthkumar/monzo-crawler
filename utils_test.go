package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetHostname(t *testing.T) {
	host, err := GetHostname("http://ashwanthkumar.in/resume.pdf")
	assert.NoError(t, err)
	assert.Equal(t, "ashwanthkumar.in", host)
}

func TestIsSameHost(t *testing.T) {
	assert.True(t, IsSameHostName("http://ashwanthkumar.in/resume.pdf", "ashwanthkumar.in"))
	assert.False(t, IsSameHostName("http://ashwanthkumar.in/resume.pdf", "www.ashwanthkumar.in"))
}

func TestResolveUrl(t *testing.T) {
	assert.Equal(t, "http://ashwanthkumar.in/", ResolveUrl("/", "http://ashwanthkumar.in/about"))
	assert.Equal(t, "http://ashwanthkumar.in/archive", ResolveUrl("/archive", "http://ashwanthkumar.in/about"))
}

func TestToUrl(t *testing.T) {
	assert.Equal(t, "http://ashwanthkumar.in", DomainToUrl("ashwanthkumar.in"))
}
