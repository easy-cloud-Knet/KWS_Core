package service

import (
	"fmt"
	"net/http"
)






func CommonErrorHelper(w http.ResponseWriter,err error, statusCode int ,message string, ){
	w.WriteHeader(statusCode)
	fmt.Fprintf(w,"{error: %v, message:%s}",err,message )

}