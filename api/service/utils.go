package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// func response(w http.ResponseWriter, err error, statusCode int, message string)error{
// 	return nil
// }
type BaseResponse struct{
	Information interface{} `json:"information,omitempty"`
	Message string `json:"message"`
	Errors *error `json:"error"`
}

func ResponseGen(message string)*BaseResponse{
	return &BaseResponse{
		Information: nil,
		Message: fmt.Sprintf("%s operation", message),
		Errors: nil,
	}
}



func HttpDecoder[T any](w http.ResponseWriter, r *http.Request,param *T) error{
	w.Header().Set("Content-Type", "application/json")
	
	body, err:= io.ReadAll(r.Body)
	if err!=nil{
		return fmt.Errorf("error occured while decoding Body")		
	}
	if err:= json.Unmarshal(body, param); err!=nil{
		return fmt.Errorf("error occured while decoding Body")		
	}
	return nil
}
//미들웨어 사용 고려
func (BR * BaseResponse)ResponseWriteErr(w http.ResponseWriter,err error, statusCode int){
	data, ERR :=json.Marshal(BR)
	if ERR!=nil{
		//이전 에러 포함 방법 고민
		http.Error(w, "failed in Marshaling output data", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(statusCode)
	BR.Message+="failed"
	w.Write(data)
}


func (BR * BaseResponse)ResponseWriteOK(w http.ResponseWriter, info interface{}){
	data, err :=json.Marshal(BR)
	if err!=nil{
		http.Error(w, "failed in Marshaling output data", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	BR.Message+="success"	
	w.Write(data)
}
// type generic can be implemented in further needs
