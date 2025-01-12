sudo apt install -y --no-install-recommends qemu-system libvirt-clients libvirt-daemon-system libvirt-dev dnsmasq virtinst pkg-config whois qemu-guest-agent
sudo apt install -y  gcc

sudo apt install -y  cloud-init genisoimage

mkdir -p /var/lib/kws/baseimg
mkdir /var/lib/kws/userConf

sudo wget -O /var/lib/kws/baseimg/debian-12.7.0.qcow2  https://cloud.debian.org/images/cloud/bookworm/20240901-1857/debian-12-generic-amd64-20240901-1857.qcow2




sudo wget -O /var/lib/kws/baseimg/ubuntu-cloud-24.04.img   https://cloud-images.ubuntu.com/noble/current/noble-server-cloudimg-amd64.img
