package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/sundowndev/phoneinfoga/v2/lib/number"
	"github.com/sundowndev/phoneinfoga/v2/web/v2/api"
	"net/http"
)

type AddNumberInput struct {
	Number string `json:"number" binding:"number,required"`
}

type AddNumberResponse struct {
	Valid         bool   `json:"valid"`
	RawLocal      string `json:"rawLocal"`
	Local         string `json:"local"`
	E164          string `json:"e164"`
	International string `json:"international"`
	CountryCode   int32  `json:"countryCode"`
	Country       string `json:"country"`
	Carrier       string `json:"carrier"`
}

// AddNumber is an HTTP handler
// @ID AddNumber
// @Tags Numbers
// @Summary Add a new number.
// @Description This route returns information about a given phone number.
// @Accept  json
// @Produce  json
// @Param request body AddNumberInput true "Request body"
// @Success 200 {object} AddNumberResponse
// @Success 500 {object} api.ErrorResponse
// @Router /v2/numbers [post]
func AddNumber(ctx *gin.Context) *api.Response {
	var input AddNumberInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		return &api.Response{
			Code: http.StatusBadRequest,
			JSON: true,
			Data: api.ErrorResponse{Error: "Invalid phone number: please provide an integer without any special chars"},
		}
	}

	num, err := number.NewNumber(input.Number)
	if err != nil {
		return &api.Response{
			Code: http.StatusBadRequest,
			JSON: true,
			Data: api.ErrorResponse{Error: err.Error()},
		}
	}

	return &api.Response{
		Code: http.StatusOK,
		JSON: true,
		Data: AddNumberResponse{
			Valid:         num.Valid,
			RawLocal:      num.RawLocal,
			Local:         num.Local,
			E164:          num.E164,
			International: num.International,
			CountryCode:   num.CountryCode,
			Country:       num.Country,
			Carrier:       num.Carrier,
		},
	}
}
