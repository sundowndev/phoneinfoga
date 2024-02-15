package remote_test

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/sundowndev/phoneinfoga/v2/lib/filter"
	"github.com/sundowndev/phoneinfoga/v2/lib/number"
	"github.com/sundowndev/phoneinfoga/v2/lib/remote"
	"github.com/sundowndev/phoneinfoga/v2/lib/remote/suppliers"
	"github.com/sundowndev/phoneinfoga/v2/mocks"
	"testing"
)

func TestNumverifyScanner_Metadata(t *testing.T) {
	scanner := remote.NewNumverifyScanner(&mocks.NumverifySupplier{})
	assert.Equal(t, remote.Numverify, scanner.Name())
	assert.NotEmpty(t, scanner.Description())
}

func TestNumverifyScanner(t *testing.T) {
	dummyError := errors.New("dummy")

	testcases := []struct {
		name       string
		number     *number.Number
		opts       remote.ScannerOptions
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
				s.On("Validate", "15556661212", "").Return(&suppliers.NumverifyValidateResponse{
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
				"numverify": remote.NumverifyScannerResponse{
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
				s.On("Validate", "15556661212", "").Return(nil, dummyError).Once()
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
		{
			name: "should run with options defined",
			opts: remote.ScannerOptions{
				"api_key": "secret",
			},
			number: func() *number.Number {
				n, _ := number.NewNumber("15556661212")
				return n
			}(),
			mocks: func(s *mocks.NumverifySupplier) {
				s.On("Validate", "15556661212", "secret").Return(&suppliers.NumverifyValidateResponse{
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
				"numverify": remote.NumverifyScannerResponse{
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
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			numverifySupplierMock := &mocks.NumverifySupplier{}
			tt.mocks(numverifySupplierMock)

			scanner := remote.NewNumverifyScanner(numverifySupplierMock)
			lib := remote.NewLibrary(filter.NewEngine())
			lib.AddScanner(scanner)

			got, errs := lib.Scan(tt.number, tt.opts)
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
