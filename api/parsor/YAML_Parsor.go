package parsor

import (
	"fmt"
	"os/exec"
	"strings"
)



func(u *User_data_yaml) Parse_data(param *VM_Init_Info){
	var Users_Detail []User_specific

	Users_Detail= append(Users_Detail, User_specific{Name:"default"})
	for _,User := range param.Users{
		output, err:= exec.Command("openssl", "passwd", "-6", User.PassWord).Output()
		if err!= nil{
			fmt.Println(err)
		}
		Users_Detail=append(Users_Detail, User_specific{
			Name:User.Name,
			Passwd:string(output),
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


	
}	
func (m* Meta_data_yaml) Parse_data(param *VM_Init_Info){
	m.Instance_ID=param.UUID
	m.Local_Host_Id=param.DomName
}
