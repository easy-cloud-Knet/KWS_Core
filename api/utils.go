package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	domCon "github.com/easy-cloud-Knet/KWS_Core/DomCon"
	virerr "github.com/easy-cloud-Knet/KWS_Core/error"
)



type BaseResponse[T any] struct {
	Information *T      `json:"information,omitempty"`
	Message     string `json:"message"`
	Errors      *virerr.ErrorDescriptor `json:"errors,omitempty"`
	ErrorDebug  string `json:"errrorDebug,omitempty"`
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
		return virerr.ErrorGen(virerr.FaildDeEncoding,fmt.Errorf("%w error unmarshaling body into Structure",err))
	}
	defer r.Body.Close()

	if err := json.Unmarshal(body, param); err != nil {
		return virerr.ErrorGen(virerr.FaildDeEncoding,fmt.Errorf("%w error unmarshaling body into Structure",err))
	}
	return nil
}

func (br *BaseResponse[T]) ResponseWriteErr(w http.ResponseWriter, err error, statusCode int) {
	br.Message += " failed"
	errDesc, ok:= err.(virerr.ErrorDescriptor)
	if !ok{
		http.Error(w, br.Message, http.StatusInternalServerError)
		return
	}
	br.Errors = &errDesc
	br.ErrorDebug= errDesc.Error()
	data, marshalErr := json.Marshal(br)
	if marshalErr != nil {
		http.Error(w, "failed to marshal error response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(data)
}

func (br *BaseResponse[T]) ResponseWriteOK(w http.ResponseWriter, info *T) {
	br.Message += " success"
	br.Information = info
	br.Errors=nil
	data, err := json.Marshal(br)
	if err != nil {
		http.Error(w, "failed to marshal success response", http.StatusInternalServerError)
		return
	}
	

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}


func (i *InstHandler) domainConGetter()(*domCon.DomListControl, error){
	fmt.Println(i)
	if(i.DomainControl == nil){
		return nil, fmt.Errorf("domainController not initialled yet")
	}
	return i.DomainControl,nil
}
