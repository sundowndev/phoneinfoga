package suppliers

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (fn roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req)
}

func TestXquikSupplierSearch(t *testing.T) {
	client := &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		assert.Equal(t, http.MethodGet, req.Method)
		assert.Equal(t, "xquik.com", req.URL.Host)
		assert.Equal(t, "/api/v1/x/tweets/search", req.URL.Path)
		assert.Equal(t, `"+14152229670" OR "14152229670"`, req.URL.Query().Get("q"))
		assert.Equal(t, "Latest", req.URL.Query().Get("queryType"))
		assert.Equal(t, "25", req.URL.Query().Get("limit"))
		assert.Equal(t, "secret", req.Header.Get("x-api-key"))
		assert.Equal(t, "application/json", req.Header.Get("Accept"))
		assert.NotContains(t, req.URL.String(), "secret")

		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
			Body: io.NopCloser(strings.NewReader(`{
				"tweets":[{
					"id":"1",
					"text":"public post",
					"createdAt":"2026-08-27T00:00:00Z",
					"likeCount":4,
					"author":{"id":"2","username":"analyst","name":"Analyst"}
				}],
				"has_next_page":true,
				"next_cursor":"cursor"
			}`)),
			Request: req,
		}, nil
	})}

	result, err := NewXquikSupplier(client).Search(
		"secret",
		`"+14152229670" OR "14152229670"`,
		25,
	)

	require.NoError(t, err)
	require.Len(t, result.Tweets, 1)
	assert.Equal(t, "1", result.Tweets[0].ID)
	assert.Equal(t, "Analyst", result.Tweets[0].Author.Name)
	assert.True(t, result.HasNextPage)
	assert.Equal(t, "cursor", result.NextCursor)
}

func TestXquikSupplierRejectsRedirects(t *testing.T) {
	requestCount := 0
	client := &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		requestCount++
		return &http.Response{
			StatusCode: http.StatusFound,
			Header:     http.Header{"Location": []string{"https://example.invalid/capture"}},
			Body:       io.NopCloser(strings.NewReader("redirect")),
			Request:    req,
		}, nil
	})}

	_, err := NewXquikSupplier(client).Search("secret", `"+14152229670"`, 20)

	assert.EqualError(t, err, "Xquik API returned HTTP 302")
	assert.Equal(t, 1, requestCount)
}

func TestXquikSupplierDoesNotExposeErrorBodies(t *testing.T) {
	client := &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusTooManyRequests,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader(`{"message":"private upstream detail"}`)),
			Request:    req,
		}, nil
	})}

	_, err := NewXquikSupplier(client).Search("secret", `"+14152229670"`, 20)

	assert.EqualError(t, err, "Xquik API returned HTTP 429")
	assert.NotContains(t, err.Error(), "private upstream detail")
	assert.NotContains(t, err.Error(), "secret")
}

func TestXquikSupplierRejectsOversizedResponses(t *testing.T) {
	client := &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader(strings.Repeat("x", xquikResponseMaxBytes+1))),
			Request:    req,
		}, nil
	})}

	_, err := NewXquikSupplier(client).Search("secret", `"+14152229670"`, 20)

	assert.EqualError(t, err, "Xquik API response exceeds 4 MiB")
}

func TestXquikSupplierRejectsInvalidResponses(t *testing.T) {
	testcases := []struct {
		name string
		body string
		err  string
	}{
		{name: "invalid JSON", body: "not JSON", err: "Xquik API returned invalid JSON"},
		{name: "missing tweets", body: `{}`, err: "Xquik API response is missing tweets"},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			client := &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Header:     make(http.Header),
					Body:       io.NopCloser(strings.NewReader(tt.body)),
					Request:    req,
				}, nil
			})}

			_, err := NewXquikSupplier(client).Search("secret", `"+14152229670"`, 20)

			assert.EqualError(t, err, tt.err)
		})
	}
}

func TestXquikSupplierValidatesInputsBeforeNetworkAccess(t *testing.T) {
	client := &http.Client{Transport: roundTripFunc(func(_ *http.Request) (*http.Response, error) {
		return nil, errors.New("network should not run")
	})}
	supplier := NewXquikSupplier(client)

	testcases := []struct {
		name  string
		key   string
		query string
		limit int
		err   string
	}{
		{name: "missing key", query: "query", limit: 20, err: "Xquik API key is not defined"},
		{name: "empty query", key: "secret", limit: 20, err: "Xquik search query is empty"},
		{name: "zero limit", key: "secret", query: "query", err: "Xquik result limit must be between 1 and 100"},
		{name: "large limit", key: "secret", query: "query", limit: 101, err: "Xquik result limit must be between 1 and 100"},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			_, err := supplier.Search(tt.key, tt.query, tt.limit)
			assert.EqualError(t, err, tt.err)
		})
	}
}
