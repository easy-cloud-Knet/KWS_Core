package conn

import (
	"bytes"
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
	
	parsed:= parsor.XML_Parsor(&param)
	shellPath:="/home/kws/kwsWorker/build/autoGen.sh"
	fmt.Println(shellPath, param.UUID, param.DomName, param.IPs)
	cmd:=exec.Command("bash",shellPath, param.UUID, param.DomName, param.IPs[0])
	var output bytes.Buffer
	cmd.Stdout= &output
	cmd.Stderr=&output
	fmt.Printf("%s", output)
	// output,
	err:=cmd.Run()
	if err!=nil{
		// fmt.Println(output)
		fmt.Println(err)
	}

	fmt.Println(output)
	domain,err:= i.CreateDomainWithXML(parsed)
	 if err!=nil{
	 	fmt.Printf("error", err)
	 }
	//refer client's request,create vm with diffrent options


	fmt.Fprintf(w, "Domain created: %v", domain)

}
