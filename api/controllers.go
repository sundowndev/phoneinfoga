package api

import (
	"github.com/gin-gonic/gin"
	"gopkg.in/sundowndev/phoneinfoga.v2/config"
	"gopkg.in/sundowndev/phoneinfoga.v2/scanners"
)

type scanResultResponse struct {
	JSONResponse
	Result interface{} `json:"result"`
}

type getAllNumbersResponse struct {
	JSONResponse
	Numbers []scanners.Number `json:"numbers"`
}

type healthResponse struct {
	Success bool   `json:"success"`
	Version string `json:"version"`
	Commit  string `json:"commit"`
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
	c.JSON(200, getAllNumbersResponse{
		JSONResponse: JSONResponse{Success: true},
		Numbers:      []scanners.Number{},
	})
}

// @ID validate
// @Tags Numbers
// @Summary Check if a number is valid and possible.
// @Produce  json
// @Success 200 {object} JSONResponse
// @Success 400 {object} JSONResponse
// @Router /numbers/{number}/validate [get]
// @Param number path string true "Input phone number" validate(required)
func validate(c *gin.Context) {
	c.JSON(200, successResponse("The number is valid"))
}

// @ID localScan
// @Tags Numbers
// @Summary Perform a scan using local phone number library.
// @Produce  json
// @Success 200 {object} scanResultResponse{result=scanners.Number}
// @Success 400 {object} JSONResponse
// @Router /numbers/{number}/scan/local [get]
// @Param number path string true "Input phone number" validate(required)
func localScan(c *gin.Context) {
	result, _ := c.Get("number")

	c.JSON(200, scanResultResponse{
		JSONResponse: JSONResponse{Success: true},
		Result:       result.(*scanners.Number),
	})
}

// @ID numverifyScan
// @Tags Numbers
// @Summary Perform a scan using Numverify's API.
// @Produce  json
// @Success 200 {object} scanResultResponse{result=scanners.NumverifyScannerResponse}
// @Success 400 {object} JSONResponse
// @Router /numbers/{number}/scan/numverify [get]
// @Param number path string true "Input phone number" validate(required)
func numverifyScan(c *gin.Context) {
	number, _ := c.Get("number")

	result, err := scanners.NumverifyScan(number.(*scanners.Number))
	if err != nil {
		c.JSON(500, errorResponse(err.Error()))
		return
	}

	c.JSON(200, scanResultResponse{
		JSONResponse: JSONResponse{Success: true},
		Result:       result,
	})
}

// @ID googleSearchScan
// @Tags Numbers
// @Summary Perform a scan using Google Search engine.
// @Produce  json
// @Success 200 {object} scanResultResponse{result=scanners.GoogleSearchResponse}
// @Success 400 {object} JSONResponse
// @Router /numbers/{number}/scan/googlesearch [get]
// @Param number path string true "Input phone number" validate(required)
func googleSearchScan(c *gin.Context) {
	number, _ := c.Get("number")

	result := scanners.GoogleSearchScan(number.(*scanners.Number))

	c.JSON(200, scanResultResponse{
		JSONResponse: JSONResponse{Success: true},
		Result:       result,
	})
}

// @ID ovhScan
// @Tags Numbers
// @Summary Perform a scan using OVH's API.
// @Produce  json
// @Success 200 {object} scanResultResponse{result=scanners.OVHScannerResponse}
// @Success 400 {object} JSONResponse
// @Router /numbers/{number}/scan/ovh [get]
// @Param number path string true "Input phone number" validate(required)
func ovhScan(c *gin.Context) {
	number, _ := c.Get("number")

	result, err := scanners.OVHScan(number.(*scanners.Number))

	if err != nil {
		c.JSON(500, errorResponse())
		return
	}

	c.JSON(200, scanResultResponse{
		JSONResponse: JSONResponse{Success: true},
		Result:       result,
	})
}

// @ID healthCheck
// @Tags default
// @Summary Check if service is healthy.
// @Produce  json
// @Success 200 {object} healthResponse
// @Success 500 {object} JSONResponse
// @Router / [get]
func healthHandler(c *gin.Context) {
	c.JSON(200, healthResponse{
		Success: true,
		Version: config.Version,
		Commit:  config.Commit,
	})
}
