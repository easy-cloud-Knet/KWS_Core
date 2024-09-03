package conn

import (
	"log"
	"net/http"
)


func (d *Domain)CreateVM(w http.ResponseWriter, r * http.Request){
	err:=d.Domain.Create()
	if err!=nil{
		log.Fatal(err)
	}
}
