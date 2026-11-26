package main

import (
	"ClotheScrape/db"
	"context"
	"os"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
)

func runBot(ctx context.Context,q *db.Queries,  currentJob db.ScrapingQueueRow ) error {
	token := os.Getenv("DISCORD_KEY")
	channelId := os.Getenv("DISCORD_CHANNEL_ID")

	//Call to initalize the discord bot to send request
	dg, err := discordgo.New("Bot "+ token)
	if err != nil {
		return err
	}

	//get all the values needed from the most recent process
	var FlaggedPostings []db.GetFlaggedItemRow
	FlaggedPostings, err = q.GetFlaggedItem(ctx, db.GetFlaggedItemParams{
		Brand: currentJob.DesignerNames,
		Category: currentJob.Category,
	})

	

	for _, fp := range FlaggedPostings{
		messages := discordgo.MessageEmbed{
			URL: fp.PostingUrl,
			Title: fp.Title,
			Description: "$"+strconv.Itoa(fp.Price),
			Image: &discordgo.MessageEmbedImage{
				URL: fp.ImageUrl,
			},
		}
		dg.ChannelMessageSendEmbed(channelId, &messages)
		time.Sleep(2 * time.Second)
	}

	return nil
}