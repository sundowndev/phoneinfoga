package api

import (
	"github.com/gin-gonic/gin"
	"github.com/sundowndev/phoneinfoga/pkg/config"
	"github.com/sundowndev/phoneinfoga/pkg/scanners"
)

type scanResultResponse struct {
	JSONResponse
	Result interface{} `json:"result"`
}

type listNumbersResponse struct {
	JSONResponse
	Numbers interface{} `json:"numbers"`
}

func getAllNumbers(c *gin.Context) {
	c.JSON(200, listNumbersResponse{
		JSONResponse: JSONResponse{Success: true},
		Numbers:      []scanners.Number{},
	})
}

func validate(c *gin.Context) {
	number := c.Param("number")

	_, err := scanners.LocalScan(number)

	if err != nil {
		c.JSON(400, errorResponse(err.Error()))
		return
	}

	c.JSON(200, successResponse("The number is valid"))
}

func localScan(c *gin.Context) {
	number := c.Param("number")

	result, err := scanners.LocalScan(number)

	if err != nil {
		c.JSON(500, errorResponse(err.Error()))
		return
	}

	c.JSON(200, scanResultResponse{
		JSONResponse: JSONResponse{Success: true},
		Result:       result,
	})
}

func numverifyScan(c *gin.Context) {
	number := c.Param("number")

	n, err := scanners.LocalScan(number)

	if err != nil {
		c.JSON(500, errorResponse("The number is not valid"))
		return
	}

	result, err := scanners.NumverifyScan(n)

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
	number := c.Param("number")

	n, err := scanners.LocalScan(number)

	if err != nil {
		c.JSON(500, errorResponse("The number is not valid"))
		return
	}

	result := scanners.GoogleSearchScan(n)

	c.JSON(200, scanResultResponse{
		JSONResponse: JSONResponse{Success: true},
		Result:       result,
	})
}

func ovhScan(c *gin.Context) {
	number := c.Param("number")

	n, err := scanners.LocalScan(number)

	if err != nil {
		c.JSON(500, errorResponse("The number is not valid"))
		return
	}

	result, err := scanners.OVHScan(n)

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
