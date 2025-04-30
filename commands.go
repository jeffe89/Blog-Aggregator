package main

import "errors"

// Define command struct for command arguments
type command struct {
	Name string
	Args []string
}

// Define commands struct to hold all commands the CLI can handle
type commands struct {
	registeredCommands map[string]func(*state, command) error
}

// Register method for commands struct
func (c *commands) register(name string, f func(*state, command) error) {

	// Register a new handler function for a command name
	c.registeredCommands[name] = f
}

// Run method for commands struct
func (c *commands) run(s *state, cmd command) error {

	//Check if command is found in handlers map
	handler, exists := c.registeredCommands[cmd.Name]
	if !exists {
		return errors.New("command not found")
	}

	// Call handler function and return errors (if any)
	return handler(s, cmd)
}
