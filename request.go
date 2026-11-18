package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func requestHandler(urlToSite, key string, page, currentLowPrice, currentHighPrice int) Results {

	priceRange := url.QueryEscape(fmt.Sprintf(`["price_i>=%d","price_i<=%d"]`, currentLowPrice, currentHighPrice))

	reqBody := strings.NewReader(
		`{"requests":[` +
			`{"indexName":"Listing_by_listing_quality_production","params":"` +
			`analytics=true&clickAnalytics=true&enableABTest=true` +
			`&facets=%5B%22badges%22%2C%22category_path%22%2C%22category_size%22%2C%22color%22%2C%22condition%22%2C%22designers.name%22%2C%22location%22%2C%22price_i%22%2C%22strata%22%5D` +
			`&filters=department%3Amenswear&getRankingInfo=true` +
			`&highlightPostTag=__%2Fais-highlight__&highlightPreTag=__ais-highlight__` +
			`&hitsPerPage=40&maxValuesPerFacet=165` +
			`&numericFilters=` + priceRange +
			`&page=` + strconv.Itoa(page) +
			`&query=&userToken=fc1d7cf6-5b6b-4d0e-86e5-7c2e8388b325"},` +
			`{"indexName":"Listing_by_listing_quality_production","params":"` +
			`analytics=false&clickAnalytics=false&enableABTest=true&facets=price_i` +
			`&filters=department%3Amenswear&getRankingInfo=true` +
			`&highlightPostTag=__%2Fais-highlight__&highlightPreTag=__ais-highlight__` +
			`&hitsPerPage=0&maxValuesPerFacet=165&page=0&query=` +
			`&userToken=fc1d7cf6-5b6b-4d0e-86e5-7c2e8388b325"}` +
			`]}`)

	req, _ := http.NewRequest("POST", urlToSite, reqBody)
	req.Header.Set("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 18_7 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/26.6.1 Mobile/15E148 Safari/604.1")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Origin", "https://www.grailed.com")
	req.Header.Set("Referer", "https://www.grailed.com/")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("x-algolia-application-id", "MNRWEFSS2Q")
	req.Header.Set("x-algolia-api-key", key)

	//Performs the request to the site
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal("Unexpected Status Code ", resp.StatusCode)
	}
	defer resp.Body.Close()

	//After getting readable data it needs to me unmarshaled pointed to the result struct to map to the json tags of the models
	var root Root
	resBody, err := io.ReadAll(resp.Body)
	err = json.Unmarshal(resBody, &root)
	if err != nil {
		log.Fatal("Error extracting data")
	}

	return root.Results[0]
}
