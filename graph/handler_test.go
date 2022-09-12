package graph_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/99designs/gqlgen/graphql"
	"github.com/aereal/github-graphql-proxy/graph"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-github/v47/github"
	"github.com/vektah/gqlparser/v2/gqlerror"
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
	type testCase struct {
		name                string
		responseDefinition  mockAPIResponseList
		graphqlParams       *graphql.RawParams
		wantData            map[string]any
		wantExtension       map[string]any
		assertsErrorMessage func(*testing.T, gqlerror.List)
	}
	testCases := []testCase{
		{
			"ok",
			mockAPIResponseList{
				{
					urlPath: fmt.Sprintf("/api/v3/orgs/%s", org),
					body: &github.Organization{
						Plan: &github.Plan{
							Name:        github.String("enterprise"),
							Seats:       github.Int(5),
							FilledSeats: github.Int(3),
						},
					},
				},
			},
			&graphql.RawParams{
				Query:     query,
				Variables: map[string]any{"org": org},
			},
			map[string]any{"organization": map[string]any{"plan": map[string]any{"filledSeats": float64(3), "seats": float64(5), "name": "enterprise"}}},
			nil,
			func(t *testing.T, errs gqlerror.List) {
				t.Helper()
				if msg := errs.Error(); msg != "" {
					t.Errorf("errors:\n%s", msg)
				}
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			githubClient, finite, err := newMockedGitHubClient(tc.responseDefinition)
			if err != nil {
				t.Fatal(err)
			}
			defer finite()
			resp, err, close := sendGraphqlRequest(context.Background(), tc.graphqlParams, githubClient)
			defer close()
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				t.Fatalf("response status code: %d", resp.StatusCode)
			}
			var gqlResp graphql.Response
			if err := json.NewDecoder(resp.Body).Decode(&gqlResp); err != nil {
				t.Fatalf("cannot decode to GraphQL response: %+v", err)
			}
			tc.assertsErrorMessage(t, gqlResp.Errors)
			var gotData any
			if err := json.Unmarshal(gqlResp.Data, &gotData); err != nil {
				t.Errorf("cannot decode data field: %+v", err)
			} else {
				if diff := cmp.Diff(gotData, tc.wantData); diff != "" {
					t.Errorf("data (-got, +want):\n%s", diff)
				}
			}
			if diff := cmp.Diff(gqlResp.Extensions, tc.wantExtension); diff != "" {
				t.Errorf("extension (-got, +want):\n%s", diff)
			}
		})
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

type mockAPIResponse struct {
	method  string
	urlPath string
	code    int
	body    any
}

var _ http.Handler = (*mockAPIResponse)(nil)

func (a *mockAPIResponse) match(r *http.Request) bool {
	method := a.method
	if method == "" {
		method = http.MethodGet
	}
	return fmt.Sprintf("%s %s", r.Method, r.URL.Path) == fmt.Sprintf("%s %s", method, a.urlPath)
}

func (a *mockAPIResponse) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	b, err := json.Marshal(a.body)
	if err != nil {
		w.Header().Add("content-type", "application/json")
		w.WriteHeader(599)
		fmt.Fprintln(w, `{"error":"cannot encode body as a JSON"}`)
		return
	}
	code := a.code
	if code == 0 {
		code = http.StatusOK
	}
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(code)
	fmt.Fprintln(w, string(b))
}

type mockAPIResponseList []*mockAPIResponse

var _ http.Handler = (mockAPIResponseList)(nil)

func (l mockAPIResponseList) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, def := range l {
		if def.match(r) {
			def.ServeHTTP(w, r)
			return
		}
	}
	noMatchingDefinitionFoundHandler(w, r)
}

var noMatchingDefinitionFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(599)
	fmt.Fprintf(w, `{"error":"no matching definition found"}`)
})

func sendGraphqlRequest(ctx context.Context, params *graphql.RawParams, githubClient *github.Client) (*http.Response, error, func()) {
	handlerSrv := httptest.NewServer(graph.NewHTTPHandler(githubClient))
	close := func() { handlerSrv.Close() }
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(params); err != nil {
		return nil, err, close
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, handlerSrv.URL, buf)
	if err != nil {
		return nil, err, close
	}
	req.Header.Set("content-type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err, close
	}
	return resp, nil, close
}
