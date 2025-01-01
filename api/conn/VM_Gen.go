package conn

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/easy-cloud-Knet/KWS_Core.git/api/parsor"
	"gopkg.in/yaml.v3"
)

func (i *InstHandler) CreateVM(w http.ResponseWriter, r *http.Request) {
	var param parsor.VM_Init_Info


	if err := json.NewDecoder(r.Body).Decode(&param); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		log.Printf("Error decoding JSON: %v", err)
		return
	}

	dirPath := fmt.Sprintf("/var/lib/kws/%s", param.UUID)
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		http.Error(w, "Failed to create directory", http.StatusInternalServerError)
		log.Printf("Error creating directory %s: %v", dirPath, err)
		return
	}

	// Parse user-data and meta-data
	var parsedUserYaml parsor.User_data_yaml
	var parsedMetaYaml parsor.Meta_data_yaml

	parsedUserYaml.Parse_data(&param)
	parsedMetaYaml.Parse_data(&param)

	// Marshal YAML data
	marshalledUserData, err := yaml.Marshal(parsedUserYaml)
	if err != nil {
		http.Error(w, "Failed to marshal user data", http.StatusInternalServerError)
		log.Printf("Error marshaling user data: %v", err)
		return
	}

	marshalledMetaData, err := yaml.Marshal(parsedMetaYaml)
	if err != nil {
		http.Error(w, "Failed to marshal meta data", http.StatusInternalServerError)
		log.Printf("Error marshaling meta data: %v", err)
		return
	}

	// Write user-data file
	userConfig := bytes.Buffer{}
	userConfig.WriteString("#cloud-config\n")
	userConfig.Write(marshalledUserData)
	if err := os.WriteFile(filepath.Join(dirPath, "user-data"), userConfig.Bytes(), 0644); err != nil {
		http.Error(w, "Failed to write user-data file", http.StatusInternalServerError)
		log.Printf("Error writing user-data file: %v", err)
		return
	}

	// Write meta-data file
	metaConfig := bytes.Buffer{}
	metaConfig.Write(marshalledMetaData)
	if err := os.WriteFile(filepath.Join(dirPath, "meta-data"), metaConfig.Bytes(), 0644); err != nil {
		http.Error(w, "Failed to write meta-data file", http.StatusInternalServerError)
		log.Printf("Error writing meta-data file: %v", err)
		return
	}

	// Execute qemu-img create command
	baseImage := fmt.Sprintf("/var/lib/kws/baseimg/%s", param.OS )
	targetImage := filepath.Join(dirPath, fmt.Sprintf("%s.qcow2", param.UUID))
	qemuImgCmd := exec.Command("qemu-img", "create",
		"-b", baseImage,
		"-f", "qcow2",
		"-F", "qcow2",
		targetImage, "10G",
	)

	qemuImgCmd.Stdout = os.Stdout
	qemuImgCmd.Stderr = os.Stderr

	log.Println("Creating disk image...")
	if err := qemuImgCmd.Run(); err != nil {
		http.Error(w, "Failed to create disk image", http.StatusInternalServerError)
		log.Printf("qemu-img command failed: %v", err)
		return
	}

	// Execute genisoimage command
	isoOutput := filepath.Join(dirPath, "cidata.iso")
	userDataPath := filepath.Join(dirPath, "user-data")
	metaDataPath := filepath.Join(dirPath, "meta-data")

	genisoCmd := exec.Command("genisoimage",
		"--output", isoOutput,
		"-V", "cidata",
		"-r", "-J",
		userDataPath, metaDataPath,
	)

	genisoCmd.Stdout = os.Stdout
	genisoCmd.Stderr = os.Stderr

	log.Println("Generating ISO image...")
	if err := genisoCmd.Run(); err != nil {
		http.Error(w, "Failed to generate ISO image", http.StatusInternalServerError)
		log.Printf("genisoimage command failed: %v", err)
		return
	}

	// (선택 사항) 추가적인 VM 생성 로직을 여기에 구현
	// 예: libvirt를 사용하여 도메인 생성 등

	// 성공 응답
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "VM with UUID %s created successfully.", param.UUID)
}
