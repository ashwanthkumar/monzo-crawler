package main

import "net/url"

func GetHostname(rawurl string) (string, error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return "", err
	} else {
		return u.Hostname(), nil
	}
}

func IsSameHostName(rawurl, targetHost string) bool {
	host, err := GetHostname(rawurl)
	if err != nil {
		return false
	} else {
		return host == targetHost
	}
}
