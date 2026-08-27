package suppliers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

type VeriphoneSupplierInterface interface {
	Verify(apiKey string, internationalNumber string) (*VeriphoneResponse, error)
}

// VeriphoneResponse mirrors https://api.veriphone.io/v3/verify's JSON shape
type VeriphoneResponse struct {
	Status              string `json:"status"`
	Phone               string `json:"phone"`
	PhoneValid          bool   `json:"phone_valid"`
	PhoneType           string `json:"phone_type"`
	PhoneRegion         string `json:"phone_region"`
	Country             string `json:"country"`
	CountryCode         string `json:"country_code"`
	CountryPrefix       string `json:"country_prefix"`
	InternationalNumber string `json:"international_number"`
	LocalNumber         string `json:"local_number"`
	E164                string `json:"e164"`
	Carrier             string `json:"carrier"`
	Error               string `json:"error"`
}

type VeriphoneSupplier struct {
	Uri string
}

func NewVeriphoneSupplier() *VeriphoneSupplier {
	return &VeriphoneSupplier{Uri: "https://api.veriphone.io/v3"}
}

func (s *VeriphoneSupplier) Verify(apiKey string, internationalNumber string) (*VeriphoneResponse, error) {
	reqUrl := fmt.Sprintf("%s/verify?phone=%s", s.Uri, url.QueryEscape("+"+internationalNumber))

	req, err := http.NewRequest("GET", reqUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result VeriphoneResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if result.Status != "success" {
		if result.Error != "" {
			return nil, errors.New(result.Error)
		}
		return nil, errors.New("veriphone request failed")
	}

	return &result, nil
}
