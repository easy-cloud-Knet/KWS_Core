package parsor

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

func(u *User_data_yaml) FileConfig(dirPath string) error{
	marshalledData, err := yaml.Marshal(u)
	if err!=nil{
		return err
	}
	Writer := bytes.Buffer{}
	Writer.WriteString("#cloud-config\n")
	Writer.Write(marshalledData)
	if err := os.WriteFile(filepath.Join(dirPath, "user-data"), Writer.Bytes(), 0644); err != nil {
		log.Printf("Error writing user-data file: %v", err)
		return err
	}
	return nil
}
func(u *Meta_data_yaml) FileConfig(dirPath string)error{
	marshalledData, err := yaml.Marshal(u)
	if err!=nil{
		return err
	}
	Writer := bytes.Buffer{}
	Writer.Write(marshalledData)
	if err := os.WriteFile(filepath.Join(dirPath, "meta-data"), Writer.Bytes(), 0644); err != nil {
		log.Printf("Error writing meta-data file: %v", err)
		return err
	}
	return nil
}


func(u *User_data_yaml) Parse_data(param *VM_Init_Info){

	u.PackageUpdatable=true 
	
	u.PredownProjects= []string{"qemu-guest-agent"}
	// add more packages needed



	Users_Detail:= []interface{}{"default",}

	for _,User := range param.Users{
		outputPasswd, err:= PasswdEncryption(User.PassWord)
		if err!= nil{
			fmt.Println(err)
		}
		Users_Detail=append(Users_Detail, User_specific{
			Name:User.Name,
			Passwd:outputPasswd,
			Groups: User.Groups,
			Ssh_authorized_keys: []string{"ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQC/ywMjVatnszunIy8axe43sMkzJum+Rw81UibQAID7xZouNNpDADNiQNicBW8dcuj44ScGnMZJpNmEYHgVrSCDDiC8uBC1NgzSpeURQwiSGrXZh0/sowmJaAm8cWHdvhHqFUHsIEIgSSh13iNAam2TAhajtU9MwPZreMNwNpN/qHqKHpq4FCXKn441gs7mE/VcPOj8pau6jM/9Bb8Wg9kmjhF3y1vN1YgKIXLdm0CW1x11axUKvKY7v1D7BaVL618×Md+e4zsLOCObHYw9KEsn7asOKcfUwLXScjWXNVUexv06+voltUdSA976NGHZIGZqEzvMttH+6TQVNSa78kIUls71N1A9v4yiqx"},
			SuGroup: "ALL=(ALL) NOPASSWD:ALL",
			Shell: "/bin/bash",
			Lock_passwd:false,
		}) 
	}
	var File_Appendor []User_write_file
	for index, IP := range param.IPs{
		ipCon:= strings.Split(IP, ".")
			ipAddress:=strings.Join([]string{ipCon[0],ipCon[1],ipCon[2], ipCon[3]},".")
			Gateway:= strings.Join([]string{ipCon[0],ipCon[1],ipCon[2],"1"},".")
			File_Appendor = append(File_Appendor, User_write_file{
			Path:fmt.Sprintf("/etc/systemd/network/10-enp%ds3.network",index),
			Permissions:"0644",
			Content:fmt.Sprintf("[Match]\nName=enp%ds3\n[Network]\nAddress=%s/24\nGateway=%s\nDNS=%s\nDHCP=no", index,ipAddress,Gateway,"8.8.8.8"),
		})
	}

	u.Users=Users_Detail
	u.Write_files= File_Appendor
	u.Runcmd= append(u.Runcmd, "systemctl enable systemd-networkd")
	u.Runcmd= append(u.Runcmd, "systemctl start systemd-networkd")
	u.Runcmd= append(u.Runcmd, "sudo netplan apply")
	u.Runcmd= append(u.Runcmd, "sudo systemctl start qemu-guest-agent")
	u.Runcmd= append(u.Runcmd, "sudo systemctl enable qemu-guest-agent")
	
}	
func (m* Meta_data_yaml) Parse_data(param *VM_Init_Info){
	m.Instance_ID=param.UUID
	m.Local_Host_Id=param.DomName
}


func PasswdEncryption (passwd string) (string ,error) {
		cmd := exec.Command("mkpasswd", "--method=SHA-512", "--stdin")
		
		var stdin bytes.Buffer
		stdin.WriteString(passwd)
		cmd.Stdin = &stdin
		var stdout bytes.Buffer
		cmd.Stdout = &stdout
		
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		
		err := cmd.Run()
		if err != nil {
			return "", fmt.Errorf("failed to hash password: %v, stderr: %s", err, stderr.String())
		}
		
		return strings.TrimSuffix(stdout.String(),"\n"), nil
}
