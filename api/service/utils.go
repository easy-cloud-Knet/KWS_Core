package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/easy-cloud-Knet/KWS_Core.git/api/conn"
)

// BaseResponse를 제네릭으로 변경
type BaseResponse[T any] struct {
	Information *T      `json:"information,omitempty"`
	Message     string `json:"message"`
	Errors      conn.ErrorDescriptor `json:"errors,omitempty"`
}

// 메시지를 생성하는 함수
func ResponseGen[T any](message string) *BaseResponse[T] {
	return &BaseResponse[T]{
		Information: nil, // 기본값 설정
		Message:     fmt.Sprintf("%s operation", message),

	}
}

// HTTP 요청을 디코딩하는 함수
func HttpDecoder[T any](r *http.Request, param *T) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("error reading request body: %w", err)
	}
	defer r.Body.Close()

	if err := json.Unmarshal(body, param); err != nil {
		return fmt.Errorf("error unmarshaling JSON: %w", err)
	}
	return nil
}

// 실패 응답 처리
func (br *BaseResponse[T]) ResponseWriteErr(w http.ResponseWriter, err error, statusCode int) {
	br.Message += " failed"
	errDesc, ok:= err.(conn.ErrorDescriptor)
	if !ok{
		http.Error(w, "fasdfdfsde", http.StatusInternalServerError)
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

// 성공 응답 처리
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
