package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type HandlerFunc func(ctx *gin.Context) *Response

type Response struct {
	Code    int
	Headers http.Header
	Data    interface{}
	JSON    bool
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func WrapHandler(h HandlerFunc) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				ctx.AbortWithStatusJSON(500, ErrorResponse{Error: "Unknown error"})
			}
		}()

		res := h(ctx)
		for key, values := range res.Headers {
			for _, val := range values {
				ctx.Header(key, val)
			}
		}
		if res.JSON && res.Data != nil {
			ctx.JSON(res.Code, res.Data)
			return
		}
		ctx.Writer.WriteHeader(res.Code)
		if _, ok := res.Data.([]byte); ok {
			_, _ = ctx.Writer.Write(res.Data.([]byte))
		}
	}
}
