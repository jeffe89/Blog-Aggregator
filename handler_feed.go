package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jeffe89/gator/internal/database"
)

// Function to add feed to database
func handlerAddFeed(s *state, cmd command, user database.User) error {

	// Check both arguments provided
	if len(cmd.Args) != 2 {
		return fmt.Errorf("usage: %s <name> <url>", cmd.Name)
	}

	// Gather name and url provided
	name := cmd.Args[0]
	url := cmd.Args[1]

	// Create the feed
	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		Name:      name,
		Url:       url,
	})
	if err != nil {
		return fmt.Errorf("couldn't create feed: %w", err)
	}

	// Create a feed follow
	feedFollow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("couldn't create feed follow: %w", err)
	}

	//Print success message and feed to console
	fmt.Println("Feed created successfully:")
	printFeed(feed, user)
	fmt.Println()
	fmt.Println("Feed followed successfully:")
	printFeedFollow(feedFollow.UserName, feedFollow.FeedName)
	fmt.Println("=====================================")

	return nil
}

// Function to list all feeds currently in database
func handlerListFeeds(s *state, cmd command) error {

	// Call the SQL query to retrieve all feeds
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return err
	}

	// Check if no feeds are found
	if len(feeds) == 0 {
		fmt.Println("No feeds found.")
		return nil
	}

	// Print the feeds
	fmt.Printf("Found %d feeds:\n", len(feeds))
	for _, feed := range feeds {
		user, err := s.db.GetUserById(context.Background(), feed.UserID)
		if err != nil {
			return fmt.Errorf("couldn't get user: %w", err)
		}
		printFeed(feed, user)
		fmt.Println("=====================================")
	}

	return nil
}

// Helper function to format and print feed
func printFeed(feed database.Feed, user database.User) {
	fmt.Printf("* ID:		%s\n", feed.ID)
	fmt.Printf("* Created:   %v\n", feed.CreatedAt)
	fmt.Printf("* Updated:   %v\n", feed.UpdatedAt)
	fmt.Printf("* Name:	  %s\n", feed.Name)
	fmt.Printf("* URL:	   %s\n", feed.Url)
	fmt.Printf("* User:	  %s\n", user.Name)
	fmt.Printf("* LastFetchedAt: %v\n", feed.LastFetchedAt.Time)
}
