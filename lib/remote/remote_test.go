package remote

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/sundowndev/phoneinfoga/v2/lib/filter"
	"github.com/sundowndev/phoneinfoga/v2/lib/number"
	"github.com/sundowndev/phoneinfoga/v2/mocks"
	"testing"
)

func TestRemoteLibrarySuccessScan(t *testing.T) {
	type fakeScannerResponse struct {
		Valid bool
	}

	expected := map[string]interface{}{
		"fake":  fakeScannerResponse{Valid: true},
		"fake2": fakeScannerResponse{Valid: false},
	}

	num, err := number.NewNumber("15556661212")
	if err != nil {
		t.Fatal(err)
	}

	fakeScanner := &mocks.Scanner{}
	fakeScanner.On("ShouldRun", *num).Return(true).Once()
	fakeScanner.On("Name").Return("fake").Times(2)
	fakeScanner.On("Scan", *num).Return(fakeScannerResponse{Valid: true}, nil).Once()

	fakeScanner2 := &mocks.Scanner{}
	fakeScanner2.On("ShouldRun", *num).Return(true).Once()
	fakeScanner2.On("Name").Return("fake2").Times(2)
	fakeScanner2.On("Scan", *num).Return(fakeScannerResponse{Valid: false}, nil).Once()

	lib := NewLibrary(filter.NewEngine())

	lib.AddScanner(fakeScanner)
	lib.AddScanner(fakeScanner2)

	result, errs := lib.Scan(num)
	assert.Equal(t, expected, result)
	assert.Equal(t, map[string]error{}, errs)

	fakeScanner.AssertExpectations(t)
	fakeScanner2.AssertExpectations(t)
}

func TestRemoteLibraryFailedScan(t *testing.T) {
	num, err := number.NewNumber("15556661212")
	if err != nil {
		t.Fatal(err)
	}

	dummyError := errors.New("test")

	fakeScanner := &mocks.Scanner{}
	fakeScanner.On("ShouldRun", *num).Return(true).Once()
	fakeScanner.On("Name").Return("fake").Times(2)
	fakeScanner.On("Scan", *num).Return(nil, dummyError).Once()

	lib := NewLibrary(filter.NewEngine())

	lib.AddScanner(fakeScanner)

	result, errs := lib.Scan(num)
	assert.Equal(t, map[string]interface{}{}, result)
	assert.Equal(t, map[string]error{"fake": dummyError}, errs)

	fakeScanner.AssertExpectations(t)
}

func TestRemoteLibraryEmptyScan(t *testing.T) {
	num, err := number.NewNumber("15556661212")
	if err != nil {
		t.Fatal(err)
	}

	fakeScanner := &mocks.Scanner{}
	fakeScanner.On("Name").Return("mockscanner").Times(2)
	fakeScanner.On("ShouldRun", *num).Return(false).Once()

	lib := NewLibrary(filter.NewEngine())

	lib.AddScanner(fakeScanner)

	result, errs := lib.Scan(num)
	assert.Equal(t, map[string]interface{}{}, result)
	assert.Equal(t, map[string]error{}, errs)

	fakeScanner.AssertExpectations(t)
}
