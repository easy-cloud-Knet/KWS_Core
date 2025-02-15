package timeCal

import (
	"fmt"
	"net/http"
	"time"
)



func TimeLogging (next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		start:= time.Now()

		next.ServeHTTP(w,r)

		elapsed := time.Since(start)
		fmt.Println("opeartion", elapsed)
	})
}