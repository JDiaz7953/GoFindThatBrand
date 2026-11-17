package main

import (
	"ClotheScrape/db"
	"context"
	"database/sql"
   _"github.com/mattn/go-sqlite3" // why is this neede?
)

//runs a sql transaction to instert all postings in one sql transactions rather than making multiple calls for every instert
func bulkInsertPostings(ctx context.Context, conn *sql.DB, queries *db.Queries, postings []db.GrailedPosting) error {

	//begin the transaction with error handling and rollback precaution
	tx, err := conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	qtx := queries.WithTx(tx)

	//run loop that adds each element from db.grailedPosting to the database
	for _, p := range postings {
		if err := qtx.CreatePosting(ctx, db.CreatePostingParams{
			ID:            p.ID,
			Title:         p.Title,
			DesignerNames: p.DesignerNames,
			Department:    p.Department,
			Sold:          p.Sold, //Sqlite doesnt take bool values so it must be converted to text in order to store it properly
			Price:         p.Price,
			PriceDrop: 	   p.PriceDrop,
			ImageUrl:      p.ImageUrl,
			PostingUrl:    p.PostingUrl,
		}); err != nil {
			return err
		}
	}

	// Commit the changes
	return tx.Commit()
}
