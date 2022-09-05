package graph

import (
	"net/http"

	gqlgenhandler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/aereal/github-graphql-proxy/graph/handler"
	"github.com/aereal/github-graphql-proxy/graph/resolvers"
)

func NewHTTPHandler() http.Handler {
	schema := handler.NewExecutableSchema(handler.Config{Resolvers: &resolvers.Resolver{}})
	h := gqlgenhandler.New(schema)
	h.AddTransport(transport.Options{ /* TODO: AllowedMethods */ })
	h.AddTransport(transport.GET{})
	h.AddTransport(transport.POST{})
	h.SetQueryCache(lru.New(100))
	h.Use(extension.Introspection{})
	return h
}
