package suppliers

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"time"
)

type FTCSupplierInterface interface {
	CountComplaints(rawLocal string, days int) (*FTCComplaintsResponse, error)
}

type FTCComplaint struct {
	Date     string
	State    string
	Subject  string
	Robocall string
}

type FTCComplaintsResponse struct {
	DaysChecked int
	Complaints  []FTCComplaint
}

type FTCSupplier struct {
	BaseUri string
}

func NewFTCSupplier() *FTCSupplier {
	return &FTCSupplier{BaseUri: "https://www.ftc.gov/sites/default/files"}
}

// CountComplaints scrapes the FTC's daily Do Not Call complaint CSVs, published at
// a predictable URL per weekday (no auth, no API key - the documented api.ftc.gov
// endpoint returns static sandbox data regardless of parameters, so this is the
// only way to actually filter by phone number). There's no per-number server-side
// filter, so each day's file (roughly 1MB, a few thousand rows) is matched client-side.
func (s *FTCSupplier) CountComplaints(rawLocal string, days int) (*FTCComplaintsResponse, error) {
	res := &FTCComplaintsResponse{}

	for i := 0; i < days; i++ {
		date := time.Now().AddDate(0, 0, -i)
		url := fmt.Sprintf("%s/DNC_Complaint_Numbers_%s.csv", s.BaseUri, date.Format("2006-01-02"))

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, err
		}

		// No file is published on weekends/holidays - just skip those days.
		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			continue
		}
		res.DaysChecked++

		rows, err := csv.NewReader(resp.Body).ReadAll()
		resp.Body.Close()
		if err != nil || len(rows) < 2 {
			continue
		}

		col := columnIndexer(rows[0])
		phoneIdx := col("Company_Phone_Number")
		if phoneIdx == -1 {
			continue
		}
		dateIdx, stateIdx, subjectIdx, robocallIdx := col("Created_Date"), col("Consumer_State"), col("Subject"), col("Recorded_Message_Or_Robocall")

		for _, row := range rows[1:] {
			if phoneIdx >= len(row) || !matchesLocalNumber(row[phoneIdx], rawLocal) {
				continue
			}
			res.Complaints = append(res.Complaints, FTCComplaint{
				Date:     safeCol(row, dateIdx),
				State:    safeCol(row, stateIdx),
				Subject:  safeCol(row, subjectIdx),
				Robocall: safeCol(row, robocallIdx),
			})
		}
	}

	return res, nil
}

func columnIndexer(header []string) func(name string) int {
	return func(name string) int {
		for i, h := range header {
			if h == name {
				return i
			}
		}
		return -1
	}
}

func safeCol(row []string, idx int) string {
	if idx == -1 || idx >= len(row) {
		return ""
	}
	return row[idx]
}

// matchesLocalNumber compares a CSV phone value against a 10-digit US/Canada
// local number, stripping a leading NANP country code "1" if present (an area
// code can never start with 0 or 1, so this can't misfire on a real local number).
func matchesLocalNumber(csvPhone, rawLocal string) bool {
	if len(csvPhone) == 11 && csvPhone[0] == '1' {
		csvPhone = csvPhone[1:]
	}
	return csvPhone == rawLocal
}
