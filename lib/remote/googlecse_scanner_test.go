package remote

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/sundowndev/phoneinfoga/v2/lib/filter"
	"github.com/sundowndev/phoneinfoga/v2/lib/number"
	"github.com/sundowndev/phoneinfoga/v2/test"
	"google.golang.org/api/customsearch/v1"
	"google.golang.org/api/googleapi"
	"gopkg.in/h2non/gock.v1"
	"net/http"
	"os"
	"testing"
)

func TestGoogleCSEScanner_Metadata(t *testing.T) {
	scanner := NewGoogleCSEScanner(&http.Client{})
	assert.Equal(t, GoogleCSE, scanner.Name())
	assert.NotEmpty(t, scanner.Description())
}

func TestGoogleCSEScanner_Scan_Success(t *testing.T) {
	testcases := []struct {
		name       string
		number     *number.Number
		opts       ScannerOptions
		expected   map[string]interface{}
		wantErrors map[string]error
		mocks      func()
	}{
		{
			name:   "test with no results",
			number: test.NewFakeUSNumber(),
			expected: map[string]interface{}{
				"googlecse": GoogleCSEScannerResponse{
					Homepage:          "https://cse.google.com/cse?cx=fake_search_engine_id",
					ResultCount:       0,
					TotalResultCount:  0,
					TotalRequestCount: 2,
					Items:             nil,
				},
			},
			wantErrors: map[string]error{},
			mocks: func() {
				gock.New("https://customsearch.googleapis.com").
					Get("/customsearch/v1").
					MatchParam("cx", "fake_search_engine_id").
					// TODO: the matcher below doesn't work for some reason
					//MatchParam("q", "intext:\"14152229670\" OR intext:\"+14152229670\" OR intext:\"4152229670\" OR intext:\"(415) 222-9670\"").
					MatchParam("start", "0").
					Reply(200).
					JSON(&customsearch.Search{
						ServerResponse: googleapi.ServerResponse{
							Header:         http.Header{},
							HTTPStatusCode: 200,
						},
						SearchInformation: &customsearch.SearchSearchInformation{
							FormattedSearchTime:   "0",
							FormattedTotalResults: "0",
							SearchTime:            0,
							TotalResults:          "0",
							ForceSendFields:       nil,
							NullFields:            nil,
						},
						Items: []*customsearch.Result{},
					})

				gock.New("https://customsearch.googleapis.com").
					Get("/customsearch/v1").
					MatchParam("cx", "fake_search_engine_id").
					// TODO: the matcher below doesn't work for some reason
					//MatchParam("q", "(ext:doc OR ext:docx OR ext:odt OR ext:pdf OR ext:rtf OR ext:sxw OR ext:psw OR ext:ppt OR ext:pptx OR ext:pps OR ext:csv OR ext:txt OR ext:xls) intext:\"14152229670\" OR intext:\"+14152229670\" OR intext:\"4152229670\" OR intext:\"(415)+222-9670\"").
					MatchParam("start", "0").
					Reply(200).
					JSON(&customsearch.Search{
						ServerResponse: googleapi.ServerResponse{
							Header:         http.Header{},
							HTTPStatusCode: 200,
						},
						SearchInformation: &customsearch.SearchSearchInformation{
							FormattedSearchTime:   "0",
							FormattedTotalResults: "0",
							SearchTime:            0,
							TotalResults:          "0",
							ForceSendFields:       nil,
							NullFields:            nil,
						},
						Items: []*customsearch.Result{},
					})
			},
		},
		{
			name:   "test with options and no results",
			number: test.NewFakeUSNumber(),
			opts: ScannerOptions{
				"GOOGLECSE_CX":   "custom_cx",
				"GOOGLE_API_KEY": "secret",
			},
			expected: map[string]interface{}{
				"googlecse": GoogleCSEScannerResponse{
					Homepage:          "https://cse.google.com/cse?cx=custom_cx",
					ResultCount:       0,
					TotalResultCount:  0,
					TotalRequestCount: 2,
					Items:             nil,
				},
			},
			wantErrors: map[string]error{},
			mocks: func() {
				gock.New("https://customsearch.googleapis.com").
					Get("/customsearch/v1").
					MatchParam("cx", "custom_cx").
					// TODO: ensure that custom api key is used
					// MatchHeader("Authorization", "secret").
					// TODO: the matcher below doesn't work for some reason
					//MatchParam("q", "intext:\"14152229670\" OR intext:\"+14152229670\" OR intext:\"4152229670\" OR intext:\"(415) 222-9670\"").
					MatchParam("start", "0").
					Reply(200).
					JSON(&customsearch.Search{
						ServerResponse: googleapi.ServerResponse{
							Header:         http.Header{},
							HTTPStatusCode: 200,
						},
						SearchInformation: &customsearch.SearchSearchInformation{
							FormattedSearchTime:   "0",
							FormattedTotalResults: "0",
							SearchTime:            0,
							TotalResults:          "0",
							ForceSendFields:       nil,
							NullFields:            nil,
						},
						Items: []*customsearch.Result{},
					})

				gock.New("https://customsearch.googleapis.com").
					Get("/customsearch/v1").
					MatchParam("cx", "custom_cx").
					// TODO: ensure that custom api key is used
					// MatchHeader("Authorization", "secret").
					// TODO: the matcher below doesn't work for some reason
					//MatchParam("q", "(ext:doc OR ext:docx OR ext:odt OR ext:pdf OR ext:rtf OR ext:sxw OR ext:psw OR ext:ppt OR ext:pptx OR ext:pps OR ext:csv OR ext:txt OR ext:xls) intext:\"14152229670\" OR intext:\"+14152229670\" OR intext:\"4152229670\" OR intext:\"(415)+222-9670\"").
					MatchParam("start", "0").
					Reply(200).
					JSON(&customsearch.Search{
						ServerResponse: googleapi.ServerResponse{
							Header:         http.Header{},
							HTTPStatusCode: 200,
						},
						SearchInformation: &customsearch.SearchSearchInformation{
							FormattedSearchTime:   "0",
							FormattedTotalResults: "0",
							SearchTime:            0,
							TotalResults:          "0",
							ForceSendFields:       nil,
							NullFields:            nil,
						},
						Items: []*customsearch.Result{},
					})
			},
		},
		{
			name:   "test with results",
			number: test.NewFakeUSNumber(),
			expected: map[string]interface{}{
				"googlecse": GoogleCSEScannerResponse{
					Homepage:          "https://cse.google.com/cse?cx=fake_search_engine_id",
					ResultCount:       2,
					TotalResultCount:  2,
					TotalRequestCount: 2,
					Items: []ResultItem{
						{
							Title: "Result 1",
							URL:   "https://result1.com",
						},
						{
							Title: "Result 2",
							URL:   "https://result2.com",
						},
					},
				},
			},
			wantErrors: map[string]error{},
			mocks: func() {
				gock.New("https://customsearch.googleapis.com").
					Get("/customsearch/v1").
					MatchParam("cx", "fake_search_engine_id").
					// TODO: the matcher below doesn't work for some reason
					//MatchParam("q", "intext:\"14152229670\" OR intext:\"+14152229670\" OR intext:\"4152229670\" OR intext:\"(415) 222-9670\"").
					MatchParam("start", "0").
					Reply(200).
					JSON(&customsearch.Search{
						ServerResponse: googleapi.ServerResponse{
							Header:         http.Header{},
							HTTPStatusCode: 200,
						},
						SearchInformation: &customsearch.SearchSearchInformation{
							FormattedSearchTime:   "0",
							FormattedTotalResults: "2",
							SearchTime:            0,
							TotalResults:          "2",
							ForceSendFields:       nil,
							NullFields:            nil,
						},
						Items: []*customsearch.Result{
							{
								Title: "Result 1",
								Link:  "https://result1.com",
							},
							{
								Title: "Result 2",
								Link:  "https://result2.com",
							},
						},
					})

				gock.New("https://customsearch.googleapis.com").
					Get("/customsearch/v1").
					MatchParam("cx", "fake_search_engine_id").
					// TODO: the matcher below doesn't work for some reason
					//MatchParam("q", "(ext:doc OR ext:docx OR ext:odt OR ext:pdf OR ext:rtf OR ext:sxw OR ext:psw OR ext:ppt OR ext:pptx OR ext:pps OR ext:csv OR ext:txt OR ext:xls) intext:\"14152229670\" OR intext:\"+14152229670\" OR intext:\"4152229670\" OR intext:\"(415)+222-9670\"").
					MatchParam("start", "0").
					Reply(200).
					JSON(&customsearch.Search{
						ServerResponse: googleapi.ServerResponse{
							Header:         http.Header{},
							HTTPStatusCode: 200,
						},
						SearchInformation: &customsearch.SearchSearchInformation{
							FormattedSearchTime:   "0",
							FormattedTotalResults: "0",
							SearchTime:            0,
							TotalResults:          "0",
							ForceSendFields:       nil,
							NullFields:            nil,
						},
						Items: []*customsearch.Result{},
					})
			},
		},
		{
			name:     "test with rate limit error",
			number:   test.NewFakeUSNumber(),
			expected: map[string]interface{}{},
			wantErrors: map[string]error{
				"googlecse": errors.New("rate limit exceeded, see https://developers.google.com/custom-search/v1/overview#pricing"),
			},
			mocks: func() {
				gock.New("https://customsearch.googleapis.com").
					Get("/customsearch/v1").
					MatchParam("cx", "fake_search_engine_id").
					MatchParam("start", "0").
					Reply(429).
					JSON(&googleapi.Error{
						Code:    429,
						Message: "rate limit exceeded",
						Details: nil,
						Body:    "rate limit exceeded",
						Header:  http.Header{},
						Errors:  []googleapi.ErrorItem{},
					})
			},
		},
		{
			name:     "test with basic error",
			number:   test.NewFakeUSNumber(),
			expected: map[string]interface{}{},
			wantErrors: map[string]error{
				"googlecse": &googleapi.Error{
					Code:    403,
					Message: "",
					Details: nil,
					Body:    "{\"code\":403,\"message\":\"dummy error\",\"details\":null,\"Body\":\"dummy error\",\"Header\":{},\"Errors\":[]}\n",
					Header: http.Header{
						"Content-Type": []string{"application/json"},
					},
					Errors: nil,
				},
			},
			mocks: func() {
				gock.New("https://customsearch.googleapis.com").
					Get("/customsearch/v1").
					MatchParam("cx", "fake_search_engine_id").
					MatchParam("start", "0").
					Reply(403).
					JSON(&googleapi.Error{
						Code:    403,
						Message: "dummy error",
						Details: nil,
						Body:    "dummy error",
						Header:  http.Header{},
						Errors:  []googleapi.ErrorItem{},
					})
			},
		},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			_ = os.Setenv("GOOGLECSE_CX", "fake_search_engine_id")
			_ = os.Setenv("GOOGLE_API_KEY", "fake_api_key")
			defer os.Unsetenv("GOOGLECSE_CX")
			defer os.Unsetenv("GOOGLE_API_KEY")

			tt.mocks()
			defer gock.Off() // Flush pending mocks after test execution

			scanner := NewGoogleCSEScanner(&http.Client{})
			remote := NewLibrary(filter.NewEngine())
			remote.AddScanner(scanner)

			if scanner.DryRun(*tt.number, tt.opts) != nil {
				t.Fatal("DryRun() should return nil")
			}

			got, errs := remote.Scan(tt.number, tt.opts)
			if len(tt.wantErrors) > 0 {
				assert.Equal(t, tt.wantErrors, errs)
			} else {
				assert.Len(t, errs, 0)
			}
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestGoogleCSEScanner_DryRun(t *testing.T) {
	_ = os.Setenv("GOOGLECSE_CX", "abc")
	_ = os.Setenv("GOOGLE_API_KEY", "abc")
	defer os.Unsetenv("GOOGLECSE_CX")
	defer os.Unsetenv("GOOGLE_API_KEY")
	scanner := NewGoogleCSEScanner(&http.Client{})
	assert.Nil(t, scanner.DryRun(*test.NewFakeUSNumber(), ScannerOptions{}))
}

func TestGoogleCSEScanner_DryRunWithOptions(t *testing.T) {
	errStr := "search engine ID and/or API key is not defined"

	scanner := NewGoogleCSEScanner(&http.Client{})
	assert.Nil(t, scanner.DryRun(*test.NewFakeUSNumber(), ScannerOptions{"GOOGLECSE_CX": "test", "GOOGLE_API_KEY": "secret"}))
	assert.EqualError(t, scanner.DryRun(*test.NewFakeUSNumber(), ScannerOptions{"GOOGLECSE_CX": "", "GOOGLE_API_KEY": ""}), errStr)
	assert.EqualError(t, scanner.DryRun(*test.NewFakeUSNumber(), ScannerOptions{"GOOGLECSE_CX": "test"}), errStr)
	assert.EqualError(t, scanner.DryRun(*test.NewFakeUSNumber(), ScannerOptions{"GOOGLE_API_KEY": "test"}), errStr)
}

func TestGoogleCSEScanner_DryRun_Error(t *testing.T) {
	scanner := NewGoogleCSEScanner(&http.Client{})
	assert.EqualError(t, scanner.DryRun(*test.NewFakeUSNumber(), ScannerOptions{}), "search engine ID and/or API key is not defined")
}

func TestGoogleCSEScanner_MaxResults(t *testing.T) {
	testcases := []struct {
		value    string
		expected int64
	}{
		{
			value:    "",
			expected: 10,
		},
		{
			value:    "20",
			expected: 20,
		},
		{
			value:    "120",
			expected: 100,
		},
	}

	defer os.Clearenv()

	for _, tt := range testcases {
		_ = os.Setenv("GOOGLECSE_MAX_RESULTS", tt.value)

		scanner := NewGoogleCSEScanner(nil)

		assert.Equal(t, tt.expected, scanner.(*googleCSEScanner).MaxResults)
	}
}
