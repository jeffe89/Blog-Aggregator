package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jeffe89/gator/internal/database"
)

// Function to create a new feed follow record for current user
func handlerFollow(s *state, cmd command, user database.User) error {

	// Check for correct amount of arguments
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <feed_url>", cmd.Name)
	}

	// Find feed by URL
	feed, err := s.db.GetFeedByURL(context.Background(), cmd.Args[0])
	if err != nil {
		return fmt.Errorf("couldn't get feed: %w", err)
	}

	// Create feed follow
	ffRow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("couldn't create feed follow: %w", err)
	}

	// Print result
	fmt.Println("Feed follow created:")
	printFeedFollow(ffRow.UserName, ffRow.FeedName)
	return nil
}

// Function to list all followed feeds for a user
func handlerListFeedFollows(s *state, cmd command, user database.User) error {

	// Get all feeds followed by specific user
	feedFollows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("couldn't get feed follows: %w", err)
	}

	// Check if user doesn't follow any feeds
	if len(feedFollows) == 0 {
		fmt.Println("no feed follows found for this user.")
		return nil
	}

	// Print out all feeds followed by user
	fmt.Printf("Feed follows for user %s:\n", user.Name)
	for _, ff := range feedFollows {
		fmt.Printf("* %s\n", ff.FeedName)
	}

	return nil
}

// Function to unfollow a feed for a user
func handlerUnfollow(s *state, cmd command, user database.User) error {

	// Check URL argument provided
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <feed_url", cmd.Name)
	}

	// Get feed from URL argument
	feed, err := s.db.GetFeedByURL(context.Background(), cmd.Args[0])
	if err != nil {
		return fmt.Errorf("couldn't get feed: %w", err)
	}

	// Call query on database to unfollow feed
	err = s.db.DeleteFeedFollow(context.Background(), database.DeleteFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return fmt.Errorf("couldn't delete feed follow: %w", err)
	}

	// Print success message
	fmt.Printf("%s unfollowed successfully!\n", feed.Name)
	return nil
}

// Helper function to print out feed follow
func printFeedFollow(username, feedname string) {
	fmt.Printf("* User:      %s\n", username)
	fmt.Printf("* Feed:	  %s\n", feedname)
}
