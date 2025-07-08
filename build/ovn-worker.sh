#!/bin/bash

sudo pkill ovsdb-server || true
sudo pkill ovs-vswitchd || true
sudo pkill ovn-northd || true
sudo pkill ovn-controller || true
sudo rm -f /usr/local/etc/ovn/*.lock || true
sudo rm -f /usr/local/var/run/ovn/*.sock || true
sudo rm -f /usr/local/var/run/openvswitch/*.sock || true

sudo mkdir -p /usr/local/etc/ovn
sudo mkdir -p /usr/local/var/run/ovn
sudo mkdir -p /usr/local/var/log/ovn
sudo mkdir -p /usr/local/var/log/openvswitch/

sudo ovs-vsctl --may-exist del-br br-ext
sudo ovs-vsctl add-br br-ext
sudo ovs-vsctl add-port br-ext ens3

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



sudo tee /etc/systemd/system/openvswitch.service > /dev/null << EOF
[Unit]
Description=Open vSwitch
After=network-online.target
Wants=network-online.target

[Service]
Type=forking
ExecStart=/usr/local/share/openvswitch/scripts/ovs-ctl start --system-id=random
ExecStop=/usr/local/share/openvswitch/scripts/ovs-ctl stop
ExecReload=/usr/local/share/openvswitch/scripts/ovs-ctl restart
PIDFile=/usr/local/var/run/openvswitch/ovs-vswitchd.pid
Restart=on-failure
RestartSec=5s

[Install]
WantedBy=multi-user.target
EOF

sudo tee /etc/systemd/system/ovn-controller.service > /dev/null << EOF
[Unit]
Description=OVN Controller
After=openvswitch.service
Requires=openvswitch.service

[Service]
Type=forking
ExecStart=/usr/local/bin/ovn-controller --pidfile=/usr/local/var/run/ovn/ovn-controller.pid --detach --log-file=/usr/local/var/log/ovn/ovn-controller.log
ExecStop=/usr/local/bin/ovn-appctl -t ovn-controller exit
PIDFile=/usr/local/var/run/ovn/ovn-controller.pid
Restart=on-failure
RestartSec=5s

[Install]
WantedBy=multi-user.target
EOF

sudo systemctl daemon-reload
sleep 1

sudo systemctl enable --now openvswitch.service
sudo systemctl enable --now ovn-controller.service


sudo ovs-vsctl set open_vswitch . external-ids:ovn-remote="tcp:10.5.15.39:6642"
sudo ovs-vsctl set open_vswitch . external-ids:ovn-encap-type=geneve
sudo ovs-vsctl set open_vswitch . external-ids:ovn-encap-ip="$IP_ADDRESS"
sudo ovs-vsctl set open_vswitch . external-ids:ovn-bridge-mappings=UPLINK:br-ext
sudo ovs-vsctl set open_vswitch . external-ids:system-id="$(hostname)"

sudo systemctl restart systemd-networkd
sleep 5

echo "Worker Node Setup Complete."
echo "Check status with: sudo systemctl status openvswitch ovn-controller"
echo "Also check from control node: sudo ovn-sbctl show"

