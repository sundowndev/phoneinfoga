package handlers_test

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/sundowndev/phoneinfoga/v2/web/v2/api"
	"github.com/sundowndev/phoneinfoga/v2/web/v2/api/handlers"
	"github.com/sundowndev/phoneinfoga/v2/web/v2/api/server"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAddNumber(t *testing.T) {
	type expectedResponse struct {
		Code int
		Body interface{}
	}

	testcases := []struct {
		Name     string
		Input    handlers.AddNumberInput
		Expected expectedResponse
	}{
		{
			Name:  "test successfully adding number",
			Input: handlers.AddNumberInput{Number: "14152229670"},
			Expected: expectedResponse{
				Code: 200,
				Body: handlers.AddNumberResponse{
					Valid:         true,
					RawLocal:      "4152229670",
					Local:         "(415) 222-9670",
					E164:          "+14152229670",
					International: "14152229670",
					CountryCode:   1,
					Country:       "US",
					Carrier:       "",
				},
			},
		},
		{
			Name:  "test bad params",
			Input: handlers.AddNumberInput{Number: "a14152229670"},
			Expected: expectedResponse{
				Code: 400,
				Body: api.ErrorResponse{Error: "Invalid phone number: please provide an integer without any special chars"},
			},
		},
		{
			Name:  "test invalid number",
			Input: handlers.AddNumberInput{Number: "331"},
			Expected: expectedResponse{
				Code: 400,
				Body: api.ErrorResponse{Error: "the string supplied is too short to be a phone number"},
			},
		},
	}

	for _, tt := range testcases {
		t.Run(tt.Name, func(t *testing.T) {
			r := server.NewServer()

			data, err := json.Marshal(&tt.Input)
			if err != nil {
				t.Fatal(err)
			}

			req, err := http.NewRequest(http.MethodPost, "/v2/numbers", bytes.NewReader(data))
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
		})
	}
}
