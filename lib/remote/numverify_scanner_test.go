package remote

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/sundowndev/phoneinfoga/v2/lib/number"
	"github.com/sundowndev/phoneinfoga/v2/lib/remote/suppliers"
	"github.com/sundowndev/phoneinfoga/v2/mocks"
	"testing"
)

func TestNumverifyScanner(t *testing.T) {
	dummyError := errors.New("dummy")

	testcases := []struct {
		name       string
		number     *number.Number
		mocks      func(s *mocks.NumverifySupplier)
		expected   map[string]interface{}
		wantErrors map[string]error
	}{
		{
			name: "successful scan",
			number: func() *number.Number {
				n, _ := number.NewNumber("15556661212")
				return n
			}(),
			mocks: func(s *mocks.NumverifySupplier) {
				s.On("IsAvailable").Return(true)
				s.On("Validate", "15556661212").Return(&suppliers.NumverifyValidateResponse{
					Valid:               true,
					Number:              "test",
					LocalFormat:         "test",
					InternationalFormat: "test",
					CountryPrefix:       "test",
					CountryCode:         "test",
					CountryName:         "test",
					Location:            "test",
					Carrier:             "test",
					LineType:            "test",
				}, nil).Once()
			},
			expected: map[string]interface{}{
				"numverify": NumverifyScannerResponse{
					Valid:               true,
					Number:              "test",
					LocalFormat:         "test",
					InternationalFormat: "test",
					CountryPrefix:       "test",
					CountryCode:         "test",
					CountryName:         "test",
					Location:            "test",
					Carrier:             "test",
					LineType:            "test",
				},
			},
			wantErrors: map[string]error{},
		},
		{
			name: "failed scan",
			number: func() *number.Number {
				n, _ := number.NewNumber("15556661212")
				return n
			}(),
			mocks: func(s *mocks.NumverifySupplier) {
				s.On("IsAvailable").Return(true)
				s.On("Validate", "15556661212").Return(nil, dummyError).Once()
			},
			expected: map[string]interface{}{},
			wantErrors: map[string]error{
				"numverify": dummyError,
			},
		},
		{
			name: "should not run",
			number: func() *number.Number {
				n, _ := number.NewNumber("15556661212")
				return n
			}(),
			mocks: func(s *mocks.NumverifySupplier) {
				s.On("IsAvailable").Return(false)
			},
			expected:   map[string]interface{}{},
			wantErrors: map[string]error{},
		},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			numverifySupplierMock := &mocks.NumverifySupplier{}
			tt.mocks(numverifySupplierMock)

			scanner := NewNumverifyScanner(numverifySupplierMock)
			remote := NewLibrary()
			remote.AddScanner(scanner)

			got, errs := remote.Scan(tt.number)
			if len(tt.wantErrors) > 0 {
				assert.Equal(t, tt.wantErrors, errs)
			} else {
				assert.Len(t, errs, 0)
			}
			assert.Equal(t, tt.expected, got)

			numverifySupplierMock.AssertExpectations(t)
		})
	}
}
