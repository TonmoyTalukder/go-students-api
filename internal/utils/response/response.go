package response

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Response struct {
	Status string `json:"status"`
	Error string `json:"error"`
}

type JsonResponse struct {
	Status  string      `json:"status"`
	Code    int         `json:"code"`
	Message string      `json:"message,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

const (
	StatusOK = "OK"
	StatusError = "Error"
)

func WriteJson(w http.ResponseWriter, status int, payload JsonResponse) error{   // data interface{}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(payload)
}

func SuccessResponse(code int, message string, data interface{}, meta interface{}) JsonResponse {
	return JsonResponse{
		Status:  StatusOK,
		Code:    code,
		Message: message,
		Data:    data,
		Meta:    meta,
	}
}

func ErrorResponse(code int, message string, err error) JsonResponse {
	return JsonResponse{
		Status:  StatusError,
		Code:    code,
		Message: message,
		Error:   err.Error(),
	}
}

func GeneralError(err error) Response {
	return Response {
		Status: StatusError,
		Error: err.Error(),
	}
}

func ValidationErr(errs validator.ValidationErrors) JsonResponse {
	var errMsgs []string

	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("Field %s is required field.", err.Field()))
		
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("Field %s is invalid.", err.Field()))
		}
	}

	return JsonResponse{
		Status: StatusError,
		Error: strings.Join(errMsgs, ", "),
	}
}