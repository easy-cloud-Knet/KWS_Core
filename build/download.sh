sudo apt install --no-install-recommends qemu-system libvirt-clients libvirt-daemon-system libvirt-dev dnsmasq virtinst pkg-config
sudo apt install gcc


sudo wget -O /var/lib/libvirt/images/debain-12.7.0.qcow2  https://cloud.debian.org/images/cloud/bookworm/20240901-1857/debian-12-generic-amd64-20240901-1857.qcow2
sudo wget -O /var/lib/libvirt/image/ubuntu-cloud-24.04.img https://cloud-images.ubuntu.com/noble/20240822/noble-server-cloudimg-amd64.img