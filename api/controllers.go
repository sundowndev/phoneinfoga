package api

import (
	"github.com/gin-gonic/gin"
	"gopkg.in/sundowndev/phoneinfoga.v2/pkg/config"
	"gopkg.in/sundowndev/phoneinfoga.v2/pkg/scanners"
)

type scanResultResponse struct {
	JSONResponse
	Result interface{} `json:"result"`
}

type listNumbersResponse struct {
	JSONResponse
	Numbers []scanners.Number `json:"numbers"`
}

func getAllNumbers(c *gin.Context) {
	c.JSON(200, listNumbersResponse{
		JSONResponse: JSONResponse{Success: true},
		Numbers:      []scanners.Number{},
	})
}

func validate(c *gin.Context) {
	c.JSON(200, successResponse("The number is valid"))
}

func localScan(c *gin.Context) {
	result, _ := c.Get("number")

	c.JSON(200, scanResultResponse{
		JSONResponse: JSONResponse{Success: true},
		Result:       result.(*scanners.Number),
	})
}

func numverifyScan(c *gin.Context) {
	number, _ := c.Get("number")

	result, err := scanners.NumverifyScan(number.(*scanners.Number))

	if err != nil {
		c.JSON(500, errorResponse())
		return
	}

	c.JSON(200, scanResultResponse{
		JSONResponse: JSONResponse{Success: true},
		Result:       result,
	})
}

func googleSearchScan(c *gin.Context) {
	number, _ := c.Get("number")

	result := scanners.GoogleSearchScan(number.(*scanners.Number))

	c.JSON(200, scanResultResponse{
		JSONResponse: JSONResponse{Success: true},
		Result:       result,
	})
}

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

func healthHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"success": true,
		"version": config.Version,
	})
}
