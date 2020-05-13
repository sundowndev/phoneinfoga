package api

import (
	"strings"
)

func successResponse(msg ...string) JSONResponse {
	var message string

	if len(msg) == 0 {
		message = ""
	} else {
		message = strings.Join(msg, "")
	}

	return JSONResponse{
		Success: true,
		Error:   message,
	}
}

func errorResponse(msg ...string) JSONResponse {
	var message string

	if len(msg) == 0 {
		message = "An error occurred"
	} else {
		message = strings.Join(msg, "")
	}

	return JSONResponse{
		Success: false,
		Error:   message,
	}
}
