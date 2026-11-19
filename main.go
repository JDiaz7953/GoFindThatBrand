package main

import (
	"ClotheScrape/db"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	_ "github.com/mattn/go-sqlite3" 
	"time"
)

/*
	------ What Should the program do? -------
	
	- Scrape grailed site by a specified brand 
	- Calculate the median price of a brands clothing/appearl 
	- Use that median value calculated to flag clothing and appearl that falls below the median 
	- send flagged postings to a discord bot to notify users 
	- recheck for price changes 

*/


/*


----------- My Questions -----------
	- How am i gonna check if an item has reduced in price that is already in my data
 	- How am i gonna make the point of the database useful? ----> calculate prices to find the median prices of a specifed brand clothing/appearal and flag "Snipes"
   	- recheck listings after a certain amount of time?
	- Figure how am i going to ge the median price

	
*/


/*

TO DOs:
 - Implement pauses between scrapes to avoid rate limiting (done)
 - Calculate the Median price of a brand in the database
 - Display a page of the possible brands that can be searched 
 - flag the clothing/appearal of a specifed brand that falls below the median 
 - those flagged ones get sent to the discord bot
 - recheck for price changes and update those price

 Optimizations: seaching by a specifed brand lowers the ammount of hits in a query making the scrape work faster 
 Stretch: Find a way to automate the key renewal 

*/


/* 

Current Task(s):
 - Implement pauses between scrapes to avoid rate limiting 
 - Calculate the Median price of a brand in the database
 - Figure out how to clean up the reqBody

*/

func main() {

	//Open database connection
	ctx := context.Background()
	conn, err := sql.Open("sqlite3", "./grailed_posting_db.db")
	queries := db.New(conn)
	if err != nil {
		log.Fatal("Unable to use data source name", err)
	}
	defer conn.Close()

	brandName := "Supreme"

	//Partions queries into smaller queries to avoid hitting the 1000 posting limit
	currentLowPrice, previousHighPrice := 0, 100
	const maxPrice = 100000
	for currentLowPrice <= maxPrice {
		currentHighPrice := currentLowPrice + previousHighPrice

		//Search for the current postings using the starting parameters of the price and page which is zero
		page0 := scrapePostings(0, currentLowPrice, currentHighPrice, brandName) //Searching the page

		//If the amount of hits are greater than 1000 and not 0 we divide it to reduce the amount of listings by reducing the price range to filter
		if page0.NbHits > 1000 && previousHighPrice > 0 {
			previousHighPrice /= 2
			continue
		}

		//if it bypasses the it means that the number of hits found can be collected so we save it
		saveHits(ctx, queries, conn, page0.Hits)

		//go through all the pages that are available and save em (I know they can all be saved because there are less than 1000 hits)
		for page := 1; page < page0.NbPages; page++ {
			saveHits(ctx, queries, conn, scrapePostings(page, currentLowPrice, currentHighPrice, brandName).Hits)
			
			//Prevents the program from being rate limited by the site 
			time.Sleep(2 * time.Second)
		}

		//if the hits are less than 500 widen the window
		if page0.NbHits < 500 {
			previousHighPrice *= 2
		}

		// move the binary search up one
		currentLowPrice = currentHighPrice + 1
	}


	//After all the data is in the database calculate the median




	//Flag Postings in the db that fall between the median 



	
	//Send all flagged postings the database 
	

}

// Takes in the root and narrow it down to using results and hits
func saveHits(ctx context.Context, queries *db.Queries, conn *sql.DB, hits []Hits) {

	//builds the struct of GrailedPosting
	var postings []db.GrailedPostingsDb

	// Stores each individual posting into the Hit struct
	for _, h := range hits {

		//Turn price drop into json string to properly place it in the database
		pd, err := json.Marshal(h.PriceDrop)
		if err != nil {
			log.Fatal(err)
		}

		if h.DesignerNames == "Supreme" {
			posting := db.GrailedPostingsDb{
				ID:            h.ID,
				Title:         h.Title,
				DesignerNames: h.DesignerNames,
				Department:    h.Department,
				Sold:          strconv.FormatBool(h.Sold), //Sqlite doesnt take bool values so it must be converted to text in order to store it properly
				Price:         h.Price,
				PriceDrop: sql.NullString{
					String: string(pd),
					Valid:  len(h.PriceDrop) > 0,
				},
				ImageUrl:   h.CoverPhoto.ImageUrl,
				PostingUrl: "https://www.grailed.com/listings/" + strconv.Itoa(h.ID),
			}

			postings = append(postings, posting)

		}
	}

	//store results in the database gonna have to use withtx to process multiple transactions cleanly and effeciently
	err := bulkInsertPostings(ctx, conn, queries, postings)
	if err != nil {
		log.Fatal("Error processing SQL transaction", err)
	} else {
		fmt.Println("SQL Bulk insertion succesful!")
	}
}
