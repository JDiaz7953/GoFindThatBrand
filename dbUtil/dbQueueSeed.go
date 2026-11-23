package dbUtil

import (
	"ClotheScrape/db"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"ClotheScrape/models"
)

func SeedQueue(ctx context.Context, queries *db.Queries) {
	body := strings.NewReader(`{"requests":[{"indexName":"Listing_by_listing_quality_production","params":"` +
		`facets=%5B%22designers.name%22%2C%22category_path%22%5D` +
		`&filters=department%3Amenswear&hitsPerPage=0&maxValuesPerFacet=1000&page=0&query="}]}`)

	key := "MGMyZmIyYmE2OWM5MGNmMDBmMTVjYzgyMzc0NjAzODg4OWVlZTBlMWZlOWYwNWNjNGRkNzIyN2YxOTQzNDJlN3ZhbGlkVW50aWw9MTc5NTQyMzY3OSZ1c2VyVG9rZW49OTYuMjM4LjYzLjE5"

	req, _ := http.NewRequest("POST", "https://mnrwefss2q-dsn.algolia.net/1/indexes/*/queries", body)
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Origin", "https://www.grailed.com")
	req.Header.Set("Referer", "https://www.grailed.com/")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("x-algolia-application-id", "MNRWEFSS2Q")
	req.Header.Set("x-algolia-api-key", key)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal("Failure to make request to the site", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Fatal("Request failed with status", resp.StatusCode)
	}

	var root models.Root
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("error filling root table with data")
	}

	err = json.Unmarshal(respBody, &root)
	if err != nil {
		log.Fatal("Failure to insert Queue Seed", err)
	}

	// Conditional that filters for the category and brand I want, looping through a map
	var designerNames []string
	var categories []string

	if len(root.Results) == 0 {
		log.Fatal("No data available goodbye...")
	}
	facets := root.Results[0].Facets

	for designerName := range facets.DesignerNames {
		designerNames = append(designerNames, designerName)
	}
	for category := range facets.CategoryPath {
		categories = append(categories, category)
	}

	fmt.Println(categories)
	fmt.Println(designerNames)

	// Store those results in a slice struct from the db package
	for i := 0; i < len(designerNames); i++ {
		if designerNames[i] == "Gucci" ||
			designerNames[i] == "Prada" ||
			designerNames[i] == "Balenciaga" ||
			designerNames[i] == "Carhartt" {
			for j := 0; j < len(categories); j++ {
				err := queries.InsertFilters(ctx, db.InsertFiltersParams{
					DesignerNames: designerNames[i],
					Category:      categories[j],
				})
				if err != nil {
					log.Fatal("Error adding filters to the database", err)
				}
			}
		}
	}
}