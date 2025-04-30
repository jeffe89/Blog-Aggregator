package main

import (
	"gator/internal/config"
	"log"
	"os"
)

// Define state struct which holds pointer to config
type state struct {
	cfg *config.Config
}

func main() {

	// Read the config file
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	// Store config in a new State instance
	programState := &state{
		cfg: &cfg,
	}

	// Create a new commands instance with initialized map
	cmds := commands{
		registeredCommands: make(map[string]func(*state, command) error),
	}

	// Register the login handler
	cmds.register("login", handlerLogin)

	// Check if enough arguments were provided
	if len(os.Args) < 2 {
		log.Fatal("Usage: cli <command> [args...]")
		return
	}

	// Create command instance
	cmd := command{
		Name: os.Args[1],
		Args: os.Args[2:],
	}

	// Run the command
	if err := cmds.run(programState, cmd); err != nil {
		log.Fatal(err)
	}
}
