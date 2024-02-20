package remote

import (
	"context"
	"errors"
	"fmt"
	"github.com/sundowndev/dorkgen"
	"github.com/sundowndev/dorkgen/googlesearch"
	"github.com/sundowndev/phoneinfoga/v2/lib/number"
	"google.golang.org/api/customsearch/v1"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
	"net/http"
	"os"
	"strconv"
)

const GoogleCSE = "googlecse"

type googleCSEScanner struct {
	MaxResults int64
	httpClient *http.Client
}

type ResultItem struct {
	Title string `json:"title,omitempty" console:"Title,omitempty"`
	URL   string `json:"url,omitempty" console:"URL,omitempty"`
}

type GoogleCSEScannerResponse struct {
	Homepage          string       `json:"homepage,omitempty" console:"Homepage,omitempty"`
	ResultCount       int          `json:"result_count" console:"Results shown"`
	TotalResultCount  int          `json:"total_result_count" console:"Total number of results"`
	TotalRequestCount int          `json:"total_request_count" console:"Requests made"`
	Items             []ResultItem `json:"items,omitempty" console:"Items,omitempty"`
}

func NewGoogleCSEScanner(HTTPclient *http.Client) Scanner {
	// CSE limits you to 10 pages of results with max 10 results per page
	// We only fetch the first page of results by default for each request
	maxResults := 10
	if v := os.Getenv("GOOGLECSE_MAX_RESULTS"); v != "" {
		val, err := strconv.Atoi(v)
		if err == nil {
			if val > 100 {
				val = 100
			}
			maxResults = val
		}
	}

	return &googleCSEScanner{
		MaxResults: int64(maxResults),
		httpClient: HTTPclient,
	}
}

func (s *googleCSEScanner) Name() string {
	return GoogleCSE
}

func (s *googleCSEScanner) Description() string {
	return "Googlecse searches for footprints of a given phone number on the web using Google Custom Search Engine."
}

func (s *googleCSEScanner) DryRun(_ number.Number, opts ScannerOptions) error {
	if opts.GetStringEnv("GOOGLECSE_CX") == "" || opts.GetStringEnv("GOOGLE_API_KEY") == "" {
		return errors.New("search engine ID and/or API key is not defined")
	}
	return nil
}

func (s *googleCSEScanner) Run(n number.Number, opts ScannerOptions) (interface{}, error) {
	var allItems []*customsearch.Result
	var dorks []*GoogleSearchDork
	var totalResultCount int
	var totalRequestCount int
	var cx = opts.GetStringEnv("GOOGLECSE_CX")
	var apikey = opts.GetStringEnv("GOOGLE_API_KEY")

	dorks = append(dorks, s.generateDorkQueries(n)...)

	customsearchService, err := customsearch.NewService(
		context.Background(),
		option.WithAPIKey(apikey),
		option.WithHTTPClient(s.httpClient),
	)
	if err != nil {
		return nil, err
	}

	for _, req := range dorks {
		n, items, err := s.search(customsearchService, req.Dork, cx)
		if err != nil {
			if s.isRateLimit(err) {
				return nil, errors.New("rate limit exceeded, see https://developers.google.com/custom-search/v1/overview#pricing")
			}
			return nil, err
		}
		allItems = append(allItems, items...)
		totalResultCount += n
		totalRequestCount++
	}

	var data GoogleCSEScannerResponse
	for _, item := range allItems {
		data.Items = append(data.Items, ResultItem{
			Title: item.Title,
			URL:   item.Link,
		})
	}
	data.Homepage = fmt.Sprintf("https://cse.google.com/cse?cx=%s", cx)
	data.ResultCount = len(allItems)
	data.TotalResultCount = totalResultCount
	data.TotalRequestCount = totalRequestCount

	return data, nil
}

func (s *googleCSEScanner) search(service *customsearch.Service, q string, cx string) (int, []*customsearch.Result, error) {
	var results []*customsearch.Result
	var totalResultCount int

	offset := int64(0)
	for offset < s.MaxResults {
		search := service.Cse.List()
		search.Cx(cx)
		search.Q(q)
		search.Start(offset)
		searchQuery, err := search.Do()
		if err != nil {
			return 0, nil, err
		}
		results = append(results, searchQuery.Items...)
		totalResultCount, err = strconv.Atoi(searchQuery.SearchInformation.TotalResults)
		if err != nil {
			return 0, nil, err
		}
		if totalResultCount <= int(s.MaxResults) {
			break
		}
		offset += int64(len(searchQuery.Items))
	}

	return totalResultCount, results, nil
}

func (s *googleCSEScanner) isRateLimit(theError error) bool {
	if theError == nil {
		return false
	}
	var err *googleapi.Error
	if !errors.As(theError, &err) {
		return false
	}
	if theError.(*googleapi.Error).Code != 429 {
		return false
	}
	return true
}

func (s *googleCSEScanner) generateDorkQueries(number number.Number) (results []*GoogleSearchDork) {
	var dorks = []*googlesearch.GoogleSearch{
		dorkgen.NewGoogleSearch().
			InText(number.International).
			Or().
			InText(number.E164).
			Or().
			InText(number.RawLocal).
			Or().
			InText(number.Local),
		dorkgen.NewGoogleSearch().
			Group(dorkgen.NewGoogleSearch().
				Ext("doc").
				Or().
				Ext("docx").
				Or().
				Ext("odt").
				Or().
				Ext("pdf").
				Or().
				Ext("rtf").
				Or().
				Ext("sxw").
				Or().
				Ext("psw").
				Or().
				Ext("ppt").
				Or().
				Ext("pptx").
				Or().
				Ext("pps").
				Or().
				Ext("csv").
				Or().
				Ext("txt").
				Or().
				Ext("xls")).
			InText(number.International).
			Or().
			InText(number.E164).
			Or().
			InText(number.RawLocal).
			Or().
			InText(number.Local),
	}

	for _, dork := range dorks {
		results = append(results, &GoogleSearchDork{
			Number: number.E164,
			Dork:   dork.String(),
			URL:    dork.URL(),
		})
	}

	return results
}
