package remote

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/sundowndev/phoneinfoga/v2/lib/filter"
	"github.com/sundowndev/phoneinfoga/v2/lib/remote/suppliers"
	"github.com/sundowndev/phoneinfoga/v2/test"
)

type fakeXquikSupplier struct {
	apiKey   string
	query    string
	limit    int
	response *suppliers.XquikSearchResponse
	err      error
	calls    int
}

func (s *fakeXquikSupplier) Search(apiKey string, query string, limit int) (*suppliers.XquikSearchResponse, error) {
	s.calls++
	s.apiKey = apiKey
	s.query = query
	s.limit = limit
	return s.response, s.err
}

func TestXquikScannerMetadata(t *testing.T) {
	scanner := NewXquikScanner(&fakeXquikSupplier{})

	assert.Equal(t, Xquik, scanner.Name())
	assert.NotEmpty(t, scanner.Description())
}

func TestXquikScannerIsRegistered(t *testing.T) {
	library := NewLibrary(filter.NewEngine())

	InitScanners(library)

	assert.NotNil(t, library.GetScanner(Xquik))
}

func TestXquikScannerDryRun(t *testing.T) {
	t.Setenv("XQUIK_API_KEY", "")
	t.Setenv("XQUIK_MAX_RESULTS", "")
	scanner := NewXquikScanner(&fakeXquikSupplier{})
	number := *test.NewFakeUSNumber()

	testcases := []struct {
		name string
		opts ScannerOptions
		err  string
	}{
		{name: "configured", opts: ScannerOptions{"XQUIK_API_KEY": "secret"}},
		{name: "custom limit", opts: ScannerOptions{"XQUIK_API_KEY": "secret", "XQUIK_MAX_RESULTS": "100"}},
		{name: "missing key", opts: ScannerOptions{}, err: "XQUIK_API_KEY is not defined"},
		{name: "empty key", opts: ScannerOptions{"XQUIK_API_KEY": " "}, err: "XQUIK_API_KEY is not defined"},
		{name: "invalid limit", opts: ScannerOptions{"XQUIK_API_KEY": "secret", "XQUIK_MAX_RESULTS": "many"}, err: "XQUIK_MAX_RESULTS must be between 1 and 100"},
		{name: "zero limit", opts: ScannerOptions{"XQUIK_API_KEY": "secret", "XQUIK_MAX_RESULTS": "0"}, err: "XQUIK_MAX_RESULTS must be between 1 and 100"},
		{name: "large limit", opts: ScannerOptions{"XQUIK_API_KEY": "secret", "XQUIK_MAX_RESULTS": "101"}, err: "XQUIK_MAX_RESULTS must be between 1 and 100"},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			err := scanner.DryRun(number, tt.opts)
			if tt.err == "" {
				assert.NoError(t, err)
				return
			}
			assert.EqualError(t, err, tt.err)
		})
	}
}

func TestXquikScannerSearchesExactPhoneFormats(t *testing.T) {
	supplier := &fakeXquikSupplier{response: &suppliers.XquikSearchResponse{
		Tweets: []suppliers.XquikTweet{
			{
				ID:           "1",
				Text:         "Call us\nnow\x1b[31m",
				CreatedAt:    "2026-08-27T00:00:00Z",
				LikeCount:    4,
				RetweetCount: 3,
				ReplyCount:   2,
				QuoteCount:   1,
				ViewCount:    100,
				URL:          "https://x.com/example/status/1",
				Lang:         "en",
				Author: suppliers.XquikAuthor{
					ID:          "2",
					Username:    "example",
					Name:        "Example\u202e name",
					Description: "Public profile\ttext",
					Location:    "Vienna",
					Followers:   200,
					Verified:    true,
				},
			},
		},
		HasNextPage: true,
	}}
	scanner := NewXquikScanner(supplier)

	result, err := scanner.Run(*test.NewFakeUSNumber(), ScannerOptions{
		"XQUIK_API_KEY":     "secret",
		"XQUIK_MAX_RESULTS": "25",
	})

	require.NoError(t, err)
	assert.Equal(t, 1, supplier.calls)
	assert.Equal(t, "secret", supplier.apiKey)
	assert.Equal(t, `"+14152229670" OR "14152229670" OR "4152229670" OR "(415) 222-9670"`, supplier.query)
	assert.Equal(t, 25, supplier.limit)
	assert.Equal(t, XquikScannerResponse{
		ResultCount: 1,
		HasMore:     true,
		Tweets: []XquikTweet{
			{
				ID:           "1",
				Text:         "Call us now",
				CreatedAt:    "2026-08-27T00:00:00Z",
				LikeCount:    4,
				RetweetCount: 3,
				ReplyCount:   2,
				QuoteCount:   1,
				ViewCount:    100,
				URL:          "https://x.com/example/status/1",
				Lang:         "en",
				Author: XquikAuthor{
					ID:          "2",
					Username:    "example",
					Name:        "Example name",
					Description: "Public profile text",
					Location:    "Vienna",
					Followers:   200,
					Verified:    true,
				},
			},
		},
	}, result)
}

func TestXquikScannerEnforcesItsResultLimit(t *testing.T) {
	supplier := &fakeXquikSupplier{response: &suppliers.XquikSearchResponse{
		Tweets: []suppliers.XquikTweet{{ID: "1"}, {ID: "2"}, {ID: "3"}},
	}}
	scanner := NewXquikScanner(supplier)

	result, err := scanner.Run(*test.NewFakeUSNumber(), ScannerOptions{
		"XQUIK_API_KEY":     "secret",
		"XQUIK_MAX_RESULTS": "2",
	})

	require.NoError(t, err)
	response := result.(XquikScannerResponse)
	assert.Equal(t, 2, response.ResultCount)
	assert.True(t, response.HasMore)
	assert.Equal(t, []XquikTweet{{ID: "1"}, {ID: "2"}}, response.Tweets)
}

func TestXquikScannerReturnsConfigurationErrorsWithoutCallingTheSupplier(t *testing.T) {
	t.Setenv("XQUIK_API_KEY", "")
	supplier := &fakeXquikSupplier{}
	scanner := NewXquikScanner(supplier)

	_, err := scanner.Run(*test.NewFakeUSNumber(), ScannerOptions{})

	assert.EqualError(t, err, "XQUIK_API_KEY is not defined")
	assert.Equal(t, 0, supplier.calls)
}

func TestXquikScannerReturnsSupplierErrors(t *testing.T) {
	supplier := &fakeXquikSupplier{err: errors.New("request failed")}
	scanner := NewXquikScanner(supplier)

	_, err := scanner.Run(*test.NewFakeUSNumber(), ScannerOptions{"XQUIK_API_KEY": "secret"})

	assert.EqualError(t, err, "request failed")
}
