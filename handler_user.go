package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jeffe89/gator/internal/database"
)

// Function to handle registering a new user
func handlerRegister(s *state, cmd command) error {

	// Check command has at least one argument
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}

	// Gather username and create new user in database
	name := cmd.Args[0]
	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      name,
	})
	if err != nil {
		// Check if errir is due to duplicate user
		if strings.Contains(err.Error(), "duplicate key") ||
			strings.Contains(err.Error(), "unique constraint") {
			fmt.Printf("Error: User with name '%s' already exists\n", name)
			os.Exit(1)
		}
		return fmt.Errorf("couldn't create user: %w", err)
	}

	if err := s.cfg.SetUser(user.Name); err != nil {
		return fmt.Errorf("failed to set user: %w", err)
	}

	fmt.Println("User created successfully:")
	printUser(user)
	return nil
}

// Function to handle Login
func handlerLogin(s *state, cmd command) error {

	// Check command has at least one argument
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}

	// Gather username and utilize SetUser method from config package
	userName := cmd.Args[0]

	//Check if username exists in database
	_, err := s.db.GetUser(context.Background(), userName)
	if err != nil {
		return fmt.Errorf("couldn't find user: %w", err)
	}

	if err := s.cfg.SetUser(userName); err != nil {
		return fmt.Errorf("failed to set user: %w", err)
	}

	// Print message to terminal stating the user has been set
	fmt.Printf("User set to %s\n", userName)

	return nil
}

// Function to handle displaying all users in database
func handlerUsers(s *state, cmd command) error {

	//Get all users from database
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't list users: %w", err)
	}

	// Get the current user from config
	currentUser := s.cfg.CurrentUserName

	for _, user := range users {
		if currentUser != "" && user.Name == currentUser {
			fmt.Printf("* %s (current)\n", user.Name)
		} else {
			fmt.Printf("* %s\n", user.Name)
		}
	}

	return nil
}

// Helper function to print user from database
func printUser(user database.User) {
	fmt.Printf(" * ID:		%v\n", user.ID)
	fmt.Printf(" * Name:	%v\n", user.Name)
}
