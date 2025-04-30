package main

import (
	"fmt"
)

// Function to handle Login
func handlerLogin(s *state, cmd command) error {

	// Check command has at least one argument
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}

	// Gather username and utilize SetUser method from config package
	userName := cmd.Args[0]
	if err := s.cfg.SetUser(userName); err != nil {
		return fmt.Errorf("failed to set user: %w", err)
	}

	// Print message to terminal stating the user has been set
	fmt.Printf("User set to %s\n", userName)

	return nil
}
