package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	githubgraphqlproxy "github.com/aereal/github-graphql-proxy"
	"github.com/aereal/github-graphql-proxy/authz"
	"github.com/aereal/github-graphql-proxy/resolvers"
	"github.com/google/go-github/v47/github"
	"golang.org/x/sync/semaphore"
)

func Handler() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/", playground.Handler("GraphQL playground", "/extension/query"))
	mux.Handle("/extension/query", withSemaphoreClient(int64(runtime.GOMAXPROCS(0))))
	return mux
}

func withSemaphoreClient(maxConcurrency int64) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		httpClient := authz.ProxiedHTTPClient(r.Context(), r.Header.Get("authorization"))
		rt := &semaphoreTransport{
			base: httpClient.Transport,
			sem:  semaphore.NewWeighted(maxConcurrency),
		}
		if rt.base == nil {
			rt.base = http.DefaultTransport
		}
		httpClient.Transport = rt
		h := queryHandler(github.NewClient(httpClient))
		h.ServeHTTP(w, r)
	})
}

func Start(ctx context.Context, addr string, startTimeout time.Duration) error {
	srv := &http.Server{
		Handler: Handler(),
		Addr:    addr,
	}
	go graceful(ctx, srv, startTimeout)
	log.Printf("starting server (addr=%s) ...", addr)
	err := srv.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func graceful(ctx context.Context, srv *http.Server, timeout time.Duration) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	sig := <-sigChan
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	log.Printf("shutting down server (%v) ...", sig)
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("failed to shutdown: %v", err)
	}
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

func queryHandler(githubClient *github.Client) http.Handler {
	schema := githubgraphqlproxy.NewExecutableSchema(githubgraphqlproxy.Config{Resolvers: resolvers.New(githubClient)})
	h := handler.New(schema)
	h.AddTransport(transport.Options{ /* TODO: AllowedMethods */ })
	h.AddTransport(transport.GET{})
	h.AddTransport(transport.POST{})
	h.Use(extension.Introspection{})
	return h
}
