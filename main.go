package main

import (
	"ClotheScrape/db"
	"context"
	"database/sql"
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"math"
	_ "math"
	"strconv"
	"time"
	_ "time"

	"os/exec"

	_ "github.com/mattn/go-sqlite3"

	"ClotheScrape/dbUtil"
	"ClotheScrape/models"
)

/*
	------ What Should the program do? -------
	- automate key renewal
	- Scrape grailed site by a specified brand and category that rotate on a queue (Important!! makes it a lot more effecient and reliable)
	- Set cron job to automate on a schedule
	- send flagged postings to a discord bot to notify users
	- recheck for price changes

	------DONE-----
	- Calculate the median price of a brands clothing/appearl (done)
	- Use that median value calculated to flag clothing and appearl that falls below the median (done)
*/

/*


----------- My Notes -----------
	- How am i gonna check if an item has reduced in price that is already in my data
   	- recheck listings after a certain amount of time?
	- How am I going to allow it to scale to other brands -----> Using a queue system that utalizes arry
	- Reducing search filters to fit me (Size, brand, category) produces a more effecient program may have to lean towards that

	 Optimizations: seaching by a specifed brand lowers the ammount of hits in a query making the scrape work faster
 	Stretch: Find a way to automate the key renewal
*/

/*

TO DOs:
 - Display a page of the possible brands that can be searched
 - those flagged ones get sent to the discord bot
 - recheck for price changes and update those price
 - Figure out how to establish using multiple different brands and clothing type
 - turn it into a cron job that waits for the key to be written and processes accordingly

*/

/*
Current Task(s):
	- write script so they is captured
	- write the loop
	- simplify loop
	- recheck pricing
	- Establish a connection with discord bot so it sends listing  Update notify_at (LAST)

*/

//go:embed sql/schema.sql
var ddl string

func main() {

	//Open database connection
	ctx := context.Background()
	conn, err := sql.Open("sqlite3", "./grailed_posting_db.db")
	if err != nil {
		log.Fatal("Unable to use data source name", err)
	}
	queries := db.New(conn)
	

	// create tables
	if _, err := conn.ExecContext(ctx, ddl); err != nil {
		log.Fatal("error creating tabel", err)
	}

	defer conn.Close()

	//Call the python program to wrtie the key
	cmd := exec.Command("python", "-m", "scripts.captureKey")
	err = cmd.Run()

	for {

		//sleep to avoid rate limiting
		time.Sleep(60 * time.Second)

		//Get the category and brand next on the queue
		currentJob, err := queries.ScrapingQueue(ctx)
		if err != nil {
			time.Sleep(60 * time.Second)
			continue
		}

		//Runs crawler loop that scrapes and saves postings
		err = CrawlSite(ctx, queries, conn, currentJob )
		if err != nil {
			continue
		}

		//After all the data is in the database calculate the median and pass the decimal value to medianPrice
		processMedian, err := queries.CalculateMedianPriceOfBrand(ctx, db.CalculateMedianPriceOfBrandParams{
			Brand:    currentJob.DesignerNames,
			Category: currentJob.Category,
		})
		if err != nil {
			fmt.Println("Error processing median", err)
		}

		if processMedian.Valid == true {
			medianPrice := processMedian.Float64

			//Flag Postings in the db that fall below the median (lowMedian and highMedian are used as a band to avoid collecting pure junk postings).
			lowMedian := math.Ceil(((medianPrice * 0.3) * 100)) / 100
			highMedian := math.Ceil(((medianPrice * 0.7) * 100)) / 100

			//Flags posting using SQL querying
			err = queries.FlagDeals(ctx, db.FlagDealsParams{
				Low:      lowMedian,
				High:     highMedian,
				Brand:    currentJob.DesignerNames,
				Category: currentJob.Category,
			})
			if err != nil {
				fmt.Println("Error flagging values ")
			}
		}

		//Stamp a job with when it was last completed
		err = queries.MarkJobCompleted(ctx, currentJob.ID)
		if err != nil {
			fmt.Print("Error marking job complete")
		}

	}

}

// Takes in the root and narrow it down to using results and hits
func saveHits(ctx context.Context, queries *db.Queries, conn *sql.DB, hits []models.Hits) {

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
	err := dbUtil.BulkInsertPostings(ctx, conn, queries, postings)
	if err != nil {
		log.Fatal("Error processing SQL transaction", err)
	} else {
		fmt.Println("SQL Bulk insertion succesful!")
	}
}
