package main

import (
	"context"
	"fmt"
)

// Function to handle resetting the database
func handlerReset(s *state, cmd command) error {

	//Call the database method to delete all users
	err := s.db.DeleteAllUsers(context.Background())
	// If any error encountered, display and return error
	if err != nil {
		return fmt.Errorf("couldn't delete users: %w", err)
	}

	// If successful, display a success message
	fmt.Println("Database reset successfully!")
	return nil
}
