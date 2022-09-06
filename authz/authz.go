package authz

import (
	"context"
	"net/http"
	"strings"

	"golang.org/x/oauth2"
)

const separator = "Bearer "

func ProxiedHTTPClient(ctx context.Context, authzHeader string) *http.Client {
	_, token, found := strings.Cut(authzHeader, separator)
	if !found {
		return http.DefaultClient
	}
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	return oauth2.NewClient(ctx, ts)
}
