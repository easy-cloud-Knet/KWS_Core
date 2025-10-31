package userconfig

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/easy-cloud-Knet/KWS_Core/vm/parsor"
	"gopkg.in/yaml.v3"
)

func (u *User_data_yaml) WriteFile(dirPath string) error {
	marshalledData, err := yaml.Marshal(u)
	if err != nil {
		return err
	}
	Writer := bytes.Buffer{}
	Writer.WriteString("#cloud-config\n")
	Writer.Write(marshalledData)
	if err := os.WriteFile(filepath.Join(dirPath, "user-data"), Writer.Bytes(), 0644); err != nil {
		return fmt.Errorf("error Writing user-data file %w", err)
	}
	return nil
}

func (u *Meta_data_yaml) WriteFile(dirPath string) error {
	marshalledData, err := yaml.Marshal(u)
	if err != nil {
		return err
	}
	Writer := bytes.Buffer{}
	Writer.Write(marshalledData)
	if err := os.WriteFile(filepath.Join(dirPath, "meta-data"), Writer.Bytes(), 0644); err != nil {
		return fmt.Errorf("error Writing Meta-data file %w", err)
	}
	return nil
}

func (u *User_data_yaml) ParseData(param *parsor.VM_Init_Info) error {

	u.PackageUpdatable = true

	u.PredownProjects = []string{"qemu-guest-agent"}
	// add more packages needed
	Users_Detail := []interface{}{"default"}

	for i, User := range param.Users {
		outputPasswd, err := PasswdEncryption(User.PassWord)
		if err != nil {
			return fmt.Errorf("error Encrypting password %w of User index %d", err, i)
		}

		Users_Detail = append(Users_Detail, User_specific{
			Name:                User.Name,
			Passwd:              outputPasswd,
			Groups:              User.Groups,
			Ssh_authorized_keys: []string{User.Ssh_authorized_keys[0]},
			SuGroup:             "ALL=(ALL) NOPASSWD:ALL",
			Shell:               "/bin/bash",
			Lock_passwd:         false,
		})
	}
	File_Appendor := u.configNetworkIP(param.NetConf.Ips)
	appending,err:= u.fetchos()
	if err!=nil{
		log.Printf("error from appending file, %v, ignoring", err)
	}
	File_Appendor = append(File_Appendor,appending... )
	u.Users = Users_Detail
	u.Write_files = File_Appendor

	u.configNetworkCommand()
	u.configQEMU()
	u.configSsh()

	return nil
}

func (u *User_data_yaml) fetchos() ([]User_write_file,error){    
    var fetchFile []User_write_file

	data, err := os.ReadFile("/var/lib/kws/baseimg/fetchos")
    if err != nil {
        return  fetchFile,fmt.Errorf("fetchos parsing error, ingnorable but needs report%v",err)
    }

	fetchFile= append(fetchFile, User_write_file{
		Path: "/etc/profile.d/99-my-motd.sh",
		Permissions: "0644",
		Content: string(data),
	})    
    return fetchFile,nil
}

func (u *User_data_yaml) configNetworkIP(ips []string) []User_write_file {
	var File_Appendor []User_write_file

	for index, IP := range ips {
		ipCon := strings.Split(IP, ".")
		ipAddress := strings.Join([]string{ipCon[0], ipCon[1], ipCon[2], ipCon[3]}, ".")
		Gateway := strings.Join([]string{ipCon[0], ipCon[1], ipCon[2], "1"}, ".")
		File_Appendor = append(File_Appendor, User_write_file{
			Path:        fmt.Sprintf("/etc/systemd/network/10-enp%ds3.network", index),
			Permissions: "0644",
			Content:     fmt.Sprintf("[Match]\nName=enp%ds3\n[Network]\nAddress=%s/24\nGateway=%s\nDNS=%s\nDHCP=no", index, ipAddress, Gateway, "8.8.8.8"),
		}) //
	}

	return File_Appendor
}




func (u *User_data_yaml) configSsh(){
	u.Runcmd = append(u.Runcmd, "sed -i 's/PasswordAuthentication no/PasswordAuthentication yes/g' /etc/ssh/sshd_config.d/60-cloudimg-settings.conf")
	u.Runcmd=append(u.Runcmd, "systemctl restart ssh")
	u.Runcmd=append(u.Runcmd, "systemctl enable ssh")

}



func (u *User_data_yaml) configNetworkCommand() {
	u.Runcmd = append(u.Runcmd, "systemctl enable systemd-networkd")
	u.Runcmd = append(u.Runcmd, "systemctl start systemd-networkd")
	u.Runcmd = append(u.Runcmd, "sudo systemctl disable systemd-networkd-wait-online.service")
	u.Runcmd = append(u.Runcmd, "sudo systemctl mask systemd-networkd-wait-online.service")
	
	u.Runcmd = append(u.Runcmd, "sudo netplan apply")
}

func (u *User_data_yaml) configQEMU() {
	u.Runcmd = append(u.Runcmd, "sudo systemctl start qemu-guest-agent")
	u.Runcmd = append(u.Runcmd, "sudo systemctl enable qemu-guest-agent")
}

func (m *Meta_data_yaml) ParseData(param *parsor.VM_Init_Info) error {
	m.Instance_ID = param.UUID
	m.Local_Host_Id = param.DomName
	return nil
}

func PasswdEncryption(passwd string) (string, error) {
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

	return strings.TrimSuffix(stdout.String(), "\n"), nil
}
