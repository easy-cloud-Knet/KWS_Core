#!/bin/bash

sudo mkdir -p /usr/local/etc/ovn
sudo mkdir -p /usr/local/var/run/ovn
sudo mkdir -p /usr/local/var/log/ovn
sudo mkdir -p /usr/local/var/log/openvswitch/

sudo ovsdb-tool create /usr/local/etc/ovn/ovnnb_db.db /usr/local/share/ovn/ovn-nb.ovsschema
sudo ovsdb-tool create /usr/local/etc/ovn/ovnsb_db.db /usr/local/share/ovn/ovn-sb.ovsschema


sudo ovsdb-server /usr/local/etc/ovn/ovnnb_db.db \
	--remote=punix:/usr/local/var/run/ovn/ovnnb_db.sock \
	--remote=ptcp:6641:0.0.0.0 \
	--remote=db:OVN_Northbound,NB_Global,connections \
	--pidfile=/usr/local/var/run/ovn/ovnnb-server.pid --detach \
	--log-file=/usr/local/var/log/ovn/ovnnb-server.log


sudo ovsdb-server /usr/local/etc/ovn/ovnsb_db.db \
	--remote=punix:/usr/local/var/run/ovn/ovnsb_db.sock \
	--remote=ptcp:6642:0.0.0.0 \
	--remote=db:OVN_Southbound,SB_Global,connections \
	--pidfile=/usr/local/var/run/ovn/ovnsb-server.pid --detach \
	--log-file=/usr/local/var/log/ovn/ovnsb-server.log



sudo ovn-nbctl --no-wait init
sudo ovn-sbctl init


sudo ovn-northd --pidfile --detach --log-file=/usr/local/var/log/ovn/ovn-northd.log


sudo ovs-ctl start 

IP_ADDRESS=$1
DNS=$2
FILE_PATH="/etc/systemd/network/20-br-ext.network"


cat <<EOF | sudo tee "$FILE_PATH" > /dev/null
[Match]
Name=br-ext

[Network]
Address=$IP_ADDRESS/24
DNS=$DNS
Gateway=$DNS
EOF


sudo ovs-vsctl add-br br-ext
sudo systemctl restart systemd-networkd
sudo ovs-vsctl add-port br-ext ens3

sudo ovs-vsctl set open_vswitch . external-ids:ovn-remote="tcp:10.5.15.39:6642"
sudo ovs-vsctl set open_vswitch . external-ids:ovn-encap-type=geneve
sudo ovs-vsctl set open_vswitch . external-ids:ovn-nb="tcp:10.5.15.39:6641" 
sudo ovs-vsctl set open_vswitch . external-ids:ovn-encap-ip=10.5.15.39 
sudo ovs-vsctl set open_vswitch . external-ids:ovn-bridge-mappings=UPLINK:br-ext


sudo ovn-ctl start_controller 
