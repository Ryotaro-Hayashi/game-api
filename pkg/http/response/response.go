//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=./mock_$GOPACKAGE/mock_$GOFILE

package response

import (
	"20dojo-online/pkg/myerror"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

type HttpResponse struct{}

type HttpResponseInterface interface {
	Success(writer http.ResponseWriter, response interface{})
	Failed(writer http.ResponseWriter, err error)
}

func NewHttpResponse() *HttpResponse {
	return &HttpResponse{}
}

// Success HTTPコード:200 正常終了を処理する
func (hr *HttpResponse) Success(writer http.ResponseWriter, response interface{}) {
	if response == nil {
		return
	}
	data, err := json.Marshal(response)
	if err != nil {
		log.Println(err)
		internalServerError(writer, "marshal error")
		return
	}
	writer.Write(data)
}

// Failed リクエスト失敗時のエラー処理
func (hr *HttpResponse) Failed(writer http.ResponseWriter, err error) {
	var appErr myerror.ApplicationError
	if errors.As(err, &appErr) {
		switch appErr.Code {
		case http.StatusBadRequest:
			badRequest(writer, "Bad Request")
		case http.StatusInternalServerError:
			internalServerError(writer, "Internal Server Error")
		}
	} else {
		internalServerError(writer, "Unknown Internal Server Error")
	}
}

// BadRequest HTTPコード:400 BadRequestを処理する
func badRequest(writer http.ResponseWriter, message string) {
	httpError(writer, http.StatusBadRequest, message)
}

// InternalServerError HTTPコード:500 InternalServerErrorを処理する
func internalServerError(writer http.ResponseWriter, message string) {
	httpError(writer, http.StatusInternalServerError, message)
}

// httpError エラー用のレスポンス出力を行う
func httpError(writer http.ResponseWriter, code int, message string) {
	data, _ := json.Marshal(errorResponse{
		Code:    code,
		Message: message,
	})
	writer.WriteHeader(code)
	if data != nil {
		writer.Write(data)
	}
}

type errorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
