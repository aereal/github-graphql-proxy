//go:generate go run github.com/99designs/gqlgen generate

package resolvers

import (
	"github.com/google/go-github/v47/github"
)

func New(githubClient *github.Client) *Resolver {
	return &Resolver{githubClient: githubClient}
}

type Resolver struct {
	githubClient *github.Client
}
