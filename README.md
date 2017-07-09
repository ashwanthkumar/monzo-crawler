[![Build Status](https://travis-ci.org/ashwanthkumar/monzo-crawler.svg?branch=master)](https://travis-ci.org/ashwanthkumar/monzo-crawler)

# monzo-crawler

## Usage
```
$ make setup # uses glide to install all dependencies to vendor/
$ make build && ./monzo-crawler tomblomfield.com | tee sitemap.txt
$ less sitemap.txt
```

## Testing
```
$ make test
```

## Known Issues
- Doesn't have politeness delay implemented. We'll bombard the site with `runtime.NumCPU() * 4` concurrent HTTP requests
- HTML page parsing is done using [`goquery`](https://github.com/PuerkitoBio/goquery) library which is really slow for very big HTML pages (like that of amazon.com)

---

## Problem Statement
We'd like you to write a simple web crawler in a programming language of your choice. Feel free to either choose one you're very familiar with or, if you'd like to learn some Go, you can also make this your first Go program! The crawler should be limited to one domain - so when crawling tomblomfield.com it would crawl all pages within the domain, but not follow external links, for example to the Facebook and Twitter accounts. Given a URL, it should output a site map, showing which static assets each page depends on, and the links between pages.

Ideally, write it as you would a production piece of code. Bonus points for tests and making it as fast as possible!
