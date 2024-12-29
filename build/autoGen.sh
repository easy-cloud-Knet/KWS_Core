#!/bin/bash

METADATA=$1
USERNAME=$2
IP=$3

# Validate input
if [[ -z "$METADATA" || -z "$USERNAME" || -z "$IP" ]]; then
  echo "Usage: $0 <METADATA> <USERNAME> <IP>"
  exit 1
fi

# Create directory
if [[ -d "/var/lib/kws/${METADATA}" ]]; then
  echo "Directory /var/lib/kws/${METADATA} already exists. Exiting."
  exit 1
fi
mkdir -p /var/lib/kws/${METADATA}

# Generate user-data
PASSWORD_HASH=$(openssl passwd -6 "password")
cat <<EOF > /var/lib/kws/${METADATA}/user-data
#cloud-config
users:
  - default
  - name: ${USERNAME}
    passwd: ${PASSWORD_HASH}
    lock_passwd: false
    ssh_authorized_keys:
      - ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQC/ywMjVatnszunIy8axe43sMkzJum+Rw81UibQAID7xZouNNpDADNiQNicBW8dcuj44ScGnMZJpNmEYHgVrSCDDiC8uBC1NgzSpeURQwiSGrXZh0/sowmJaAm8cWHdvhHqFUHsIEIgSSh13iNAam2TAhajtU9MwPZreMNwNpN/qHqKHpq4FCXKn441gs7mE/VcPOj8pau6jM/9Bb8Wg9kmjhF3y1vN1YgKIXLdm0CW1x11axUKvKY7v1D7BaVL618×Md+e4zsLOCObHYw9KEsn7asOKcfUwLXScjWXNVUexv06+voltUdSA976NGHZIGZqEzvMttH+6TQVNSa78kIUls71N1A9v4yiqx
    sudo: ALL=(ALL) NOPASSWD:ALL
    groups: sudo
    shell: /bin/bash

write_files:
  - path: /etc/systemd/network/10-enp0s3.network
    permissions: "0644"
    content: |
      [Match]
      Name=enp0s3

      [Network]
      Address=${IP}/24
      Gateway=10.5.12.1
      DNS=10.5.12.1
      DHCP=no

runcmd:
  - systemctl enable systemd-networkd
  - systemctl start systemd-networkd
EOF

# Generate meta-data
cat <<EOF > /var/lib/kws/${METADATA}/meta-data
instance-id: ${METADATA}
local-hostname: ${USERNAME}
EOF

# Create disk image
qemu-img create -b /var/lib/kws/baseimg/ubuntu-cloud-24.04.img -f qcow2 -F qcow2 "/var/lib/kws/${METADATA}/${METADATA}.qcow2" 10G

# Generate ISO
genisoimage --output /var/lib/kws/${METADATA}/cidata.iso -V cidata -r -J /var/lib/kws/${METADATA}/user-data /var/lib/kws/${METADATA}/meta-data
