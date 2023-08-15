package api

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWrapHandler(t *testing.T) {
	type expectedResult struct {
		Code    int
		Headers http.Header
		Body    string
	}

	testcases := []struct {
		Name     string
		Handler  HandlerFunc
		Expected expectedResult
	}{
		{
			Name: "test basic handler",
			Handler: func(ctx *gin.Context) *Response {
				return &Response{
					Code: 200,
					Headers: http.Header{
						"My-Header": []string{"val1", "val2"},
					},
					Data: map[string]string{"msg": "test"},
					JSON: true,
				}
			},
			Expected: expectedResult{
				Code: 200,
				Headers: http.Header{
					"Content-Type": []string{"application/json; charset=utf-8"},
					"My-Header":    []string{"val1", "val2"},
				},
				Body: `{"msg":"test"}`,
			},
		},
		{
			Name: "test panic in handler",
			Handler: func(ctx *gin.Context) *Response {
				panic("dummy panic")
			},
			Expected: expectedResult{
				Code:    500,
				Headers: http.Header{"Content-Type": []string{"application/json; charset=utf-8"}},
				Body:    `{"error":"Unknown error"}`,
			},
		},
		{
			Name: "test byte response in handler",
			Handler: func(ctx *gin.Context) *Response {
				return &Response{
					Code:    200,
					Headers: http.Header{},
					Data:    []byte("test"),
				}
			},
			Expected: expectedResult{
				Code:    200,
				Headers: http.Header{},
				Body:    `test`,
			},
		},
		{
			Name: "test unknown response in handler",
			Handler: func(ctx *gin.Context) *Response {
				return &Response{
					Code:    200,
					Headers: http.Header{},
					Data:    23,
				}
			},
			Expected: expectedResult{
				Code:    200,
				Headers: http.Header{},
				Body:    ``,
			},
		},
		{
			Name: "test nil response in handler",
			Handler: func(ctx *gin.Context) *Response {
				ctx.Writer.WriteHeader(403)
				_, _ = ctx.Writer.Write([]byte("test"))
				return nil
			},
			Expected: expectedResult{
				Code:    403,
				Headers: http.Header{},
				Body:    `test`,
			},
		},
	}

	for _, tt := range testcases {
		t.Run(tt.Name, func(t *testing.T) {
			e := gin.New()
			e.GET("/test", WrapHandler(tt.Handler))

			req, err := http.NewRequest(http.MethodGet, "/test", nil)
			if err != nil {
				t.Fatal(err)
			}
			w := httptest.NewRecorder()
			e.ServeHTTP(w, req)

			assert.Equal(t, tt.Expected.Code, w.Code)
			assert.Equal(t, tt.Expected.Headers, w.Header())
			assert.Equal(t, tt.Expected.Body, w.Body.String())
		})
	}
}
