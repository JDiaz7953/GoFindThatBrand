package main

import (
	"ClotheScrape/db"
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3" // why is this neede?
)

//runs a sql transaction to instert all postings in one sql transactions rather than making multiple calls for every instert 
func bulkInsertPostings(ctx context.Context, db *sql.DB, postings []db.GrailedPosting) error {
	//begin the transaction with error handling and rollback precaution


	//run loop that adds each element from db.grailedPosting to the database 


	// Commit the changes

	return nil
}