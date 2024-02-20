package web

import (
	"github.com/gin-gonic/gin"
	"github.com/sundowndev/phoneinfoga/v2/build"
	"github.com/sundowndev/phoneinfoga/v2/lib/number"
	"github.com/sundowndev/phoneinfoga/v2/lib/remote"
	"github.com/sundowndev/phoneinfoga/v2/lib/remote/suppliers"
	"github.com/sundowndev/phoneinfoga/v2/web/errors"
	"net/http"
)

type ScanResultResponse struct {
	JSONResponse
	Result interface{} `json:"result"`
}

type getAllNumbersResponse struct {
	JSONResponse
	Numbers []number.Number `json:"numbers"`
}

type healthResponse struct {
	Success bool   `json:"success"`
	Version string `json:"version"`
	Commit  string `json:"commit"`
	Demo    bool   `json:"demo"`
}

// @ID getAllNumbers
// @Tags Numbers
// @Summary Fetch all previously scanned numbers.
// @Description This route is actually not used yet.
// @Deprecated
// @Produce  json
// @Success 200 {object} getAllNumbersResponse
// @Router /numbers [get]
func getAllNumbers(c *gin.Context) {
	c.JSON(http.StatusOK, getAllNumbersResponse{
		JSONResponse: JSONResponse{Success: true},
		Numbers:      []number.Number{},
	})
}

// @ID validate
// @Tags Numbers
// @Summary Check if a number is valid and possible.
// @Produce  json
// @Deprecated
// @Success 200 {object} JSONResponse
// @Success 400 {object} JSONResponse
// @Router /numbers/{number}/validate [get]
// @Param number path string true "Input phone number" validate(required)
func validate(c *gin.Context) {
	_, err := number.NewNumber(c.Param("number"))
	if err != nil {
		handleError(c, errors.NewBadRequest(err))
		return
	}
	c.JSON(http.StatusOK, successResponse("The number is valid"))
}

// @ID localScan
// @Tags Numbers
// @Summary Perform a scan using local phone number library.
// @Produce  json
// @Deprecated
// @Success 200 {object} ScanResultResponse{result=number.Number}
// @Success 400 {object} JSONResponse
// @Router /numbers/{number}/scan/local [get]
// @Param number path string true "Input phone number" validate(required)
func localScan(c *gin.Context) {
	num, err := number.NewNumber(c.Param("number"))
	if err != nil {
		handleError(c, errors.NewBadRequest(err))
		return
	}

	result, err := remote.NewLocalScanner().Run(*num, make(remote.ScannerOptions))
	if err != nil {
		handleError(c, errors.NewInternalError(err))
		return
	}

	c.JSON(http.StatusOK, ScanResultResponse{
		JSONResponse: JSONResponse{Success: true},
		Result:       result,
	})
}

// @ID numverifyScan
// @Tags Numbers
// @Summary Perform a scan using Numverify's API.
// @Deprecated
// @Produce  json
// @Success 200 {object} ScanResultResponse{result=remote.NumverifyScannerResponse}
// @Success 400 {object} JSONResponse
// @Router /numbers/{number}/scan/numverify [get]
// @Param number path string true "Input phone number" validate(required)
func numverifyScan(c *gin.Context) {
	num, err := number.NewNumber(c.Param("number"))
	if err != nil {
		handleError(c, errors.NewBadRequest(err))
		return
	}

	result, err := remote.NewNumverifyScanner(suppliers.NewNumverifySupplier()).Run(*num, make(remote.ScannerOptions))
	if err != nil {
		handleError(c, errors.NewInternalError(err))
		return
	}

	c.JSON(http.StatusOK, ScanResultResponse{
		JSONResponse: JSONResponse{Success: true},
		Result:       result,
	})
}

// @ID googleSearchScan
// @Tags Numbers
// @Summary Perform a scan using Google Search engine.
// @Deprecated
// @Produce  json
// @Success 200 {object} ScanResultResponse{result=remote.GoogleSearchResponse}
// @Success 400 {object} JSONResponse
// @Router /numbers/{number}/scan/googlesearch [get]
// @Param number path string true "Input phone number" validate(required)
func googleSearchScan(c *gin.Context) {
	num, err := number.NewNumber(c.Param("number"))
	if err != nil {
		handleError(c, errors.NewBadRequest(err))
		return
	}

	result, err := remote.NewGoogleSearchScanner().Run(*num, make(remote.ScannerOptions))
	if err != nil {
		handleError(c, errors.NewInternalError(err))
		return
	}

	c.JSON(http.StatusOK, ScanResultResponse{
		JSONResponse: JSONResponse{Success: true},
		Result:       result,
	})
}

// @ID ovhScan
// @Tags Numbers
// @Summary Perform a scan using OVH's API.
// @Deprecated
// @Produce  json
// @Success 200 {object} ScanResultResponse{result=remote.OVHScannerResponse}
// @Success 400 {object} JSONResponse
// @Router /numbers/{number}/scan/ovh [get]
// @Param number path string true "Input phone number" validate(required)
func ovhScan(c *gin.Context) {
	num, err := number.NewNumber(c.Param("number"))
	if err != nil {
		handleError(c, errors.NewBadRequest(err))
		return
	}

	result, err := remote.NewOVHScanner(suppliers.NewOVHSupplier()).Run(*num, make(remote.ScannerOptions))
	if err != nil {
		handleError(c, errors.NewInternalError(err))
		return
	}

	c.JSON(http.StatusOK, ScanResultResponse{
		JSONResponse: JSONResponse{Success: true},
		Result:       result,
	})
}

// @ID healthCheck
// @Tags General
// @Summary Check if service is healthy.
// @Produce  json
// @Success 200 {object} healthResponse
// @Success 500 {object} JSONResponse
// @Router / [get]
func healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, healthResponse{
		Success: true,
		Version: build.Version,
		Commit:  build.Commit,
		Demo:    build.IsDemo(),
	})
}
