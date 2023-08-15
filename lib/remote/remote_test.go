package remote

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/sundowndev/phoneinfoga/v2/lib/filter"
	"github.com/sundowndev/phoneinfoga/v2/lib/number"
	"github.com/sundowndev/phoneinfoga/v2/mocks"
	"testing"
)

func TestRemoteLibrary_SuccessScan(t *testing.T) {
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
	fakeScanner.On("DryRun", *num).Return(nil).Once()
	fakeScanner.On("Name").Return("fake").Times(2)
	fakeScanner.On("Run", *num).Return(fakeScannerResponse{Valid: true}, nil).Once()

	fakeScanner2 := &mocks.Scanner{}
	fakeScanner2.On("DryRun", *num).Return(nil).Once()
	fakeScanner2.On("Name").Return("fake2").Times(2)
	fakeScanner2.On("Run", *num).Return(fakeScannerResponse{Valid: false}, nil).Once()

	lib := NewLibrary(filter.NewEngine())

	lib.AddScanner(fakeScanner)
	lib.AddScanner(fakeScanner2)

	result, errs := lib.Scan(num)
	assert.Equal(t, expected, result)
	assert.Equal(t, map[string]error{}, errs)

	fakeScanner.AssertExpectations(t)
	fakeScanner2.AssertExpectations(t)
}

func TestRemoteLibrary_FailedScan(t *testing.T) {
	num, err := number.NewNumber("15556661212")
	if err != nil {
		t.Fatal(err)
	}

	dummyError := errors.New("test")

	fakeScanner := &mocks.Scanner{}
	fakeScanner.On("DryRun", *num).Return(nil).Once()
	fakeScanner.On("Name").Return("fake").Times(2)
	fakeScanner.On("Run", *num).Return(nil, dummyError).Once()

	lib := NewLibrary(filter.NewEngine())

	lib.AddScanner(fakeScanner)

	result, errs := lib.Scan(num)
	assert.Equal(t, map[string]interface{}{}, result)
	assert.Equal(t, map[string]error{"fake": dummyError}, errs)

	fakeScanner.AssertExpectations(t)
}

func TestRemoteLibrary_EmptyScan(t *testing.T) {
	num, err := number.NewNumber("15556661212")
	if err != nil {
		t.Fatal(err)
	}

	fakeScanner := &mocks.Scanner{}
	fakeScanner.On("Name").Return("mockscanner").Times(2)
	fakeScanner.On("DryRun", *num).Return(errors.New("dummy error")).Once()

	lib := NewLibrary(filter.NewEngine())

	lib.AddScanner(fakeScanner)

	result, errs := lib.Scan(num)
	assert.Equal(t, map[string]interface{}{}, result)
	assert.Equal(t, map[string]error{}, errs)

	fakeScanner.AssertExpectations(t)
}

func TestRemoteLibrary_PanicRun(t *testing.T) {
	num, err := number.NewNumber("15556661212")
	if err != nil {
		t.Fatal(err)
	}

	fakeScanner := &mocks.Scanner{}
	fakeScanner.On("Name").Return("fake")
	fakeScanner.On("DryRun", *num).Return(nil).Once()
	fakeScanner.On("Run", *num).Panic("dummy panic").Once()

	lib := NewLibrary(filter.NewEngine())

	lib.AddScanner(fakeScanner)

	result, errs := lib.Scan(num)
	assert.Equal(t, map[string]interface{}{}, result)
	assert.Equal(t, map[string]error{"fake": errors.New("panic occurred while running scan, see debug logs")}, errs)

	fakeScanner.AssertExpectations(t)
}

func TestRemoteLibrary_PanicDryRun(t *testing.T) {
	num, err := number.NewNumber("15556661212")
	if err != nil {
		t.Fatal(err)
	}

	fakeScanner := &mocks.Scanner{}
	fakeScanner.On("Name").Return("fake")
	fakeScanner.On("DryRun", *num).Panic("dummy panic").Once()

	lib := NewLibrary(filter.NewEngine())

	lib.AddScanner(fakeScanner)

	result, errs := lib.Scan(num)
	assert.Equal(t, map[string]interface{}{}, result)
	assert.Equal(t, map[string]error{"fake": errors.New("panic occurred while running scan, see debug logs")}, errs)

	fakeScanner.AssertExpectations(t)
}

func TestRemoteLibrary_GetAllScanners(t *testing.T) {
	fakeScanner := &mocks.Scanner{}
	fakeScanner.On("Name").Return("fake")

	fakeScanner2 := &mocks.Scanner{}
	fakeScanner2.On("Name").Return("fake2")

	lib := NewLibrary(filter.NewEngine())

	lib.AddScanner(fakeScanner)
	lib.AddScanner(fakeScanner2)

	assert.Equal(t, []Scanner{fakeScanner, fakeScanner2}, lib.GetAllScanners())
}

func TestRemoteLibrary_AddIgnoredScanner(t *testing.T) {
	fakeScanner := &mocks.Scanner{}
	fakeScanner.On("Name").Return("fake")

	fakeScanner2 := &mocks.Scanner{}
	fakeScanner2.On("Name").Return("fake2")

	f := filter.NewEngine()
	f.AddRule("fake2")
	lib := NewLibrary(f)

	lib.AddScanner(fakeScanner)
	lib.AddScanner(fakeScanner2)

	assert.Equal(t, []Scanner{fakeScanner}, lib.GetAllScanners())
}
