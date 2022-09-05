package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/aereal/github-graphql-proxy/server"
)

var (
	addr         string
	startTimeout time.Duration
)

func init() {
	flag.StringVar(&addr, "addr", ":8080", "server listening address")
	flag.DurationVar(&startTimeout, "start-timeout", time.Second*5, "timeout to wait server spin-up")
}

func main() {
	flag.Parse()
	ctx := context.Background()
	if err := server.Start(ctx, addr, startTimeout); err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
}
