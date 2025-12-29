package cli

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/ireoluwa12345/gator/internal/database"
	"github.com/ireoluwa12345/gator/internal/domain"
)

func HandlerLogin(s *domain.State, cmd Command) error {
	if len(cmd.CommandArgs) < 1 {
		return fmt.Errorf("the handler expects a username")
	}

	name := cmd.CommandArgs[0]

	user, err := s.DB.GetUser(context.Background(), name)
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

func HandlerRegister(s *domain.State, cmd Command) error {
	if len(cmd.CommandArgs) < 1 {
		return fmt.Errorf("the handler expects a username")
	}

	name := cmd.CommandArgs[0]

	_, err := s.DB.CreateUser(context.Background(), database.CreateUserParams{
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
