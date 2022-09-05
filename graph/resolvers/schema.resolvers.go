package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/aereal/github-graphql-proxy/graph/dto"
	"github.com/aereal/github-graphql-proxy/graph/handler"
)

// Organization is the resolver for the organization field.
func (r *queryResolver) Organization(ctx context.Context, login string) (*dto.Organization, error) {
	return &dto.Organization{Login: login}, nil
}

// Query returns handler.QueryResolver implementation.
func (r *Resolver) Query() handler.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
