package output

import (
	"bytes"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/sundowndev/phoneinfoga/v2/lib/remote"
	"github.com/sundowndev/phoneinfoga/v2/test/goldenfile"
	"os"
	"testing"
)

func TestConsoleOutput(t *testing.T) {
	type FakeScannerResponse struct {
		Format string `console:"Number Format"`
	}

	type FakeScannerResponseRecursive2 struct {
		Response struct {
			Format FakeScannerResponse `console:"Format"`
		} `console:"Response"`
	}

	testcases := []struct {
		name    string
		dirName string
		result  map[string]interface{}
		errs    map[string]error
		wantErr error
	}{
		{
			name:    "should produce empty output",
			dirName: "testdata/console_empty.txt",
			result:  map[string]interface{}{},
			errs:    map[string]error{},
		},
		{
			name:    "should produce valid output",
			dirName: "testdata/console_valid.txt",
			result: map[string]interface{}{
				"numverify": remote.NumverifyScannerResponse{
					Valid:               true,
					Number:              "test",
					LocalFormat:         "test",
					InternationalFormat: "test",
					CountryPrefix:       "test",
					CountryCode:         "test",
					CountryName:         "test",
					Location:            "test",
					Carrier:             "test",
					LineType:            "test",
				},
			},
			errs: map[string]error{},
		},
		{
			name:    "should produce valid output with errors",
			dirName: "testdata/console_valid_with_errors.txt",
			result: map[string]interface{}{
				"testscanner": nil,
				"numverify": remote.NumverifyScannerResponse{
					Valid:               true,
					Number:              "test",
					LocalFormat:         "test",
					InternationalFormat: "test",
					CountryPrefix:       "test",
					CountryCode:         "test",
					CountryName:         "test",
					Location:            "test",
					Carrier:             "test",
					LineType:            "test",
				},
				"testscanner2": FakeScannerResponse{
					Format: "test",
				},
			},
			errs: map[string]error{
				"googlesearch": errors.New("dummy error"),
				"fakescanner":  errors.New("dummy error 2"),
			},
		},
		{
			name:    "should follow recursive paths",
			dirName: "testdata/console_valid_recursive.txt",
			result: map[string]interface{}{
				"testscanner": remote.GoogleSearchResponse{
					SocialMedia: []*remote.GoogleSearchDork{
						{
							URL:    "http://example.com?q=111-555-1212",
							Number: "111-555-1212",
							Dork:   "intext:\"111-555-1212\"",
						},
						{
							URL:    "http://example.com?q=222-666-2323",
							Number: "222-666-2323",
							Dork:   "intext:\"222-666-2323\"",
						},
					},
				},
				"testscanner2": FakeScannerResponseRecursive2{
					Response: struct {
						Format FakeScannerResponse `console:"Format"`
					}{Format: FakeScannerResponse{Format: "test"}},
				},
			},
			errs: map[string]error{},
		},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			shouldUpdate := tt.dirName == *goldenfile.Update

			expected, err := os.ReadFile(tt.dirName)
			if err != nil && !shouldUpdate {
				t.Fatal(err)
			}

			got := new(bytes.Buffer)
			err = GetOutput(Console, got).Write(tt.result, tt.errs)
			if tt.wantErr != nil {
				assert.EqualError(t, err, tt.wantErr.Error())
			} else {
				assert.Nil(t, err)
			}

			if shouldUpdate {
				err = os.WriteFile(tt.dirName, got.Bytes(), 0644)
				if err != nil {
					t.Fatal(err)
				}
				expected, err = os.ReadFile(tt.dirName)
				if err != nil {
					t.Fatal(err)
				}
			}

			assert.Equal(t, string(expected), got.String())
		})
	}
}
