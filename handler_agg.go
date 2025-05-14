package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jeffe89/gator/internal/database"
)

func handlerAgg(s *state, cmd command) error {

	// Check command arguments for time between each requests
	if len(cmd.Args) < 1 || len(cmd.Args) > 2 {
		return fmt.Errorf("usage: %v <time_between_reqs", cmd.Name)
	}

	timeBetweenRequests, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("invalid duration: %w", err)
	}

	// Print message to console displaying feed fetch times
	log.Printf("Collecting feeds every %s...", timeBetweenRequests)

	// Setup ticker for infinite loop
	ticker := time.NewTicker(timeBetweenRequests)

	// Infinite loop for scaping feeds
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}

// Function for aggregation
func scrapeFeeds(s *state) {

	// Get next feed to fetch
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		log.Println("Couldn't get the next feeds to fetch", err)
		return
	}

	//Print success message and call helper function to scrape the feed
	log.Println("Found a feed to fetch!")
	scrapeFeed(s.db, feed)
}

// Helper function to scrape a individual feed
func scrapeFeed(db *database.Queries, feed database.Feed) {

	// Mark feed as fetched
	_, err := db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		log.Printf("Couldn't mark feed %s fetched: %v", feed.Name, err)
		return
	}

	// Collect specific feed data
	feedData, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		log.Printf("Couldn't collect feed %s: %v", feed.Name, err)
		return
	}

	// Loop over each feed item and create posts
	for _, item := range feedData.Channel.Item {
		publishedAt := sql.NullTime{}
		if t, err := time.Parse(time.RFC1123Z, item.PubDate); err == nil {
			publishedAt = sql.NullTime{
				Time:  t,
				Valid: true,
			}
		}

		// Use sql query to create and add post to database
		_, err = db.CreatePost(context.Background(), database.CreatePostParams{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			FeedID:    feed.ID,
			Title:     item.Title,
			Description: sql.NullString{
				String: item.Description,
				Valid:  true,
			},
			Url:         item.Link,
			PublishedAt: publishedAt,
		})
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value violates unique constaint") {
				continue
			}
			log.Printf("Couldn't create post: %v", err)
			continue
		}
	}

	// Log data on how many feeds collected and posts found
	log.Printf("Feed %s collected, %v posts found", feed.Name, len(feedData.Channel.Item))
}
