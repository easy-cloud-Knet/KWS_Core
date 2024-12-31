package parsor

import (
	"fmt"
	"strings"
)






func(u *User_data_yaml) Parse_data(param *VM_Init_Info){
	var Users_Detail []User_specific
	
	for _,User := range param.Users{
		Users_Detail=append(Users_Detail, User_specific{
			Name:User.Name,
			Passwd:User.PassWord,
			Groups: User.Groups,
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
			Content:fmt.Sprintf(`[Match]
      Name=enp%ds3
      [Network]
      Address=%s/24
      Gateway=%s 
      DNS=%s
      DHCP=no`, index,ipAddress,Gateway,"8.8.8.8"),
		})
	}

	u.Users=Users_Detail
	u.Write_files= File_Appendor
	u.Runcmd= append(u.Runcmd, "systemctl enable systemd-networkd")
	u.Runcmd= append(u.Runcmd, "systemctl start systemd-networkd")


	
}	


func (m* Meta_data_yaml) Parse_data(param *VM_Init_Info){

}