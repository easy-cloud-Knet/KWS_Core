package conn

import (
	"encoding/json"
	_ "encoding/xml"
	"fmt"
	"net/http"
	"os"

	"github.com/easy-cloud-Knet/KWS_Core.git/api/parsor"
	"gopkg.in/yaml.v3"
)


func (i *InstHandler)CreateVM(w http.ResponseWriter, r * http.Request){
	
	var param parsor.VM_Init_Info
	if err:= json.NewDecoder(r.Body).Decode(&param);err!=nil{
		fmt.Printf("error",err)
	}
	
	// parsedXML:= parsor.XML_Parsor(&param)
	//need to replace with go
	var parsed_User_Yaml parsor.User_data_yaml
	var parsed_Meta_Yaml parsor.Meta_data_yaml
	
	parsed_User_Yaml.Parse_data(&param)
	parsed_Meta_Yaml.Parse_data(&param)

	data, err := yaml.Marshal(parsed_User_Yaml)
	if err!=nil{
		fmt.Println("error while unmarshaling struct")
	}
	fmt.Println(data)
	
	dirPath:= fmt.Sprintf("/var/lib/kws/%s", param.UUID)
	os.MkdirAll(dirPath, 0755)

	// shellPath:="/home/kws/kwsWorker/build/autoGen.sh"
	// fmt.Println(shellPath, param.UUID, param.DomName, param.IPs)
	// cmd:=exec.Command("bash",shellPath, param.UUID, param.DomName, param.IPs[0])
	// err:=cmd.Run()
	// if err!=nil{
	// 	// fmt.Println(output)
	// 	fmt.Println(err)
	// }

	// domain,err:= i.CreateDomainWithXML(parsedXML)
	//  if err!=nil{
	//  	fmt.Printf("error", err)
	//  }
	//refer client's request,create vm with diffrent options


	// fmt.Fprintf(w, "Domain created: %v", domain)

}
