package cli

import (
	"context"
	"fmt"

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
