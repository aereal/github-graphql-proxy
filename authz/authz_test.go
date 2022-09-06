package authz_test

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aereal/github-graphql-proxy/authz"
)

func TestProxiedHTTPClient(t *testing.T) {
	type testCase struct {
		name            string
		sentAuthzHeader string
		wantAuthzHeader string
	}
	testCases := []testCase{
		{"Bearer", "Bearer 0xdeadbeaf", "Bearer 0xdeadbeaf"},
		{"basic", fmt.Sprintf("Basic %s", base64.URLEncoding.EncodeToString([]byte("admin:pass"))), ""},
		{"no header", "", ""},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var gotHeader string
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				gotHeader = r.Header.Get("authorization")
			}))
			defer srv.Close()
			ctx := context.Background()
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, srv.URL, nil)
			if err != nil {
				t.Fatal(err)
			}
			httpClient := authz.ProxiedHTTPClient(ctx, tc.sentAuthzHeader)
			resp, err := httpClient.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			if gotHeader != tc.wantAuthzHeader {
				t.Errorf("authorization header: got=%q want=%q", tc.wantAuthzHeader, gotHeader)
			}
			if resp.StatusCode != http.StatusOK {
				t.Errorf("[bug] response status code: got=%d", resp.StatusCode)
			}
		})
	}
}
