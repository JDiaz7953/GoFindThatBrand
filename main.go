package main

import (
	"ClotheScrape/db"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "io"
	"log"
	_ "net/http"
	"os"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3" // why is this neede?
)

/*
 ---- What I have so far -----
 - The program will scrape the site for 3 hoursish before needing a new key


*/

/*
	------ Classify as finished When the Following is done -------

	- Pivot into a deal finder take user input before running the program on what exactly what I am looking for or hard code it
	- connect it to a discord bot
	- input data to a DB
	- flag specifed listings probably by criteria


*/

/*

TO DO:
 - Connect it to a database and only input a specifed brand (done)
 - go through multiplie pages

 - flag not a specific clothing item but rather and item of that brand that falls well within the median
 - those flagged ones get sent to the discord bot

 - How am i gonna check if an item has reduced in price that is already in my data
 - How am i gonna make the point of the database useful?
   - recheck listings after a certain amount of time?
-  Figure how am i going to ge the median price


- Current sprint: input specifed data into the db

*/

func readKey(path string) (string, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return "", nil
	}

	return strings.TrimSpace(string(b)), nil
}

func main() {

	//Open database connection
	ctx := context.Background()
	conn, err := sql.Open("sqlite3", "./grailed_posting_db.db")
	queries := db.New(conn)
	if err != nil {
		log.Fatal("Unable to use data source name", err)
	}
	defer conn.Close()

	// Reads sessionKey for the program to access Grailed
	key, err := readKey(`C:\Users\jaide\Documents\Workspace\Projects\GoFindThatBrand\key.txt`)
	if err != nil {

		log.Fatal(err)
	}

	//Builds the request with the proper headers to communicate with the site
	url := "https://mnrwefss2q-dsn.algolia.net/1/indexes/*/queries"
	//req, _ := requestHandler(url, key, "", "", "")

	//Performs the request to the site
	//resp, err := http.DefaultClient.Do(req)
	//if err != nil {
	//	log.Fatal("Unexpected Status Code ", resp.StatusCode)
	//}
	//defer resp.Body.Close()

	//After getting readable data it needs to me unmarshaled pointed to the result struct to map to the json tags of the models
	//var root Root
	//resBody, err := io.ReadAll(resp.Body)
	//err = json.Unmarshal(resBody, &root)
	//if err != nil {
	//	log.Fatal("Error extracting data")
	//}

	//index [0] is the location of all the postings
	//data := root.Results[0]

	//builds the struct of GrailedPosting
	//var postings []db.GrailedPostingsDb

	//for _, h := range data.Hits {

	//Turn price drop into json string to properly place it in the database
	// 	pd, err := json.Marshal(h.PriceDrop)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}

	// 	if h.DesignerNames == "Supreme" {
	// 		posting := db.GrailedPostingsDb{
	// 			ID:            h.ID,
	// 			Title:         h.Title,
	// 			DesignerNames: h.DesignerNames,
	// 			Department:    h.Department,
	// 			Sold:          strconv.FormatBool(h.Sold), //Sqlite doesnt take bool values so it must be converted to text in order to store it properly
	// 			Price:         h.Price,
	// 			PriceDrop: sql.NullString{
	// 				String: string(pd),
	// 				Valid:  len(h.PriceDrop) > 0,
	// 			},
	// 			ImageUrl:   h.CoverPhoto.ImageUrl,
	// 			PostingUrl: "https://www.grailed.com/listings/" + strconv.Itoa(h.ID),
	// 		}

	// 		postings = append(postings, posting)

	// 	}
	// }

	// //store results in the database gonna have to use withtx to process multiple transactions cleanly and effeciently
	// err = bulkInsertPostings(ctx, conn, queries, postings)
	// if err != nil {
	// 	log.Fatal("Error processing SQL transaction", err)
	// } else {
	// 	fmt.Println("SQL Bulk insertion succesful!")
	// }

	/*
		//

		---- Looping pages ------
		- A new request and response needs to be processed everything
		- the page query number needs to change
		- only allowed 1000 hits per page I have to bypass this, each query is a range to gather 1000 hits
		- the solutions is to have a max price and high and low high is how high


	*/

	currentLowPrice, previousHighPrice := 0, 100
	const maxPrice = 100000
	for currentLowPrice <= maxPrice {
		currentHighPrice := currentLowPrice + previousHighPrice

		//Search for the current parameters of the price and page which is zero
		page0 := requestHandler(url, key, 0, currentLowPrice, currentHighPrice) //Searching the page

		//If the amount of hits are greater than 1000 and not 0 we divide it to reduce the amount of listings by reducing the price range to filter
		if page0.NbHits > 1000 && previousHighPrice > 0 {
			previousHighPrice /= 2
			continue
		}

		//if it bypasses the it means that the number of hits found can be collected so we save it
		saveHits(ctx, queries, conn, page0.Hits)

		//go through all the pages that are available and save em (I know they can all be saved because there are less than 1000 hits)
		for page := 1; page < page0.NbPages; page++ {
			saveHits(ctx, queries, conn, requestHandler(url, key, page, currentLowPrice, currentHighPrice).Hits)
		}

		//if the hits are less than 500 widen the window
		if page0.NbHits < 500 {
			previousHighPrice *= 2
		}

		// move the binary search up one
		currentLowPrice = currentHighPrice + 1

	}

}

// Takes in the root and narrow it down to using results and hits
func saveHits(ctx context.Context, queries *db.Queries, conn *sql.DB, hits []Hits) {

	//builds the struct of GrailedPosting
	var postings []db.GrailedPostingsDb

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
