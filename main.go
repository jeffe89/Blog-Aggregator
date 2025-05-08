package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/jeffe89/gator/internal/config"
	"github.com/jeffe89/gator/internal/database"
	_ "github.com/lib/pq"
)

// Define state struct which holds pointer to config
type state struct {
	db  *database.Queries
	cfg *config.Config
}

func main() {

	// Read the config file
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	// Open connection to database
	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		log.Fatalf("error connecting to database: %v", err)
	}

	// Creat queries from the SQLC-generated code
	defer db.Close()
	dbQueries := database.New(db)

	// Store config in a new State instance
	programState := &state{
		db:  dbQueries,
		cfg: &cfg,
	}

	// Create a new commands instance with initialized map
	cmds := commands{
		registeredCommands: make(map[string]func(*state, command) error),
	}

	// Register the handler functions
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", handlerAddFeed)
	cmds.register("feeds", handlerListFeeds)

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
