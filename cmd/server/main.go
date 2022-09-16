package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aereal/github-graphql-proxy/server"
)

func main() {
	srv := server.New(server.WithPort(os.Getenv("PORT")), server.WithStartTimeoutLiteral("START_TIMEOUT"))
	ctx := context.Background()
	if err := srv.Start(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
}
