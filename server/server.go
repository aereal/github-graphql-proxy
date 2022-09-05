package server

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/aereal/github-graphql-proxy/graph"
)

func Handler() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/", playground.Handler("GraphQL playground", "/query"))
	mux.Handle("/query", graph.NewHTTPHandler(runtime.GOMAXPROCS(0)))
	return mux
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
