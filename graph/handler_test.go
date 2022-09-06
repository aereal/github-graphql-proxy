package graph_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/99designs/gqlgen/graphql"
	"github.com/aereal/github-graphql-proxy/graph"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-github/v47/github"
)

var (
	query = `
		query myQuery($org: String!) {
			organization(login: $org) {
				plan {
					name
					seats
					filledSeats
				}
			}
		}
	`
)

func TestHandler(t *testing.T) {
	org := "test-org"
	githubHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")
		sig := fmt.Sprintf("%s %s", r.Method, r.URL.Path)
		switch sig {
		case fmt.Sprintf("GET /api/v3/orgs/%s", org):
			org := &github.Organization{
				Plan: &github.Plan{
					Name:        github.String("enterprise"),
					Seats:       github.Int(5),
					FilledSeats: github.Int(3),
				},
			}
			if err := json.NewEncoder(w).Encode(org); err != nil {
				t.Error(err)
			}
		default:
			w.WriteHeader(599)
			fmt.Fprintln(w, `{"error":"unhandled request"}`)
			return
		}
	})
	githubClient, finite, err := newMockedGitHubClient(githubHandler)
	if err != nil {
		t.Fatal(err)
	}
	defer finite()
	handlerSrv := httptest.NewServer(graph.NewHTTPHandler(githubClient))
	defer handlerSrv.Close()
	ctx := context.Background()
	params := &graphql.RawParams{
		Query:     query,
		Variables: map[string]any{"org": org},
	}
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(params); err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, handlerSrv.URL, buf)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("content-type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("response status code: %d body=%s", resp.StatusCode, string(b))
	}
	var gqlResp graphql.Response
	if err := json.NewDecoder(resp.Body).Decode(&gqlResp); err != nil {
		t.Fatal(err)
	}
	if errMsg := gqlResp.Errors.Error(); errMsg != "" {
		t.Errorf("errors: %s", errMsg)
	}
	var gotData any
	if err := json.Unmarshal(gqlResp.Data, &gotData); err != nil {
		t.Fatal(err)
	}
	want := map[string]any{"organization": map[string]any{"plan": map[string]any{"filledSeats": float64(3), "seats": float64(5), "name": "enterprise"}}}
	if diff := cmp.Diff(gotData, want); diff != "" {
		t.Errorf("data (-got, +want):\n%s", diff)
	}
}

func newMockedGitHubClient(h http.Handler) (*github.Client, func(), error) {
	srv := httptest.NewServer(h)
	client, err := github.NewEnterpriseClient(srv.URL, srv.URL, srv.Client())
	if err != nil {
		return nil, func() {}, err
	}
	return client, func() { srv.Close() }, nil
}
