package authz

import (
	"log"
	"net/http"
	"strings"

	"golang.org/x/oauth2"
)

const separator = "Bearer "

func ProxiedHTTPClient(r *http.Request) *http.Client {
	_, token, found := strings.Cut(r.Header.Get("authorization"), separator)
	log.Printf("token=%q found=%v", token, found)
	if !found {
		return http.DefaultClient
	}
	ctx := r.Context()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	return oauth2.NewClient(ctx, ts)
}
