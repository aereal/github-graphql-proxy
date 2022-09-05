package graph

import (
	"fmt"
	"net/http"

	gqlgenhandler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/aereal/github-graphql-proxy/authz"
	"github.com/aereal/github-graphql-proxy/graph/handler"
	"github.com/aereal/github-graphql-proxy/graph/resolvers"
	"github.com/google/go-github/v47/github"
	"golang.org/x/sync/semaphore"
)

func NewHTTPHandler(maxConcurrency int) http.Handler {
	cache := lru.New(100)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		httpClient := authz.ProxiedHTTPClient(r)
		httpClient.Transport = &semaphoreTransport{
			base: httpClient.Transport,
			sem:  semaphore.NewWeighted(int64(maxConcurrency)),
		}
		githubClient := github.NewClient(httpClient)
		schema := handler.NewExecutableSchema(handler.Config{Resolvers: resolvers.New(githubClient)})
		h := gqlgenhandler.New(schema)
		h.AddTransport(transport.Options{ /* TODO: AllowedMethods */ })
		h.AddTransport(transport.GET{})
		h.AddTransport(transport.POST{})
		h.SetQueryCache(cache)
		h.Use(extension.Introspection{})
		h.ServeHTTP(w, r)
	})
}

type semaphoreTransport struct {
	base http.RoundTripper
	sem  *semaphore.Weighted
}

var _ http.RoundTripper = (*semaphoreTransport)(nil)

func (t *semaphoreTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	ctx := r.Context()
	if err := t.sem.Acquire(ctx, 1); err != nil {
		return nil, fmt.Errorf("request cancelled: %w", err)
	}
	defer t.sem.Release(1)
	return t.base.RoundTrip(r)
}
