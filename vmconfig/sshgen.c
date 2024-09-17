#include <stdio.h>
#include <stdlib.h>
#include <libvirt/libvirt.h>
#include <libvirt/virterror.h>

// 사용자 정보 구조체 정의
typedef struct {
    char username[50];
    char hostname[100];
    char ssh_pub_key[100];
    char ssh_priv_key[100];
} user_info;


void setup_network(){
     system("sudo virsh net-start default");
     system("sudo virsh net-autostart default");

}

void create_vm_user(const char *username, const char *hostname) {
    char buffer[256];

    // SSH를 통해 VM에 접속해 계정 생성 명령 실행
    snprintf(buffer, sizeof(buffer), "ssh root@%s 'sudo adduser --disabled-password --gecos \"\" %s'", hostname, username);
    int status = system(buffer);

    if (status == -1) {
        printf("Failed to create user %s on the VM\n", username);
    } else {
        printf("User %s created successfully on the VM\n", username);
    }
}




// create inner vm
void create_vm(const char *xml_path) {
    virConnectPtr conn;//connection libvirt to qemu
    virDomainPtr dom;//domain of inner vm

    conn = virConnectOpen("qemu:///system");

    if (conn == NULL) {
        printf("Failed to open connection to qemu:///system\n");
        exit(1);
    }

    FILE *xml_file = fopen(xml_path, "r");//read .xml
    if (!xml_file) {
        printf("Could not open the VM XML file.\n");
        virConnectClose(conn);
        exit(1);
    }

    //check xml_file size
    fseek(xml_file, 0, SEEK_END);
    long xml_size = ftell(xml_file);
    rewind(xml_file);



    //read xml_file
    char *xml_content = malloc(xml_size + 1);
    fread(xml_content, 1, xml_size, xml_file);
    xml_content[xml_size] = '\0';
    fclose(xml_file);



    //create inner vm
    dom = virDomainCreateXML(conn, xml_content, 0);


    //delete xml_content
    free(xml_content);

    if (dom == NULL) {
        printf("Failed to create the VM from XML definition.\n");
    } else {
        printf("VM created successfully!\n");
        virDomainFree(dom);
    }

    virConnectClose(conn);
}


// 2. SSH 설정 함수
void configure_ssh(user_info *user) {
    char buffer[256];
    
    // SSH 키 생성
    snprintf(buffer, sizeof(buffer), "ssh-keygen -t rsa -b 4096 -f ~/.ssh/%s", user->username);
    system(buffer);
    
    // SSH 공개키 전송
    snprintf(buffer, sizeof(buffer), "ssh-copy-id -i ~/.ssh/%s.pub %s@%s", user->username, user->username, user->hostname);
    int status = system(buffer);
    
    if (status == -1) {
        printf("SSH configuration failed\n");
    } else {
        printf("SSH keys configured successfully!\n");
    }
}

// 3. 사용자 매핑 XML 생성 함수
// void create_user_mapping(user_info *user) {
//     FILE *file = fopen("user-mapping.xml", "a");

//     if (file == NULL) {
//         printf("Failed to open user-mapping.xml\n");
//         exit(1);
//     }

//     fprintf(file, "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n");
//     fprintf(file, "<user-mapping>\n");
//     fprintf(file, "    <authorize username=\"%s\">\n", user->username);
//     fprintf(file, "        <connection name=\"%s Server\">\n", user->hostname);
//     fprintf(file, "            <protocol>ssh</protocol>\n");
//     fprintf(file, "            <param name=\"hostname\">%s</param>\n", user->hostname);
//     fprintf(file, "        </connection>\n");
//     fprintf(file, "    </authorize>\n");
//     fprintf(file, "</user-mapping>\n");

//     fclose(file);
//     printf("user-mapping.xml created.\n");
// }

int main() {
    user_info user = {
        .username = "test-vm",
        .hostname = "192.168.42.134",
        .ssh_pub_key = "~/.ssh/test-vm.pub",
        .ssh_priv_key = "~/.ssh/test-vm"
    };

    // VM 생성 및 SSH 설정 자동화
    create_vm("/home/debian/KWS_Core/xmlFile/vm1.xml"); 
    create_vm_user(user.username, user.hostname);
    configure_ssh(&user);             // SSH 설정
    create_user_mapping(&user);       // XML 매핑 생성

    setup_network();

    return 0;
}

