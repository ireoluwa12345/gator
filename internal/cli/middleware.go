package cli

import (
	"context"

	"github.com/ireoluwa12345/gator/internal/database"
	"github.com/ireoluwa12345/gator/internal/domain"
)

func MiddlewareLoggedIn(handler func(s *domain.State, cmd Command, user database.GetUserRow) error) func(*domain.State, Command) error {
	return func(s *domain.State, cmd Command) error {
		user, err := s.DB.GetUser(context.Background(), s.Config.User)
		if err != nil {
			return err
		}

		return handler(s, cmd, user)
	}
}
