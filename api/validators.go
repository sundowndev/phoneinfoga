package api

import (
	"github.com/gin-gonic/gin"
	"github.com/sundowndev/phoneinfoga/pkg/utils"
)

// JSONResponse is the default API response type
type JSONResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

type scanURL struct {
	Number uint `uri:"number" binding:"required,min=1,max=999679368229"`
}

// ValidateScanURL validates scan URLs
func ValidateScanURL(c *gin.Context) {
	number := c.Param("number")

	if valid := utils.IsValid(number); valid != true {
		c.JSON(500, errorResponse("The number is not valid"))
		c.Abort()
	}
}
