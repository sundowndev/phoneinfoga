package remote

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/sundowndev/phoneinfoga/v2/lib/number"
	"github.com/sundowndev/phoneinfoga/v2/lib/remote/suppliers"
)

var xquikANSISequence = regexp.MustCompile(`\x1b\[[0-?]*[ -/]*[@-~]`)

const (
	Xquik                  = "xquik"
	xquikDefaultMaxResults = 20
	xquikMaxResults        = 100
)

type xquikScanner struct {
	supplier suppliers.XquikSupplierInterface
}

// XquikAuthor is the public author information returned by the scanner.
type XquikAuthor struct {
	ID          string `json:"id,omitempty" console:"ID,omitempty"`
	Username    string `json:"username,omitempty" console:"Username,omitempty"`
	Name        string `json:"name,omitempty" console:"Name,omitempty"`
	Description string `json:"description,omitempty" console:"Description,omitempty"`
	Location    string `json:"location,omitempty" console:"Location,omitempty"`
	Followers   int64  `json:"followers,omitempty" console:"Followers,omitempty"`
	Verified    bool   `json:"verified" console:"Verified"`
}

// XquikTweet is one sanitized public post returned by the scanner.
type XquikTweet struct {
	ID           string      `json:"id" console:"ID"`
	Text         string      `json:"text" console:"Text"`
	CreatedAt    string      `json:"created_at,omitempty" console:"Created at,omitempty"`
	LikeCount    int64       `json:"like_count" console:"Likes"`
	RetweetCount int64       `json:"retweet_count" console:"Reposts"`
	ReplyCount   int64       `json:"reply_count" console:"Replies"`
	QuoteCount   int64       `json:"quote_count" console:"Quotes"`
	ViewCount    int64       `json:"view_count" console:"Views"`
	URL          string      `json:"url,omitempty" console:"URL,omitempty"`
	Lang         string      `json:"lang,omitempty" console:"Language,omitempty"`
	Author       XquikAuthor `json:"author" console:"Author"`
}

// XquikScannerResponse is the bounded result set returned by one scan.
type XquikScannerResponse struct {
	ResultCount int          `json:"result_count" console:"Results shown"`
	HasMore     bool         `json:"has_more" console:"More results available"`
	Tweets      []XquikTweet `json:"tweets" console:"Tweets,omitempty"`
}

// NewXquikScanner creates the optional public X footprint scanner.
func NewXquikScanner(supplier suppliers.XquikSupplierInterface) Scanner {
	if supplier == nil {
		supplier = suppliers.NewXquikSupplier(nil)
	}
	return &xquikScanner{supplier: supplier}
}

func (s *xquikScanner) Name() string {
	return Xquik
}

func (s *xquikScanner) Description() string {
	return "Search public X posts for exact forms of a phone number through Xquik."
}

func (s *xquikScanner) DryRun(_ number.Number, opts ScannerOptions) error {
	_, _, err := xquikConfig(opts)
	return err
}

func (s *xquikScanner) Run(n number.Number, opts ScannerOptions) (interface{}, error) {
	apiKey, limit, err := xquikConfig(opts)
	if err != nil {
		return nil, err
	}

	response, err := s.supplier.Search(
		apiKey,
		xquikPhoneQuery(n),
		limit,
	)
	if err != nil {
		return nil, err
	}

	tweetCount := len(response.Tweets)
	if tweetCount > limit {
		tweetCount = limit
	}
	tweets := make([]XquikTweet, 0, tweetCount)
	for _, tweet := range response.Tweets[:tweetCount] {
		tweets = append(tweets, XquikTweet{
			ID:           safeXquikText(tweet.ID),
			Text:         safeXquikText(tweet.Text),
			CreatedAt:    safeXquikText(tweet.CreatedAt),
			LikeCount:    tweet.LikeCount,
			RetweetCount: tweet.RetweetCount,
			ReplyCount:   tweet.ReplyCount,
			QuoteCount:   tweet.QuoteCount,
			ViewCount:    tweet.ViewCount,
			URL:          safeXquikText(tweet.URL),
			Lang:         safeXquikText(tweet.Lang),
			Author: XquikAuthor{
				ID:          safeXquikText(tweet.Author.ID),
				Username:    safeXquikText(tweet.Author.Username),
				Name:        safeXquikText(tweet.Author.Name),
				Description: safeXquikText(tweet.Author.Description),
				Location:    safeXquikText(tweet.Author.Location),
				Followers:   tweet.Author.Followers,
				Verified:    tweet.Author.Verified,
			},
		})
	}

	return XquikScannerResponse{
		ResultCount: len(tweets),
		HasMore:     response.HasNextPage || len(response.Tweets) > limit,
		Tweets:      tweets,
	}, nil
}

func xquikConfig(opts ScannerOptions) (string, int, error) {
	apiKey := strings.TrimSpace(opts.GetStringEnv("XQUIK_API_KEY"))
	if apiKey == "" {
		return "", 0, errors.New("XQUIK_API_KEY is not defined")
	}
	limit, err := xquikResultLimit(opts)
	return apiKey, limit, err
}

func xquikResultLimit(opts ScannerOptions) (int, error) {
	raw := strings.TrimSpace(opts.GetStringEnv("XQUIK_MAX_RESULTS"))
	if raw == "" {
		return xquikDefaultMaxResults, nil
	}
	limit, err := strconv.Atoi(raw)
	if err != nil || limit < 1 || limit > xquikMaxResults {
		return 0, fmt.Errorf("XQUIK_MAX_RESULTS must be between 1 and %d", xquikMaxResults)
	}
	return limit, nil
}

func xquikPhoneQuery(n number.Number) string {
	seen := map[string]bool{}
	parts := make([]string, 0, 4)
	for _, value := range []string{n.E164, n.International, n.RawLocal, n.Local} {
		value = strings.TrimSpace(value)
		if value == "" || seen[value] {
			continue
		}
		seen[value] = true
		escaped := strings.NewReplacer(`\`, `\\`, `"`, `\"`).Replace(value)
		parts = append(parts, `"`+escaped+`"`)
	}
	return strings.Join(parts, " OR ")
}

func safeXquikText(value string) string {
	value = xquikANSISequence.ReplaceAllString(value, "")
	value = strings.Map(func(char rune) rune {
		if unicode.IsControl(char) || unicode.In(char, unicode.Cf) {
			return ' '
		}
		return char
	}, value)
	return strings.Join(strings.Fields(value), " ")
}
