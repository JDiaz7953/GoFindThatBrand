package main

import (
	"encoding/json"
	_ "encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

/*

1. How do we scrape the data?

2. What Packages and struct models need to be created

3. Where is this data being sent to? (discord, database)?

4. is bypassing anti bots procedures possible?

*/

/*
	------ Classify as finished When the Following is done -------

	- Does the webscraper process without the need of human intervention

	- Does the webscraper only returns clothing from a specified price range

	- Does the webscraper only returns clothing from a specified brand

	- Does the webscraper only returns clothing for a specifed size

	- Does the webscraper display readable/link to view/buy specifed clothing

*/

// Focus Point: How do we scrape data from depop using go?
// A: Depop makes the listing of the clothing public after hitting 24 items, grab the url and use that and we are simply using the json to scrape through all the stuff

/*
TO DO:
- Refactor stucts to take the data we want
- begin building the program to get specifed brands and clothing type
- get only a specifed brand
- create function to get a new api key


*/

type Root struct {
	Results []Results `json:"results,omitempty"`
}

type Hits struct {
	Badges       []interface{} `json:"badges,omitempty"`
	BumpedAt     string        `json:"bumped_at,omitempty"`
	DesignerName string        `json:"designer_names"`
	Sold         bool          `json:"sold"`
}

type Results struct {
	Hits []Hits `json:"hits,omitempty"`
}

func readKey(path string) (string, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return "", nil
	}

	return strings.TrimSpace(string(b)), nil
}

func main() {
	key, err := readKey(`C:\Users\jaide\Documents\Workspace\Projects\GoFindThatBrand\key.txt`)
	if err != nil{
		
		log.Fatal(err)
	}
	url := "https://mnrwefss2q-dsn.algolia.net/1/indexes/*/queries"

	// research what newReader does exactly
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

	//header.set is explicit in the header that get put in the body
	req, _ := http.NewRequest("POST", url, reqBody)
	req.Header.Set("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 18_7 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/26.6.1 Mobile/15E148 Safari/604.1")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Origin", "https://www.grailed.com")
	req.Header.Set("Referer", "https://www.grailed.com/")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("x-algolia-application-id", "MNRWEFSS2Q")
	req.Header.Set("x-algolia-api-key", key)

	//Perform the request
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Fatal("Unexpected Status Code ", resp.StatusCode)
	}
	defer resp.Body.Close()

	resBody, err := io.ReadAll(resp.Body)

	//After getting readable data it needs to me unmarshaled pointed to the result struct
	var r Root
	err = json.Unmarshal(resBody, &r)
	if err != nil {
		log.Fatal("Error extracting data")
	}

	data := r.Results[0]
	//fmt.Println(r)
	for _, v := range data.Hits {
		if v.DesignerName == "Supreme" {
			fmt.Println(v)
		}
	}

	//fmt.Print(string(resBody))

}
