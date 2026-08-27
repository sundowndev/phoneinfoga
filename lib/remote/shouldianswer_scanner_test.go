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

type shouldianswerSupplierMock struct {
	response *suppliers.ShouldIAnswerScannerResponse
	err      error
}

func (m *shouldianswerSupplierMock) Lookup(_ string) (*suppliers.ShouldIAnswerScannerResponse, error) {
	return m.response, m.err
}

func TestShouldIAnswerScanner_Metadata(t *testing.T) {
	scanner := remote.NewShouldIAnswerScanner(&shouldianswerSupplierMock{})
	assert.Equal(t, remote.ShouldIAnswer, scanner.Name())
	assert.NotEmpty(t, scanner.Description())
}

func TestShouldIAnswerScanner(t *testing.T) {
	dummyError := errors.New("dummy")
	usNumber, _ := number.NewNumber("17027515054")
	frNumber, _ := number.NewNumber("33365174444")

	testcases := []struct {
		name       string
		number     *number.Number
		supplier   *shouldianswerSupplierMock
		expected   map[string]interface{}
		wantErrors map[string]error
	}{
		{
			name:   "successful scan",
			number: usNumber,
			supplier: &shouldianswerSupplierMock{
				response: &suppliers.ShouldIAnswerScannerResponse{Score: "unknown", Summary: "no reports"},
			},
			expected: map[string]interface{}{
				"shouldianswer": remote.ShouldIAnswerScannerResponse{Score: "unknown", Summary: "no reports"},
			},
			wantErrors: map[string]error{},
		},
		{
			name:   "failed scan",
			number: usNumber,
			supplier: &shouldianswerSupplierMock{
				err: dummyError,
			},
			expected: map[string]interface{}{},
			wantErrors: map[string]error{
				"shouldianswer": dummyError,
			},
		},
		{
			name:       "country not supported",
			number:     frNumber,
			supplier:   &shouldianswerSupplierMock{},
			expected:   map[string]interface{}{},
			wantErrors: map[string]error{},
		},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			scanner := remote.NewShouldIAnswerScanner(tt.supplier)
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
