package utils

import "github.com/gin-gonic/gin"

type Response struct {
	StatusCode int         `json:"statusCode"`
	Method     string      `json:"method"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
}

type ErrorResponse struct {
	StatusCode int         `json:"statusCode"`
	Method     string      `json:"method"`
	Error      interface{} `json:"error"`
}

func APIResponse(ctx *gin.Context, ErrMsg string, StatusCode int, Method string, Data interface{}) {
	jsonResponse := Response{
		StatusCode: StatusCode,
		Method:     Method,
		Message:    ErrMsg,
		Data:       Data,
	}

	if StatusCode >= 400 {
		ctx.JSON(StatusCode, jsonResponse)
		defer ctx.AbortWithStatus(StatusCode)
	} else {
		ctx.JSON(StatusCode, jsonResponse)
	}
}

func ValidatorErrorResponse(ctx *gin.Context, StatusCode int, Method string, Error interface{}) {
	errResponse := ErrorResponse{
		StatusCode: StatusCode,
		Method:     Method,
		Error:      Error,
	}

	ctx.JSON(StatusCode, errResponse)
	defer ctx.AbortWithStatus(StatusCode)
}
