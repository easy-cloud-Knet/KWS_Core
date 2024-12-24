package conn

import (
	"encoding/json"
	_ "encoding/xml"
	"fmt"
	"net/http"
	"os/exec"

	"github.com/easy-cloud-Knet/KWS_Core.git/api/parsor"
)


func (i *InstHandler)CreateVM(w http.ResponseWriter, r * http.Request){
	
	var param parsor.VM_Init_Info
	if err:= json.NewDecoder(r.Body).Decode(&param);err!=nil{
		fmt.Printf("error",err)
	}
	
	parsor.XML_Parsor(&param)

	cmd:=exec.Command("bash", "./build/augoGen.sh", param.UUID, param.DomName, param.IP)

	output,err:=cmd.CombinedOutput()
	if err!=nil{
		fmt.Println(err)
	}

	fmt.Println(output)
	// domain,err:= i.CreateDomainWithXML()
	// if err!=nil{
	// 	fmt.Printf("error", err)
	// }
	//refer client's request,create vm with diffrent options


	// fmt.Fprintf(w, "Domain created: %v", domain)

}
