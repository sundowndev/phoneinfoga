package web

import (
	"github.com/gin-gonic/gin"
	"github.com/sundowndev/phoneinfoga/v2/web/errors"
)

func handleError(c *gin.Context, e *errors.Error) {
	c.JSON(e.Status(), JSONResponse{Success: false, Error: e.String()})
	c.Abort()
}
