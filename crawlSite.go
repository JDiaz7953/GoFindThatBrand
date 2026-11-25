package main

import (
	"ClotheScrape/db"
	"context"
	"database/sql"
	"log"
	"os/exec"
	"time"
)

func WriteKey() {

	//if the program cant run the python script then its better to end the program
	cmd := exec.Command("python", "-m", "scripts.captureKey")
	err := cmd.Run()
	if err != nil {
		log.Fatal("Unable to run python script canceling program")
	}
}

func CrawlSite(ctx context.Context, queries *db.Queries, conn *sql.DB, currentJob db.ScrapingQueueRow) error {
	
	//Captures the max amount for the specifed filter inputted into query
	res, err := scrapePostings(0, 0, 1000000, currentJob.DesignerNames, currentJob.Category) // if results is empty return an error to try to get a key again
	if err != nil {
		WriteKey()
		return err
	}

	//Partions queries into smaller queries in order to collect all the neccessary postings without having to worry about hitting the 1000 posting limit
	currentLowPrice, priceStep := 0, 100
	maxPrice := res.FacetStats.Price_i.Max
	for currentLowPrice <= maxPrice {
		currentHighPrice := currentLowPrice + priceStep - 1

		//Search for the current postings using the starting parameters of the price and page which is zero
		page0, err := scrapePostings(0, currentLowPrice, currentHighPrice, currentJob.DesignerNames, currentJob.Category) //Searching the page
		if err != nil {
			return err
		}

		//If the amount of hits are greater than 1000 then reduce the highPrice
		if page0.NbHits > 1000 && priceStep > 1 {
			priceStep /= 2
			continue
		}

		//if it bypasses the it means that the number of hits found can be collected so we save it
		saveHits(ctx, queries, conn, page0.Hits)

		//go through all the pages that are available and save em (I know they can all be saved because there are less than 1000 hits)
		for page := 1; page < page0.NbPages; page++ {
			res, err := scrapePostings(page, currentLowPrice, currentHighPrice, currentJob.DesignerNames, currentJob.Category)
			if err != nil {
				return err
			}

			saveHits(ctx, queries, conn, res.Hits)

			//Prevents the program from being rate limited by the site
			time.Sleep(3 * time.Second)
		}

		//if the hits are less than 500 widen the window
		if page0.NbHits < 500 {
			priceStep *= 2
		}

		// move the binary search up one
		currentLowPrice = currentHighPrice + 1
	}

	return nil
}
