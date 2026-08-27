package suppliers

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
)

type ShouldIAnswerSupplierInterface interface {
	Lookup(rawLocal string) (*ShouldIAnswerScannerResponse, error)
}

// ShouldIAnswerScannerResponse is the Should I Answer scanner response
type ShouldIAnswerScannerResponse struct {
	Score   string // "unknown", "positive" or "negative"
	Summary string
}

type ShouldIAnswerSupplier struct {
	Uri string
}

func NewShouldIAnswerSupplier() *ShouldIAnswerSupplier {
	return &ShouldIAnswerSupplier{Uri: "https://www.shouldianswer.com"}
}

// The page carries the community verdict as a CSS class (score unknown/positive/negative)
// and a plain-text summary paragraph - no login wall, no API needed.
var shouldianswerScoreRegexp = regexp.MustCompile(`class="score (\w+)"`)
var shouldianswerInfoxRegexp = regexp.MustCompile(`(?s)class="infox">\s*(.+?)\s*</div>`)
var htmlTagRegexp = regexp.MustCompile(`<[^>]+>`)

func (s *ShouldIAnswerSupplier) Lookup(rawLocal string) (*ShouldIAnswerScannerResponse, error) {
	if len(rawLocal) != 10 {
		return nil, fmt.Errorf("expected a 10-digit US/Canada number, got %q", rawLocal)
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/phone-number/%s", s.Uri, rawLocal), nil)
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

	res := &ShouldIAnswerScannerResponse{Score: "unknown"}

	if m := shouldianswerScoreRegexp.FindSubmatch(body); m != nil {
		res.Score = string(m[1])
	}
	if m := shouldianswerInfoxRegexp.FindSubmatch(body); m != nil {
		res.Summary = strings.TrimSpace(htmlTagRegexp.ReplaceAllString(string(m[1]), ""))
	}

	return res, nil
}
