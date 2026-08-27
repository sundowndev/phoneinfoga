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

type nomoroboSupplierMock struct {
	response *suppliers.NomoroboScannerResponse
	err      error
}

func (m *nomoroboSupplierMock) Lookup(_ string) (*suppliers.NomoroboScannerResponse, error) {
	return m.response, m.err
}

func TestNomoroboScanner_Metadata(t *testing.T) {
	scanner := remote.NewNomoroboScanner(&nomoroboSupplierMock{})
	assert.Equal(t, remote.Nomorobo, scanner.Name())
	assert.NotEmpty(t, scanner.Description())
}

func TestNomoroboScanner(t *testing.T) {
	dummyError := errors.New("dummy")
	usNumber, _ := number.NewNumber("17027515054")
	frNumber, _ := number.NewNumber("33365174444")

	testcases := []struct {
		name       string
		number     *number.Number
		supplier   *nomoroboSupplierMock
		expected   map[string]interface{}
		wantErrors map[string]error
	}{
		{
			name:   "successful scan",
			number: usNumber,
			supplier: &nomoroboSupplierMock{
				response: &suppliers.NomoroboScannerResponse{Description: "Unknown Caller"},
			},
			expected: map[string]interface{}{
				"nomorobo": remote.NomoroboScannerResponse{Description: "Unknown Caller"},
			},
			wantErrors: map[string]error{},
		},
		{
			name:   "failed scan",
			number: usNumber,
			supplier: &nomoroboSupplierMock{
				err: dummyError,
			},
			expected: map[string]interface{}{},
			wantErrors: map[string]error{
				"nomorobo": dummyError,
			},
		},
		{
			name:       "country not supported",
			number:     frNumber,
			supplier:   &nomoroboSupplierMock{},
			expected:   map[string]interface{}{},
			wantErrors: map[string]error{},
		},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			scanner := remote.NewNomoroboScanner(tt.supplier)
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
