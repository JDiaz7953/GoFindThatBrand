package dbUtil

import (
	"ClotheScrape/db"
	"context"
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3" // why is this neede?
)

// runs a sql transaction to instert all postings in one sql transactions rather than making multiple calls for every instert
func BulkInsertPostings(ctx context.Context, conn *sql.DB, queries *db.Queries, postings []db.GrailedPostingsDb) error {

	//begin the transaction with error handling and rollback precaution
	tx, err := conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	qtx := queries.WithTx(tx)

	//run loop that adds each element from db.grailedPosting to the database
	for _, p := range postings {
		fmt.Println(p)
		if err := qtx.CreatePosting(ctx, db.CreatePostingParams{
			ID:            p.ID,
			Title:         p.Title,
			DesignerNames: p.DesignerNames,
			Department:    p.Department,
			Category:      p.Category,
			Sold:          p.Sold,
			Price:         p.Price,
			PriceDrop:     sql.NullString{String: p.PriceDrop.String, Valid: p.PriceDrop.Valid},
			ImageUrl:      p.ImageUrl,
			PostingUrl:    p.PostingUrl,
			Flagged:       p.Flagged,
		}); err != nil {
			return err
		}
	}

	// Commit the changes
	return tx.Commit()
}
