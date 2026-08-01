package suppliers

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
)

type NomoroboSupplierInterface interface {
	Lookup(rawLocal string) (*NomoroboScannerResponse, error)
}

// NomoroboScannerResponse is the Nomorobo scanner response
type NomoroboScannerResponse struct {
	Description string
}

type NomoroboSupplier struct {
	Uri string
}

func NewNomoroboSupplier() *NomoroboSupplier {
	return &NomoroboSupplier{Uri: "https://www.nomorobo.com"}
}

// Nomorobo puts the actual verdict in the page's meta description, e.g.
// "(702) 751-5054 is an Unknown Caller. We have no other information about this number."
var nomoroboDescriptionRegexp = regexp.MustCompile(`name="description" content="([^"]+)"`)

func (s *NomoroboSupplier) Lookup(rawLocal string) (*NomoroboScannerResponse, error) {
	if len(rawLocal) != 10 {
		return nil, fmt.Errorf("expected a 10-digit US/Canada number, got %q", rawLocal)
	}
	dashed := fmt.Sprintf("%s-%s-%s", rawLocal[0:3], rawLocal[3:6], rawLocal[6:10])

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/lookup/%s", s.Uri, dashed), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	match := nomoroboDescriptionRegexp.FindSubmatch(body)
	if match == nil {
		return &NomoroboScannerResponse{}, nil
	}

	return &NomoroboScannerResponse{Description: string(match[1])}, nil
}
