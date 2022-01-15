package remote

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/sundowndev/phoneinfoga/v2/lib/number"
	"github.com/sundowndev/phoneinfoga/v2/lib/remote/suppliers"
	"github.com/sundowndev/phoneinfoga/v2/mocks"
	"testing"
)

func TestOVHScanner(t *testing.T) {
	dummyError := errors.New("dummy")

	dummyNumber, _ := number.NewNumber("33365174444")

	testcases := []struct {
		name       string
		number     *number.Number
		mocks      func(s *mocks.OVHSupplier)
		expected   map[string]interface{}
		wantErrors map[string]error
	}{
		{
			name: "successful scan",
			number: func() *number.Number {
				return dummyNumber
			}(),
			mocks: func(s *mocks.OVHSupplier) {
				s.On("Search", *dummyNumber).Return(&suppliers.OVHScannerResponse{
					Found: false,
				}, nil).Once()
			},
			expected: map[string]interface{}{
				"ovh": OVHScannerResponse{
					Found: false,
				},
			},
			wantErrors: map[string]error{},
		},
		{
			name: "failed scan",
			number: func() *number.Number {
				return dummyNumber
			}(),
			mocks: func(s *mocks.OVHSupplier) {
				s.On("Search", *dummyNumber).Return(nil, dummyError).Once()
			},
			expected: map[string]interface{}{},
			wantErrors: map[string]error{
				"ovh": dummyError,
			},
		},
		{
			name: "country not supported",
			number: func() *number.Number {
				num, _ := number.NewNumber("15556661212")
				return num
			}(),
			mocks:      func(s *mocks.OVHSupplier) {},
			expected:   map[string]interface{}{},
			wantErrors: map[string]error{},
		},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			OVHSupplierMock := &mocks.OVHSupplier{}
			tt.mocks(OVHSupplierMock)

			scanner := NewOVHScanner(OVHSupplierMock)
			remote := NewLibrary()
			remote.AddScanner(scanner)

			got, errs := remote.Scan(tt.number)
			if len(tt.wantErrors) > 0 {
				assert.Equal(t, tt.wantErrors, errs)
			} else {
				assert.Len(t, errs, 0)
			}
			assert.Equal(t, tt.expected, got)

			OVHSupplierMock.AssertExpectations(t)
		})
	}
}
