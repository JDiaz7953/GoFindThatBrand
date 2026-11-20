package main

import (
	"ClotheScrape/db"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"strconv"
	"time"
)

/*
	------ What Should the program do? -------

	- Scrape grailed site by a specified brand (done)..... can be more dynamic
	- Calculate the median price of a brands clothing/appearl (done)
	- Use that median value calculated to flag clothing and appearl that falls below the median (done)
	- send flagged postings to a discord bot to notify users
	- recheck for price changes
	- Make filters more dynamic

*/

/*


----------- My Notes -----------
	- How am i gonna check if an item has reduced in price that is already in my data
   	- recheck listings after a certain amount of time?
	- How am I going to allow it to scale to other brands -----> Using a queue system that utalizes arry

	 Optimizations: seaching by a specifed brand lowers the ammount of hits in a query making the scrape work faster
 	Stretch: Find a way to automate the key renewal
*/

/*

TO DOs:
 - Implement pauses between scrapes to avoid rate limiting (done)
 - Calculate the Median price of a brand in the database (done)
 - Implement the median under a band so it can find less junk (done)

 - Display a page of the possible brands that can be searched
 - flag the clothing/appearal of a specifed brand that falls below the median (Done)
 - those flagged ones get sent to the discord bot
 - recheck for price changes and update those price
 - Figure out how to establish using multiple different brands and clothing type

*/

/*
Current Task(s):
	- look into deduping ids (x)
	- Test scraper to collect median and flag postings
	- Establish a connection with discord bot so it sends listing
	- Update notify_at
	- recheck pricing??


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
	maxPrice := scrapePostings(0, 0, 0, brandName).FacetStats.Price_i.Max
	for currentLowPrice <= maxPrice {
		currentHighPrice := currentLowPrice + previousHighPrice - 1

		//Search for the current postings using the starting parameters of the price and page which is zero
		page0 := scrapePostings(0, currentLowPrice, currentHighPrice, brandName) //Searching the page

		//If the amount of hits are greater than 1000 and not 0 we divide it to reduce the amount of listings by reducing the price range to filter
		if page0.NbHits > 1000 && previousHighPrice > 1 {
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

	//After all the data is in the database calculate the median and pass the decimal value to medianPrice
	processMedian, err := queries.CalculateMedianPriceOfBrand(ctx)
	if err != nil {
		log.Fatal("Error Processing the price median for the collected postings!", err)
	}
	medianPrice := processMedian.Float64

	//Flag Postings in the db that fall below the median (lowMedian and highMedian are used as a band to avoid collecting pure junk postings).
	lowMedian := medianPrice * 0.3
	highMedian := medianPrice * 0.6
	fmt.Println(lowMedian)
	fmt.Println(highMedian)

	//Flags posting using SQL querying
	err = queries.FlagDeals(ctx, db.FlagDealsParams{
		Low:  lowMedian,
		High: highMedian,
	})
	if err != nil {
		log.Fatal("Error flagging postings", err)
	}

	//Integrate discord bot to recieve flagged postings




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
				Category:      h.Catagory,
				Sold:          strconv.FormatBool(h.Sold), //Sqlite doesnt take bool values so it must be converted to text in order to store it properly
				Price:         h.Price,
				PriceDrop:     sql.NullString{String: string(pd), Valid: len(h.PriceDrop) > 0},
				ImageUrl:      h.CoverPhoto.ImageUrl,
				PostingUrl:    "https://www.grailed.com/listings/" + strconv.Itoa(h.ID),
				Flagged:       strconv.FormatBool(false),
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
