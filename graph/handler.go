package graph

import (
	"net/http"

	gqlgenhandler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/aereal/github-graphql-proxy/authz"
	"github.com/aereal/github-graphql-proxy/graph/handler"
	"github.com/aereal/github-graphql-proxy/graph/resolvers"
	"github.com/google/go-github/v47/github"
)

func NewHTTPHandler() http.Handler {
	cache := lru.New(100)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		githubClient := github.NewClient(authz.ProxiedHTTPClient(r))
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
