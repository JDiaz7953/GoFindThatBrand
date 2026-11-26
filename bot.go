package main

import (
	"ClotheScrape/db"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func runBot(ctx context.Context, q *db.Queries, currentJob db.ScrapingQueueRow) error {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, falling back to system environment")
	}
	token := os.Getenv("DISCORD_KEY")
	channelId := os.Getenv("DISCORD_CHANNEL_ID")

	//Call to initalize the discord bot to send request
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		return err
	}

	//get all the values needed from the most recent process
	var FlaggedPostings []db.GetFlaggedItemRow
	FlaggedPostings, err = q.GetFlaggedItem(ctx, db.GetFlaggedItemParams{
		Brand:    currentJob.DesignerNames,
		Category: currentJob.Category,
	})
	if err != nil {
		return err
	}
	amountFlagged := len(FlaggedPostings)
	for i, fp := range FlaggedPostings {
		messages := discordgo.MessageEmbed{
			URL:         fp.PostingUrl,
			Title:       fp.Title + " (" + fp.Category + ")",
			Description: "$" + strconv.Itoa(fp.Price),
			Image: &discordgo.MessageEmbedImage{
				URL: fp.ImageUrl,
			},
		}
		dg.ChannelMessageSendEmbed(channelId, &messages)
		fmt.Println(i, "/", amountFlagged)
		time.Sleep(10 * time.Second)
	}

	return nil
}
