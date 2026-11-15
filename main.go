package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3" // why is this neede?
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	
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
 - Connect it to a database and only input a specifed brand
 - flag not a specific clothing item but rather and item of that brand that falls well within the median
 - those flagged ones get sent to the discord bot

 - How am i gonna check if an item has reduced in price that is already in my data
 - How am i gonna make the point of the database useful?
   - recheck listings after a certain amount of time?
-  Figure how am i going to ge the median price


- Current sprint: input specifed data

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
	conn, err := sql.Open("sqlite3", "./")
	if err != nil{
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
	req, _ := requestHandler(url, key)

	//Performs the request to the site
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal("Unexpected Status Code ", resp.StatusCode)
	}
	defer resp.Body.Close()

	resBody, err := io.ReadAll(resp.Body)

	//After getting readable data it needs to me unmarshaled pointed to the result struct
	// build struct with URL to the posting PostingUrl:https://www.grailed.com/listings/<Hit.ID> ---> takes to posting
	var root Root
	err = json.Unmarshal(resBody, &root)
	if err != nil {
		log.Fatal("Error extracting data")
	}

	//create loop to continously go through pages

	data := root.Results[0]
	//fmt.Println(r)
	for _, v := range data.Hits {
		if v.DesignerName == "Supreme" {
			fmt.Println(v)
		}
	}

}
