package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/easy-cloud-Knet/KWS_Core.git/api/conn"
)

type BaseResponse[T any] struct {
	Information *T      `json:"information,omitempty"`
	Message     string `json:"message"`
	Errors      error `json:"errors,omitempty"`
}

func ResponseGen[T any](message string) *BaseResponse[T] {
	return &BaseResponse[T]{
		Message:     fmt.Sprintf("%s operation", message),
	}
}

// HTTP 요청을 디코딩하는 함수
func HttpDecoder[T any](r *http.Request, param *T) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return conn.ErrorGen(conn.FaildDeEncoding,fmt.Errorf("%w error unmarshaling body into Structure",err))
	}
	defer r.Body.Close()

	if err := json.Unmarshal(body, param); err != nil {
		return conn.ErrorGen(conn.FaildDeEncoding,fmt.Errorf("%w error unmarshaling body into Structure",err))
	}
	return nil
}

func (br *BaseResponse[T]) ResponseWriteErr(w http.ResponseWriter, err error, statusCode int) {
	br.Message += " failed"
	errDesc, ok:= err.(conn.ErrorDescriptor)
	if !ok{
		http.Error(w, br.Message, http.StatusInternalServerError)
		return
	}
	br.Errors = errDesc

	data, marshalErr := json.Marshal(br)
	if marshalErr != nil {
		fmt.Println(marshalErr)
		http.Error(w, "failed to marshal error response", http.StatusInternalServerError)
		return
	}
	fmt.Println("error occured in RESPONSE ERR ", br)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(data)
}

func (br *BaseResponse[T]) ResponseWriteOK(w http.ResponseWriter, info *T) {
	br.Message += " success"
	br.Information = info
	data, err := json.Marshal(br)
	if err != nil {
		http.Error(w, "failed to marshal success response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
