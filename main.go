package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/ireoluwa12345/gator/internal/config"
	"github.com/ireoluwa12345/gator/internal/database"
	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.Read()
	db, err := sql.Open("postgres", cfg.DatabaseURL)
	dbQueries := database.New(db)
	if err != nil {
		log.Fatalf("failed to read config: %v", err)
		return
	}

	s := &state{
		Config: cfg,
		db:     dbQueries,
	}

	commands := commands{
		commandHandlers: make(map[string]func(*state, command) error),
	}
	commands.register("login", handlerLogin)
	commands.register("register", handlerRegister)
	commands.register("reset", handleReset)
	commands.register("users", handleUsers)
	commands.register("agg", handleAgg)
	commands.register("addfeed", middlewareLoggedIn(handleAddFeed))
	commands.register("feeds", handleFeeds)
	commands.register("follow", middlewareLoggedIn(handleFollow))
	commands.register("following", middlewareLoggedIn(handleFollowing))
	commands.register("unfollow", middlewareLoggedIn(handleUnfollow))

	userInput := os.Args
	if len(userInput) < 2 {
		log.Fatal("no command provided")
		return
	}

	cmd := command{
		name:        userInput[1],
		commandArgs: userInput[2:],
	}

	if err := commands.run(s, cmd); err != nil {
		log.Fatalf("command execution failed: %v", err)
		return
	}
}
