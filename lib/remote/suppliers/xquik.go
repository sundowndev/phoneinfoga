package suppliers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	xquikSearchURL        = "https://xquik.com/api/v1/x/tweets/search"
	xquikResponseMaxBytes = 4 << 20
)

// XquikSupplierInterface searches public posts through a bounded API call.
type XquikSupplierInterface interface {
	Search(string, string, int) (*XquikSearchResponse, error)
}

// XquikAuthor contains the public profile fields used by the scanner.
type XquikAuthor struct {
	ID          string `json:"id"`
	Username    string `json:"username"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Location    string `json:"location"`
	Followers   int64  `json:"followers"`
	Verified    bool   `json:"verified"`
}

// XquikTweet contains the public post fields used by the scanner.
type XquikTweet struct {
	ID           string      `json:"id"`
	Text         string      `json:"text"`
	CreatedAt    string      `json:"createdAt"`
	LikeCount    int64       `json:"likeCount"`
	RetweetCount int64       `json:"retweetCount"`
	ReplyCount   int64       `json:"replyCount"`
	QuoteCount   int64       `json:"quoteCount"`
	ViewCount    int64       `json:"viewCount"`
	URL          string      `json:"url"`
	Lang         string      `json:"lang"`
	Author       XquikAuthor `json:"author"`
}

// XquikSearchResponse is the bounded Tweet Search response.
type XquikSearchResponse struct {
	Tweets      []XquikTweet `json:"tweets"`
	HasNextPage bool         `json:"has_next_page"`
	NextCursor  string       `json:"next_cursor"`
}

// XquikSupplier requests the fixed Xquik Tweet Search endpoint.
type XquikSupplier struct {
	client *http.Client
}

// NewXquikSupplier creates a supplier with redirects disabled and a bounded timeout.
func NewXquikSupplier(client *http.Client) *XquikSupplier {
	if client == nil {
		client = http.DefaultClient
	}

	safeClient := *client
	safeClient.CheckRedirect = func(_ *http.Request, _ []*http.Request) error {
		return http.ErrUseLastResponse
	}
	if safeClient.Timeout == 0 {
		safeClient.Timeout = 30 * time.Second
	}

	return &XquikSupplier{client: &safeClient}
}

// Search returns one bounded page of public posts for an exact query.
func (s *XquikSupplier) Search(apiKey string, query string, limit int) (*XquikSearchResponse, error) {
	if strings.TrimSpace(apiKey) == "" {
		return nil, errors.New("Xquik API key is not defined")
	}
	if strings.TrimSpace(query) == "" {
		return nil, errors.New("Xquik search query is empty")
	}
	if limit < 1 || limit > 100 {
		return nil, errors.New("Xquik result limit must be between 1 and 100")
	}

	endpoint, err := url.Parse(xquikSearchURL)
	if err != nil {
		return nil, err
	}
	params := endpoint.Query()
	params.Set("q", query)
	params.Set("queryType", "Latest")
	params.Set("limit", fmt.Sprintf("%d", limit))
	endpoint.RawQuery = params.Encode()

	req, err := http.NewRequest(http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("x-api-key", apiKey)

	response, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("Xquik API returned HTTP %d", response.StatusCode)
	}

	body, err := io.ReadAll(io.LimitReader(response.Body, xquikResponseMaxBytes+1))
	if err != nil {
		return nil, err
	}
	if len(body) > xquikResponseMaxBytes {
		return nil, errors.New("Xquik API response exceeds 4 MiB")
	}

	var result XquikSearchResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, errors.New("Xquik API returned invalid JSON")
	}
	if result.Tweets == nil {
		return nil, errors.New("Xquik API response is missing tweets")
	}

	return &result, nil
}
