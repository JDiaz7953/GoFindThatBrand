package main

import (
	"net/http"
	"strings"
)

func requestHandler(url, key string) (*http.Request, error) {

	reqBody := strings.NewReader(
		`{"requests":[` +
			`{"indexName":"Listing_by_listing_quality_production","params":"` +
			`analytics=true&clickAnalytics=true&enableABTest=true` +
			`&facets=%5B%22badges%22%2C%22category_path%22%2C%22category_size%22%2C%22color%22%2C%22condition%22%2C%22designers.name%22%2C%22location%22%2C%22price_i%22%2C%22strata%22%5D` +
			`&filters=department%3Amenswear&getRankingInfo=true` +
			`&highlightPostTag=__%2Fais-highlight__&highlightPreTag=__ais-highlight__` +
			`&hitsPerPage=40&maxValuesPerFacet=165` +
			`&numericFilters=%5B%22price_i%3E%3D0%22%2C%22price_i%3C%3D1000000%22%5D` +
			`&page=1&query=&userToken=fc1d7cf6-5b6b-4d0e-86e5-7c2e8388b325"},` +
			`{"indexName":"Listing_by_listing_quality_production","params":"` +
			`analytics=false&clickAnalytics=false&enableABTest=true&facets=price_i` +
			`&filters=department%3Amenswear&getRankingInfo=true` +
			`&highlightPostTag=__%2Fais-highlight__&highlightPreTag=__ais-highlight__` +
			`&hitsPerPage=0&maxValuesPerFacet=165&page=0&query=` +
			`&userToken=fc1d7cf6-5b6b-4d0e-86e5-7c2e8388b325"}` +
			`]}`)

	req, _ := http.NewRequest("POST", url, reqBody)
	req.Header.Set("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 18_7 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/26.6.1 Mobile/15E148 Safari/604.1")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Origin", "https://www.grailed.com")
	req.Header.Set("Referer", "https://www.grailed.com/")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("x-algolia-application-id", "MNRWEFSS2Q")
	req.Header.Set("x-algolia-api-key", key)

	return req, nil

}
