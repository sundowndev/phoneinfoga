package remote

import (
	"encoding/json"
	"fmt"
	"github.com/sundowndev/phoneinfoga/v2/lib/number"
	"net/http"
)

const Tempophone = "tempophone"

type TempophoneScanner struct{}

type TempophoneScannerResponse struct {
	Found    bool   `json:"found" console:"Found"`
	Id       int    `json:"Id" console:"Id,omitempty"`
	EventURL string `json:"event_url" console:"Event URL,omitempty"`
}

type tempophonePhone struct {
	Id    int    `json:"id"`
	Phone string `json:"phone"`
}

type tempophonePhonesResponse struct {
	Meta struct {
		Limit      int `json:"limit"`
		Offset     int `json:"offset"`
		TotalCount int `json:"total_count"`
	}
	Objects []tempophonePhone `json:"objects"`
}

func NewTempophoneScanner() *TempophoneScanner {
	return &TempophoneScanner{}
}

func (s *TempophoneScanner) Name() string {
	return Tempophone
}

func (s *TempophoneScanner) Description() string {
	return "Check if number is owned by tempophone.com"
}

func (s *TempophoneScanner) DryRun(n number.Number) error {
	if n.CountryCode != 1 {
		return fmt.Errorf("country code %d is not supported", n.CountryCode)
	}
	return nil
}

func (s *TempophoneScanner) Run(n number.Number) (interface{}, error) {
	var phones []tempophonePhone
	var body tempophonePhonesResponse

	offset := -1
	for {
		if offset != -1 && len(phones) == body.Meta.TotalCount {
			break
		}

		req, err := http.NewRequest(http.MethodGet, "https://tempophone.com/api/v1/phones", nil)
		if err != nil {
			return nil, err
		}

		req.URL.Query().Add("offset", fmt.Sprintf("%d", offset))

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}

		err = json.NewDecoder(resp.Body).Decode(&body)
		if err != nil {
			return nil, err
		}

		offset = body.Meta.Offset
		phones = append(phones, body.Objects...)

		err = resp.Body.Close()
		if err != nil {
			return nil, err
		}
	}

	data := TempophoneScannerResponse{}
	for _, phone := range phones {
		if phone.Phone == n.International {
			data.Found = true
			data.Id = phone.Id
			data.EventURL = fmt.Sprintf("https://tempophone.com/phone/%d/events", phone.Id)
			break
		}
	}

	return data, nil
}
