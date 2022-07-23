package web

import (
	errors2 "errors"
	"github.com/gin-gonic/gin"
	"github.com/sundowndev/phoneinfoga/v2/web/errors"
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
		handleError(c, errors.NewBadRequest(errors2.New("the given phone number is not valid")))
	}
}
