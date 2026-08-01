package remote

import (
	"errors"
	"github.com/sundowndev/phoneinfoga/v2/lib/number"
	"github.com/sundowndev/phoneinfoga/v2/lib/remote/suppliers"
)

const FTC = "ftc"

// ftcDaysToCheck bounds how many calendar days back to scan - each day is a
// ~1MB CSV fetched in full (no server-side phone filter exists), so this is a
// deliberate cost/recency tradeoff, not an arbitrary number.
const ftcDaysToCheck = 7

type ftcScanner struct {
	client suppliers.FTCSupplierInterface
}

type FTCComplaintEntry struct {
	Date     string `json:"date" console:"Date,omitempty"`
	State    string `json:"state" console:"State,omitempty"`
	Subject  string `json:"subject" console:"Subject,omitempty"`
	Robocall string `json:"robocall" console:"Robocall,omitempty"`
}

type FTCComplaintScannerResponse struct {
	DaysChecked      int                 `json:"days_checked" console:"Days checked"`
	ComplaintCount   int                 `json:"complaint_count" console:"Complaints reported"`
	RecentComplaints []FTCComplaintEntry `json:"recent_complaints,omitempty" console:"Recent complaints,omitempty"`
}

func NewFTCScanner(s suppliers.FTCSupplierInterface) Scanner {
	return &ftcScanner{client: s}
}

func (s *ftcScanner) Name() string {
	return FTC
}

func (s *ftcScanner) Description() string {
	return "Check the FTC's public Do Not Call complaint data (last published week) for reports against this number."
}

func (s *ftcScanner) DryRun(n number.Number, _ ScannerOptions) error {
	if n.CountryCode != 1 || len(n.RawLocal) != 10 {
		return errors.New("only US/Canada numbers are supported")
	}
	return nil
}

func (s *ftcScanner) Run(n number.Number, _ ScannerOptions) (interface{}, error) {
	res, err := s.client.CountComplaints(n.RawLocal, ftcDaysToCheck)
	if err != nil {
		return nil, err
	}

	// leave nil (not an empty slice) when there are no complaints, so the
	// "omitempty" console tag actually suppresses the field
	var entries []FTCComplaintEntry
	for _, c := range res.Complaints {
		entries = append(entries, FTCComplaintEntry{
			Date:     c.Date,
			State:    c.State,
			Subject:  c.Subject,
			Robocall: c.Robocall,
		})
	}

	return FTCComplaintScannerResponse{
		DaysChecked:      res.DaysChecked,
		ComplaintCount:   len(entries),
		RecentComplaints: entries,
	}, nil
}
