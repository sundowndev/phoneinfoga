package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/sundowndev/phoneinfoga/v2/lib/filter"
	"github.com/sundowndev/phoneinfoga/v2/lib/remote"
	"github.com/sundowndev/phoneinfoga/v2/mocks"
	"github.com/sundowndev/phoneinfoga/v2/test"
	"github.com/sundowndev/phoneinfoga/v2/web/v2/api"
	"github.com/sundowndev/phoneinfoga/v2/web/v2/api/handlers"
	"github.com/sundowndev/phoneinfoga/v2/web/v2/api/server"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetAllScanners(t *testing.T) {
	type expectedResponse struct {
		Code int
		Body interface{}
	}

	testcases := []struct {
		Name     string
		Expected expectedResponse
	}{
		{
			Name: "test getting all scanners",
			Expected: expectedResponse{
				Code: 200,
				Body: handlers.GetAllScannersResponse{
					Scanners: []handlers.Scanner{
						{
							Name:        "fakeScanner",
							Description: "fakeScanner description",
						},
					},
				},
			},
		},
	}

	for _, tt := range testcases {
		t.Run(tt.Name, func(t *testing.T) {
			fakeScanner := &mocks.Scanner{}
			fakeScanner.On("Name").Return("fakeScanner")
			fakeScanner.On("Description").Return("fakeScanner description")
			handlers.RemoteLibrary = remote.NewLibrary(filter.NewEngine())
			handlers.RemoteLibrary.AddScanner(fakeScanner)

			r := server.NewServer()

			req, err := http.NewRequest(http.MethodGet, "/v2/scanners", nil)
			if err != nil {
				t.Fatal(err)
			}
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			b, err := json.Marshal(tt.Expected.Body)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, tt.Expected.Code, w.Code)
			assert.Equal(t, string(b), w.Body.String())
			fakeScanner.AssertExpectations(t)
		})
	}
}

func TestDryRunScanner(t *testing.T) {
	type expectedResponse struct {
		Code int
		Body interface{}
	}

	type params struct {
		Supplier string
	}

	testcases := []struct {
		Name     string
		Params   params
		Body     interface{}
		Expected expectedResponse
		Mocks    func(*mocks.Scanner)
	}{
		{
			Name:   "test dry running scanner",
			Params: params{Supplier: "fakeScanner"},
			Body:   handlers.DryRunScannerInput{Number: "14152229670"},
			Expected: expectedResponse{
				Code: 200,
				Body: handlers.DryRunScannerResponse{Success: true},
			},
			Mocks: func(s *mocks.Scanner) {
				s.On("Name").Return("fakeScanner")
				s.On("DryRun", *test.NewFakeUSNumber(), remote.ScannerOptions{}).Return(nil)
			},
		},
		{
			Name:   "test dry running scanner with options",
			Params: params{Supplier: "fakeScanner"},
			Body: handlers.DryRunScannerInput{
				Number:  "14152229670",
				Options: remote.ScannerOptions{"api_key": "secret"},
			},
			Expected: expectedResponse{
				Code: 200,
				Body: handlers.DryRunScannerResponse{Success: true},
			},
			Mocks: func(s *mocks.Scanner) {
				s.On("Name").Return("fakeScanner")
				s.On("DryRun", *test.NewFakeUSNumber(), remote.ScannerOptions{"api_key": "secret"}).Return(nil)
			},
		},
		{
			Name:   "test dry running scanner with empty options",
			Params: params{Supplier: "fakeScanner"},
			Body: handlers.DryRunScannerInput{
				Number:  "14152229670",
				Options: remote.ScannerOptions{},
			},
			Expected: expectedResponse{
				Code: 200,
				Body: handlers.DryRunScannerResponse{Success: true},
			},
			Mocks: func(s *mocks.Scanner) {
				s.On("Name").Return("fakeScanner")
				s.On("DryRun", *test.NewFakeUSNumber(), remote.ScannerOptions{}).Return(nil)
			},
		},
		{
			Name:   "test dry running scanner with error",
			Params: params{Supplier: "fakeScanner"},
			Body:   handlers.DryRunScannerInput{Number: "14152229670"},
			Expected: expectedResponse{
				Code: 400,
				Body: handlers.DryRunScannerResponse{Success: false, Error: "dummy error"},
			},
			Mocks: func(s *mocks.Scanner) {
				s.On("Name").Return("fakeScanner")
				s.On("DryRun", *test.NewFakeUSNumber(), make(remote.ScannerOptions)).Return(errors.New("dummy error"))
			},
		},
		{
			Name:   "test invalid number",
			Params: params{Supplier: "fakeScanner"},
			Body:   handlers.DryRunScannerInput{Number: "1.4152229670"},
			Expected: expectedResponse{
				Code: 400,
				Body: api.ErrorResponse{Error: "Invalid phone number: please provide an integer without any special chars"},
			},
			Mocks: func(s *mocks.Scanner) {
				s.On("Name").Return("fakeScanner")
			},
		},
		{
			Name:   "test scanner not found",
			Params: params{Supplier: "test"},
			Body:   handlers.DryRunScannerInput{Number: "14152229670"},
			Expected: expectedResponse{
				Code: 404,
				Body: api.ErrorResponse{Error: "Scanner not found"},
			},
			Mocks: func(s *mocks.Scanner) {
				s.On("Name").Return("fakeScanner")
			},
		},
		{
			Name:   "test invalid number",
			Params: params{Supplier: "fakeScanner"},
			Body:   handlers.DryRunScannerInput{Number: "222"},
			Expected: expectedResponse{
				Code: 400,
				Body: api.ErrorResponse{Error: "the string supplied is too short to be a phone number"},
			},
			Mocks: func(s *mocks.Scanner) {
				s.On("Name").Return("fakeScanner")
			},
		},
	}

	for _, tt := range testcases {
		t.Run(tt.Name, func(t *testing.T) {
			fakeScanner := &mocks.Scanner{}
			tt.Mocks(fakeScanner)
			handlers.RemoteLibrary = remote.NewLibrary(filter.NewEngine())
			handlers.RemoteLibrary.AddScanner(fakeScanner)

			data, err := json.Marshal(&tt.Body)
			if err != nil {
				t.Fatal(err)
			}

			req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("/v2/scanners/%s/dryrun", tt.Params.Supplier), bytes.NewReader(data))
			if err != nil {
				t.Fatal(err)
			}
			w := httptest.NewRecorder()
			server.NewServer().ServeHTTP(w, req)

			b, err := json.Marshal(tt.Expected.Body)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, tt.Expected.Code, w.Code)
			assert.Equal(t, string(b), w.Body.String())
			fakeScanner.AssertExpectations(t)
		})
	}
}

