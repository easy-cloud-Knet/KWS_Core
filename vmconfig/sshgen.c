#include <stdio.h>
#include <stdlib.h>
#include <libvirt/libvirt.h>
#include <libvirt/virterror.h>

// Define user info structure
typedef struct {
    char username[50];
    char hostname[100];
    char ssh_pub_key[100];
    char ssh_priv_key[100];
} user_info;

void create_directory(user_info *user){

    snprintf(buffer, sizeof(buffer), "mkdir -p /var/lib/libvirt/images/%s", user->username);
    int dir_status = system(buffer);
    if (dir_status == -1) {
        printf("Failed to create directory for user: %s\n", user->username);
        return;
    }

}


// Create and write meta_data file
void meta_data_file(user_info *user){
    char meta_filename[100];

    snprintf(meta_filename, sizeof(meta_filename), "/var/lib/libvirt/images/%s/meta_data", user->username);
    FILE* meta_file = fopen(meta_filename, "a");
    if (meta_file == NULL) {
        printf("Failed to open meta_data file: %s\n", meta_filename);
        return;
    }
    fprintf(meta_file, "instance-id: %s\n", user->username);
    fprintf(meta_file, "local-hostname: kws\n");
    fclose(meta_file);

}

// SSH key generation
void ssh_key_gen(user_info *user){
    char buffer[256];

    snprintf(buffer, sizeof(buffer), "ssh-keygen -t rsa -b 4096 -f /var/lib/libvirt/images/%s/ssh -N ''", user->username);
    int keygen_status = system(buffer);
    if (keygen_status == -1) {
        printf("Failed to generate SSH keys for user: %s\n", user->username);
        return;
    }

    // Transfer the public key to the remote server
    snprintf(buffer, sizeof(buffer), "scp /var/lib/libvirt/images/%s/ssh.pub %s@%s:/home/%s/.ssh/authorized_keys", 
             user->username, user->username, user->hostname, user->username);
    int scp_status = system(buffer);
    
    if (scp_status == -1) {
        printf("SSH key transfer failed\n");
    } else {
        printf("SSH keys configured successfully!\n");
    }


}

// Create and write user_data file
void user_data_file(user_info *user){
    char user_filename[100];
    char pubkey_filename[150];
    char public_key[4096];


    snprintf(user_filename, sizeof(user_filename), "/var/lib/libvirt/images/%s/user_data", user->username);
    FILE* user_file = fopen(user_filename, "a");
    if (user_file == NULL) {
        printf("Failed to open user_data file: %s\n", user_filename);
        return;
    }

    // Read the public key from the generated ssh.pub file
    snprintf(pubkey_filename, sizeof(pubkey_filename), "/var/lib/libvirt/images/%s/ssh.pub", user->username);
    FILE* pubkey_file = fopen(pubkey_filename, "r");
    if (pubkey_file == NULL) {
        printf("Failed to open public key file: %s\n", pubkey_filename);
        fclose(user_file);
        return;
    }

    if (fgets(public_key, sizeof(public_key), pubkey_file) != NULL) {
        // Write the public key to the user_data file
        fprintf(user_file, "users:\n");
        fprintf(user_file, "  - name: %s\n", user->username);
        fprintf(user_file, "    ssh_authorized_keys:\n");
        fprintf(user_file, "      - %s", public_key);
    } else {
        printf("Error occurred while reading the public key.\n");
    }


    fprintf(user_file,"    sudo: [\"ALL=(ALL) NOPASSWD:ALL\"]\n");
    fprintf(user_file,"    groups: sudo\n");
    fprintf(user_file,"    shell: /bin/bash");

    fclose(pubkey_file);
    fclose(user_file);

}

void make_img(user_info *user){
    system("qemu-img create -b /var/lib/libvirt/images/baseimg/ubuntu-cloud-24.04.img -f qcow2 -F qcow2 /var/lib/libvirt/images/%s/debian-%s-qcow2.qcow2 10G",user->username,user->username);
}

void make_iso(user_info *user){
    system("genisoimage --output cidata.iso -V cidata -r -J /var/lib/libvirt/images/%s/user_data /var/lib/libvirt/images/%s/meta_data",user->username,user->username);
}




int main() {
    user_info user = {
        .username = "test-vm",
        .hostname = "192.168.42.134",
        .ssh_pub_key = "~/.ssh/test-vm.pub",
        .ssh_priv_key = "~/.ssh/test-vm"
    };

    configure_ssh(&user);  // Configure SSH
    create_directory(&user);
    make_img(&user);
    meta_data_file(&user);
    ssh_key_gen(&user);
    user_data_file(&user);
    
    return 0;
}
