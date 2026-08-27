package remote_test

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/sundowndev/phoneinfoga/v2/lib/filter"
	"github.com/sundowndev/phoneinfoga/v2/lib/number"
	"github.com/sundowndev/phoneinfoga/v2/lib/remote"
	"github.com/sundowndev/phoneinfoga/v2/lib/remote/suppliers"
	"testing"
)

type ftcSupplierMock struct {
	response *suppliers.FTCComplaintsResponse
	err      error
}

func (m *ftcSupplierMock) CountComplaints(_ string, _ int) (*suppliers.FTCComplaintsResponse, error) {
	return m.response, m.err
}

func TestFTCScanner_Metadata(t *testing.T) {
	scanner := remote.NewFTCScanner(&ftcSupplierMock{})
	assert.Equal(t, remote.FTC, scanner.Name())
	assert.NotEmpty(t, scanner.Description())
}

func TestFTCScanner(t *testing.T) {
	dummyError := errors.New("dummy")
	usNumber, _ := number.NewNumber("17027515054")
	frNumber, _ := number.NewNumber("33365174444")

	testcases := []struct {
		name       string
		number     *number.Number
		supplier   *ftcSupplierMock
		expected   map[string]interface{}
		wantErrors map[string]error
	}{
		{
			name:   "no complaints found",
			number: usNumber,
			supplier: &ftcSupplierMock{
				response: &suppliers.FTCComplaintsResponse{DaysChecked: 5},
			},
			expected: map[string]interface{}{
				"ftc": remote.FTCComplaintScannerResponse{
					DaysChecked:    5,
					ComplaintCount: 0,
				},
			},
			wantErrors: map[string]error{},
		},
		{
			name:   "complaints found",
			number: usNumber,
			supplier: &ftcSupplierMock{
				response: &suppliers.FTCComplaintsResponse{
					DaysChecked: 5,
					Complaints: []suppliers.FTCComplaint{
						{Date: "2026-07-29", State: "Nevada", Subject: "Other", Robocall: "Y"},
					},
				},
			},
			expected: map[string]interface{}{
				"ftc": remote.FTCComplaintScannerResponse{
					DaysChecked:    5,
					ComplaintCount: 1,
					RecentComplaints: []remote.FTCComplaintEntry{
						{Date: "2026-07-29", State: "Nevada", Subject: "Other", Robocall: "Y"},
					},
				},
			},
			wantErrors: map[string]error{},
		},
		{
			name:       "failed scan",
			number:     usNumber,
			supplier:   &ftcSupplierMock{err: dummyError},
			expected:   map[string]interface{}{},
			wantErrors: map[string]error{"ftc": dummyError},
		},
		{
			name:       "country not supported",
			number:     frNumber,
			supplier:   &ftcSupplierMock{},
			expected:   map[string]interface{}{},
			wantErrors: map[string]error{},
		},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			scanner := remote.NewFTCScanner(tt.supplier)
			lib := remote.NewLibrary(filter.NewEngine())
			lib.AddScanner(scanner)

			got, errs := lib.Scan(tt.number, remote.ScannerOptions{})
			if len(tt.wantErrors) > 0 {
				assert.Equal(t, tt.wantErrors, errs)
			} else {
				assert.Len(t, errs, 0)
			}
			assert.Equal(t, tt.expected, got)
		})
	}
}
