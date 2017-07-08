package main

import "testing"
import "github.com/stretchr/testify/assert"

func TestGetHostname(t *testing.T) {
	host, err := GetHostname("http://ashwanthkumar.in/resume.pdf")
	assert.NoError(t, err)
	assert.Equal(t, "ashwanthkumar.in", host)
}

func TestIsSameHost(t *testing.T) {
	assert.True(t, IsSameHostName("http://ashwanthkumar.in/resume.pdf", "ashwanthkumar.in"))
	assert.False(t, IsSameHostName("http://ashwanthkumar.in/resume.pdf", "www.ashwanthkumar.in"))
}