func TestRunScanner(t *testing.T) {
	type FakeScannerResponse struct {
		Info string `json:"info"`
	}

	type expectedResponse struct {
		Code int
		Body interface{}
	}

	type params struct {
		Supplier string
	}

	testcases := []struct {
		Name     string
		Params   params
		Body     interface{}
		Expected expectedResponse
		Mocks    func(*mocks.Scanner)
	}{
		{
			Name:   "test running scanner",
			Params: params{Supplier: "fakeScanner"},
			Body:   handlers.RunScannerInput{Number: "14152229670"},
			Expected: expectedResponse{
				Code: 200,
				Body: handlers.RunScannerResponse{
					Result: FakeScannerResponse{Info: "test"},
				},
			},
			Mocks: func(s *mocks.Scanner) {
				s.On("Name").Return("fakeScanner")
				s.On("Run", *test.NewFakeUSNumber(), remote.ScannerOptions{}).Return(FakeScannerResponse{Info: "test"}, nil)
			},
		},
		{
			Name:   "test running scanner with error",
			Params: params{Supplier: "fakeScanner"},
			Body:   handlers.RunScannerInput{Number: "14152229670"},
			Expected: expectedResponse{
				Code: 500,
				Body: api.ErrorResponse{Error: "dummy error"},
			},
			Mocks: func(s *mocks.Scanner) {
				s.On("Name").Return("fakeScanner")
				s.On("Run", *test.NewFakeUSNumber(), remote.ScannerOptions{}).Return(nil, errors.New("dummy error"))
			},
		},
		{
			Name:   "test invalid number",
			Params: params{Supplier: "fakeScanner"},
			Body:   handlers.RunScannerInput{Number: "1.4152229670"},
			Expected: expectedResponse{
				Code: 400,
				Body: api.ErrorResponse{Error: "Invalid phone number: please provide an integer without any special chars"},
			},
			Mocks: func(s *mocks.Scanner) {
				s.On("Name").Return("fakeScanner")
			},
		},
		{
			Name:   "test scanner not found",
			Params: params{Supplier: "test"},
			Body:   handlers.RunScannerInput{Number: "14152229670"},
			Expected: expectedResponse{
				Code: 404,
				Body: api.ErrorResponse{Error: "Scanner not found"},
			},
			Mocks: func(s *mocks.Scanner) {
				s.On("Name").Return("fakeScanner")
			},
		},
		{
			Name:   "test invalid number",
			Params: params{Supplier: "fakeScanner"},
			Body:   handlers.RunScannerInput{Number: "222"},
			Expected: expectedResponse{
				Code: 400,
				Body: api.ErrorResponse{Error: "the string supplied is too short to be a phone number"},
			},
			Mocks: func(s *mocks.Scanner) {
				s.On("Name").Return("fakeScanner")
			},
		},
	}

	for _, tt := range testcases {
		t.Run(tt.Name, func(t *testing.T) {
			fakeScanner := &mocks.Scanner{}
			tt.Mocks(fakeScanner)
			handlers.RemoteLibrary = remote.NewLibrary(filter.NewEngine())
			handlers.RemoteLibrary.AddScanner(fakeScanner)

			data, err := json.Marshal(&tt.Body)
			if err != nil {
				t.Fatal(err)
			}

			req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("/v2/scanners/%s/run", tt.Params.Supplier), bytes.NewReader(data))
			if err != nil {
				t.Fatal(err)
			}
			w := httptest.NewRecorder()
			server.NewServer().ServeHTTP(w, req)

			b, err := json.Marshal(tt.Expected.Body)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, tt.Expected.Code, w.Code)
			assert.Equal(t, string(b), w.Body.String())
			fakeScanner.AssertExpectations(t)
		})
	}
}
