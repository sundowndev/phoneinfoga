package remote

import (
	"errors"
	"github.com/sundowndev/phoneinfoga/v2/lib/number"
	"github.com/sundowndev/phoneinfoga/v2/lib/remote/suppliers"
)

const ShouldIAnswer = "shouldianswer"

type shouldianswerScanner struct {
	client suppliers.ShouldIAnswerSupplierInterface
}

type ShouldIAnswerScannerResponse struct {
	Score   string `json:"score" console:"Score"`
	Summary string `json:"summary,omitempty" console:"Summary,omitempty"`
}

func NewShouldIAnswerScanner(s suppliers.ShouldIAnswerSupplierInterface) Scanner {
	return &shouldianswerScanner{client: s}
}

func (s *shouldianswerScanner) Name() string {
	return ShouldIAnswer
}

func (s *shouldianswerScanner) Description() string {
	return "Look up a US/Canada number's community rating on shouldianswer.com (free, no API key)."
}

func (s *shouldianswerScanner) DryRun(n number.Number, _ ScannerOptions) error {
	if n.CountryCode != 1 || len(n.RawLocal) != 10 {
		return errors.New("only US/Canada numbers are supported")
	}
	return nil
}

func (s *shouldianswerScanner) Run(n number.Number, _ ScannerOptions) (interface{}, error) {
	res, err := s.client.Lookup(n.RawLocal)
	if err != nil {
		return nil, err
	}

	return ShouldIAnswerScannerResponse{Score: res.Score, Summary: res.Summary}, nil
}
