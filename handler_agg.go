package main

import (
	"context"
	"fmt"
)

func handlerAgg(s *state, cmd command) error {

	// Call fetchFeed with specified URL
	feed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return fmt.Errorf("couldn't fetch feed: %w", err)
	}

	// Print the feed to console
	fmt.Printf("Feed: %+v\n", feed)

	return nil
}
