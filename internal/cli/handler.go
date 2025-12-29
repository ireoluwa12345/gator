package cli

import (
	"context"
	"fmt"
	"time"

	"github.com/ireoluwa12345/gator/internal/aggregator"
	"github.com/ireoluwa12345/gator/internal/domain"
)

type Command struct {
	Name        string
	CommandArgs []string
}

type Commands struct {
	CommandHandlers map[string]func(*domain.State, Command) error
}

func (c *Commands) Run(s *domain.State, cmd Command) error {
	if handler, ok := c.CommandHandlers[cmd.Name]; ok {
		return handler(s, cmd)
	}
	return fmt.Errorf("unknown command: %s", cmd.Name)
}

func (c *Commands) Register(name string, f func(*domain.State, Command) error) {
	c.CommandHandlers[name] = f
}

func HandleReset(s *domain.State, cmd Command) error {
	err := s.DB.DeleteUsers(context.Background())
	if err != nil {
		return fmt.Errorf("error occurred while resetting users: %w", err)
	}

	fmt.Println("Users have been reset successfully")

	return nil
}

func HandleUsers(s *domain.State, cmd Command) error {
	users, err := s.DB.GetUsers(context.Background())
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

func HandleAgg(s *domain.State, cmd Command) error {
	if len(cmd.CommandArgs) < 1 {
		return fmt.Errorf("the handler requires a time_between_reqs")
	}

	timeBetweenReqs, err := time.ParseDuration(cmd.CommandArgs[0])

	if err != nil {
		return fmt.Errorf("error occurred while parsing time_between_reqs: %w", err)
	}

	ticker := time.NewTicker(timeBetweenReqs)
	fmt.Printf("Collecting feeds every %v\n\n", timeBetweenReqs)
	for ; ; <-ticker.C {
		err := aggregator.ScrapeFeeds(s)
		if err != nil {
			fmt.Printf("Error scraping feeds: %v\n", err)
			break
		}
	}

	return nil
}
