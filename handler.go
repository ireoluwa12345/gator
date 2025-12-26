package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/ireoluwa12345/gator/internal/config"
	"github.com/ireoluwa12345/gator/internal/database"
)

type state struct {
	db *database.Queries
	*config.Config
}

type command struct {
	name        string
	commandArgs []string
}

type commands struct {
	commandHandlers map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	if handler, ok := c.commandHandlers[cmd.name]; ok {
		return handler(s, cmd)
	}
	return fmt.Errorf("unknown command: %s", cmd.name)
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.commandHandlers[name] = f
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.commandArgs) < 1 {
		return fmt.Errorf("the handler expects a username")
	}

	name := cmd.commandArgs[0]

	user, err := s.db.GetUser(context.Background(), name)
	if err != nil {
		return fmt.Errorf("error occurred while logging in: %w", err)
	}

	if user.ID == uuid.Nil {
		return fmt.Errorf("user not found, please register first")
	}

	s.Config.SetUser(name)
	fmt.Println("Login successful")

	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.commandArgs) < 1 {
		return fmt.Errorf("the handler expects a username")
	}

	name := cmd.commandArgs[0]

	_, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		Name:      name,
		ID:        uuid.New(),
		CreatedAt: sql.NullTime{Time: time.Now(), Valid: true},
		UpdatedAt: sql.NullTime{Time: time.Now(), Valid: true},
	})

	if err != nil {
		return fmt.Errorf("error occurred while registering user: %w", err)
	}

	s.Config.SetUser(name)

	fmt.Println("The user has been registered successfully")

	return nil
}

func handleReset(s *state, cmd command) error {
	err := s.db.DeleteUsers(context.Background())
	if err != nil {
		return fmt.Errorf("error occurred while resetting users: %w", err)
	}

	fmt.Println("Users have been reset successfully")

	return nil
}

func handleUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("error occurred while fetching users: %w", err)
	}

	for _, user := range users {
		if user.Name == s.Config.User {
			fmt.Printf("- %s (current)  \n", user.Name)
			continue
		}
		fmt.Printf("- %s  \n", user.Name)
	}

	return nil
}

func handleAgg(s *state, cmd command) error {
	if len(cmd.commandArgs) < 1 {
		return fmt.Errorf("the handler requires a time_between_reqs")
	}

	timeBetweenReqs, err := time.ParseDuration(cmd.commandArgs[0])

	if err != nil {
		return fmt.Errorf("error occurred while parsing time_between_reqs: %w", err)
	}

	ticker := time.NewTicker(timeBetweenReqs)
	fmt.Printf("Collecting feeds every %v\n\n", timeBetweenReqs)
	for ; ; <-ticker.C {
		err := scrapeFeeds(s)
		if err != nil {
			fmt.Printf("Error scraping feeds: %v\n", err)
			break
		}
	}

	return nil
}

func handleAddFeed(s *state, cmd command, user database.GetUserRow) error {
	if len(cmd.commandArgs) < 2 {
		return fmt.Errorf("the handler expects the name and url")
	}

	name := cmd.commandArgs[0]
	url := cmd.commandArgs[1]

	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		Name:      name,
		Url:       url,
		UserID:    user.ID,
		CreatedAt: sql.NullTime{Time: time.Now(), Valid: true},
		UpdatedAt: sql.NullTime{Time: time.Now(), Valid: true},
	})

	if err != nil {
		return fmt.Errorf("error occurred while creating the feed: %w", err)
	}

	_, err = s.db.GetFollowedFeedByID(context.Background(), feed.ID)

	if err != nil {
		return fmt.Errorf("error occurred while fetching followed feed: %w", err)
	} else {
		_, err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{ID: uuid.New(), CreatedAt: sql.NullTime{Time: time.Now(), Valid: true}, UpdatedAt: sql.NullTime{Time: time.Now(), Valid: true}, UserID: feed.UserID, FeedID: feed.ID})

		if err != nil {
			return fmt.Errorf("error occurred while following the feed: %w", err)
		}
	}

	fmt.Println(feed)

	return nil
}

func handleFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetAllFeeds(context.Background())

	if err != nil {
		return fmt.Errorf("error occurred while fetching feeds: %w", err)
	}

	for _, feed := range feeds {
		if feed.UserName.Valid {
			fmt.Printf("- %s by %s (%s)  \n", feed.Name, feed.UserName.String, feed.Url)
			continue
		}
		fmt.Printf("- %s (%s)  \n", feed.Name, feed.Url)
	}

	return nil
}

func handleFollow(s *state, cmd command, user database.GetUserRow) error {
	if len(cmd.commandArgs) < 1 {
		return fmt.Errorf("the handler expects the feed name")
	}

	feedUrl := cmd.commandArgs[0]

	feed, err := s.db.GetFeedByUrl(context.Background(), feedUrl)
	if err != nil {
		return fmt.Errorf("error occurred while fetching feed: %w", err)
	}

	feedFollow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: sql.NullTime{Time: time.Now(), Valid: true},
		UpdatedAt: sql.NullTime{Time: time.Now(), Valid: true},
		UserID:    user.ID,
		FeedID:    feed.ID,
	})

	if err != nil {
		return fmt.Errorf("error occurred while following the feed: %w", err)
	}

	fmt.Printf("%s now following %s\n", feedFollow.UserName, feedFollow.FeedName)

	return nil
}

func handleFollowing(s *state, cmd command, user database.GetUserRow) error {
	followingFeeds, err := s.db.GetFeedsFollowForUser(context.Background(), user.ID)

	if err != nil {
		return fmt.Errorf("error occurred while fetching following feeds: %w", err)
	}

	for _, followingFeed := range followingFeeds {
		fmt.Printf("%s", followingFeed.FeedName)
	}

	return nil
}

func handleUnfollow(s *state, cmd command, user database.GetUserRow) error {
	if len(cmd.commandArgs) < 1 {
		return fmt.Errorf("the handler expects the feed url")
	}

	feedUrl := cmd.commandArgs[0]

	feed, err := s.db.GetFeedByUrl(context.Background(), feedUrl)
	if err != nil {
		return fmt.Errorf("error occurred while fetching feed: %w", err)
	}

	err = s.db.UnfollowFeed(context.Background(), database.UnfollowFeedParams{
		UserID: user.ID,
		FeedID: feed.ID,
	})

	if err != nil {
		return fmt.Errorf("error occurred while unfollowing the feed: %w", err)
	}

	fmt.Printf("%s has unfollowed %s\n", user.Name, feed.Name)

	return nil
}
