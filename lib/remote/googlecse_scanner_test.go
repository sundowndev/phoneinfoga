package remote

import (
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

func TestGoogleCSEScanner(t *testing.T) {
	testcases := []struct {
		name       string
		number     *number.Number
		expected   map[string]interface{}
		wantErrors map[string]error
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
		},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			_ = os.Setenv("GOOGLECSE_CX", "fake_search_engine_id")
			_ = os.Setenv("GOOGLE_API_KEY", "fake_api_key")
			defer os.Clearenv()

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
			defer gock.Off() // Flush pending mocks after test execution

			scanner := NewGoogleCSEScanner(&http.Client{})
			remote := NewLibrary(filter.NewEngine())
			remote.AddScanner(scanner)

			if !scanner.ShouldRun(*tt.number) {
				t.Fatal("ShouldRun() should be truthy")
			}

			got, errs := remote.Scan(tt.number)
			if len(tt.wantErrors) > 0 {
				assert.Equal(t, tt.wantErrors, errs)
			} else {
				assert.Len(t, errs, 0)
			}
			assert.Equal(t, tt.expected, got)
		})
	}
}
