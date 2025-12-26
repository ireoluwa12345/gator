package main

import (
	"context"

	"github.com/ireoluwa12345/gator/internal/database"
)

func middlewareLoggedIn(handler func(s *state, cmd command, user database.GetUserRow) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.Config.User)
		if err != nil {
			return err
		}

		return handler(s, cmd, user)
	}
}
