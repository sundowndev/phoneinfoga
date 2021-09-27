package api

import (
	"gopkg.in/sundowndev/phoneinfoga.v2/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"gopkg.in/sundowndev/phoneinfoga.v2/scanners"
)

// JSONResponse is the default API response type
type JSONResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
}

type scanURL struct {
	Number uint `uri:"number" binding:"required,min=2"`
}

// ValidateScanURL validates scan URLs
func ValidateScanURL(c *gin.Context) {
	var v scanURL

	if err := c.ShouldBindUri(&v); err != nil {
		errorHandling(c, "the given phone number is not valid")
		return
	}

	number, err := scanners.LocalScan(c.Param("number"))

	if err != nil {
		utils.LoggerService.Errorln("Validation error", c.Param("number"), err.Error())
		errorHandling(c, err.Error())
		return
	}

	c.Set("number", number)
}

func errorHandling(c *gin.Context, msg string) {
	c.JSON(http.StatusBadRequest, JSONResponse{Success: false, Error: msg})
	c.Abort()
}
