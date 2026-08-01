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

type veriphoneSupplierMock struct {
	response *suppliers.VeriphoneResponse
	err      error
}

func (m *veriphoneSupplierMock) Verify(_ string, _ string) (*suppliers.VeriphoneResponse, error) {
	return m.response, m.err
}

func TestVeriphoneScanner_Metadata(t *testing.T) {
	scanner := remote.NewVeriphoneScanner(&veriphoneSupplierMock{})
	assert.Equal(t, remote.Veriphone, scanner.Name())
	assert.NotEmpty(t, scanner.Description())
}

func TestVeriphoneScanner(t *testing.T) {
	dummyError := errors.New("dummy")
	dummyNumber, _ := number.NewNumber("17027515054")

	testcases := []struct {
		name       string
		opts       remote.ScannerOptions
		supplier   *veriphoneSupplierMock
		expected   map[string]interface{}
		wantErrors map[string]error
	}{
		{
			name: "successful scan",
			opts: remote.ScannerOptions{"VERIPHONE_API_KEY": "secret"},
			supplier: &veriphoneSupplierMock{
				response: &suppliers.VeriphoneResponse{
					PhoneValid:  true,
					PhoneType:   "voip",
					Country:     "United States",
					CountryCode: "US",
					Carrier:     "Bandwidth.com",
				},
			},
			expected: map[string]interface{}{
				"veriphone": remote.VeriphoneScannerResponse{
					PhoneValid:  true,
					PhoneType:   "voip",
					Country:     "United States",
					CountryCode: "US",
					Carrier:     "Bandwidth.com",
					SpoofRisk:   "elevated - line type reported as VOIP",
				},
			},
			wantErrors: map[string]error{},
		},
		{
			name: "failed scan",
			opts: remote.ScannerOptions{"VERIPHONE_API_KEY": "secret"},
			supplier: &veriphoneSupplierMock{
				err: dummyError,
			},
			expected: map[string]interface{}{},
			wantErrors: map[string]error{
				"veriphone": dummyError,
			},
		},
		{
			name:       "should not run without an API key",
			opts:       remote.ScannerOptions{},
			supplier:   &veriphoneSupplierMock{},
			expected:   map[string]interface{}{},
			wantErrors: map[string]error{},
		},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			scanner := remote.NewVeriphoneScanner(tt.supplier)
			lib := remote.NewLibrary(filter.NewEngine())
			lib.AddScanner(scanner)

			got, errs := lib.Scan(dummyNumber, tt.opts)
			if len(tt.wantErrors) > 0 {
				assert.Equal(t, tt.wantErrors, errs)
			} else {
				assert.Len(t, errs, 0)
			}
			assert.Equal(t, tt.expected, got)
		})
	}
}
