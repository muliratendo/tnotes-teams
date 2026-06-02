package middleware

import (
	"github.com/tendo-mulira/tnotes-teams/internal/config"
	"github.com/tendo-mulira/tnotes-teams/internal/repository"
)

// Middleware holds shared dependencies for all middleware.
type Middleware struct {
	cfg   *config.Config
	repos *repository.Repositories
}

// NewMiddleware creates a new Middleware instance.
func NewMiddleware(cfg *config.Config, repos *repository.Repositories) *Middleware {
	return &Middleware{
		cfg:   cfg,
		repos: repos,
	}
}
