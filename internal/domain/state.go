package domain

import (
	"github.com/ireoluwa12345/gator/internal/config"
	"github.com/ireoluwa12345/gator/internal/database"
)

type State struct {
	DB     *database.Queries
	Config *config.Config
}
