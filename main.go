package main

import (
	"ClotheScrape/db"
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"log"
	"math"
	_ "math"
	"time"
	_ "time"

	_ "github.com/mattn/go-sqlite3"
)

/*
	------ What Should the program do? -------
	- send flagged postings to a discord bot to notify users

	------DONE-----
	- Calculate the median price of a brands clothing/appearl (done)
	- Use that median value calculated to flag clothing and appearl that falls below the median (done)
	- automate key renewal
	- Scrape grailed site by a specified brand and category that rotate on a queue (Important!! makes it a lot more effecient and reliable)
	- recheck for price changes
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

	

*/

/*
Current Task(s):
	- Check for sold items
	-  add size to the discord bot

*/

//go:embed sql/schema.sql
var ddl string

func main() {

	//Open database connection
	ctx := context.Background()
	conn, err := sql.Open("sqlite3", "./sql/grailed_posting_db.db")
	if err != nil {
		log.Fatal("Unable to use data source name", err)
	}

	// create tables
	if _, err := conn.ExecContext(ctx, ddl); err != nil {
		log.Fatal("error creating tabel", err)
	}
	queries := db.New(conn)

	defer conn.Close()

	fmt.Println("Starting webcrawler")
	for {

		//sleep to avoid rate limiting
		time.Sleep(60 * time.Second)

		//Get the category and brand next on the queue
		currentJob, err := queries.ScrapingQueue(ctx)
		if err != nil {
			time.Sleep(60 * time.Second)
			continue
		}
		fmt.Println("Processing: " + currentJob.DesignerNames + " and " + currentJob.Category)

		//Runs crawler loop that scrapes and saves postings
		err = CrawlSite(ctx, queries, conn, currentJob)
		if err != nil {
			continue
		}

		fmt.Println("Crawler succesfully processed: " + currentJob.DesignerNames + " and " + currentJob.Category + "!!")

		fmt.Println("Calculating Median for " + currentJob.DesignerNames + " and " + currentJob.Category + "....")
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

		fmt.Println("Processing Flagged Postings to GoodGrailed bot")
		err = runBot(ctx, queries, currentJob)
		if err != nil {
			fmt.Println("Error sending data to discord channel")
			continue
		}
		fmt.Println("All messages processed Succesfully to GoodGrailed bot!!")

	}

}
