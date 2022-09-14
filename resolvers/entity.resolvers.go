package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"strings"

	githubgraphqlproxy "github.com/aereal/github-graphql-proxy"
)

// FindOrganizationByLogin is the resolver for the findOrganizationByLogin field.
func (r *entityResolver) FindOrganizationByLogin(ctx context.Context, login string) (*githubgraphqlproxy.Organization, error) {
	return &githubgraphqlproxy.Organization{Login: login}, nil
}

// FindRepositoryByNameWithOwner is the resolver for the findRepositoryByNameWithOwner field.
func (r *entityResolver) FindRepositoryByNameWithOwner(ctx context.Context, nameWithOwner string) (*githubgraphqlproxy.Repository, error) {
	parts := strings.Split(nameWithOwner, "/")
	return &githubgraphqlproxy.Repository{
		Owner:         parts[0],
		Name:          parts[1],
		NameWithOwner: nameWithOwner,
	}, nil
}

// Entity returns githubgraphqlproxy.EntityResolver implementation.
func (r *Resolver) Entity() githubgraphqlproxy.EntityResolver { return &entityResolver{r} }

type entityResolver struct{ *Resolver }
