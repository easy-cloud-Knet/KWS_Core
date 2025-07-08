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

sudo ovsdb-tool create /usr/local/etc/ovn/ovnnb_db.db /usr/local/share/ovn/ovn-nb.ovsschema || true
sudo ovsdb-tool create /usr/local/etc/ovn/ovnsb_db.db /usr/local/share/ovn/ovn-sb.ovsschema || true

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





sudo tee /etc/systemd/system/ovn-nb.service > /dev/null << EOF
[Unit]
Description=OVN Northbound Database
After=network-online.target openvswitch.service
Requires=openvswitch.service

[Service]
Type=forking
ExecStart=/usr/local/sbin/ovsdb-server /usr/local/etc/ovn/ovnnb_db.db \
	    --remote=punix:/usr/local/var/run/ovn/ovnnb_db.sock \
	        --remote=ptcp:6641:0.0.0.0 \
		    --remote=db:OVN_Northbound,NB_Global,connections \
		        --pidfile=/usr/local/var/run/ovn/ovnnb-server.pid --detach \
			    --log-file=/usr/local/var/log/ovn/ovnnb-server.log
ExecStop=/usr/local/bin/ovs-appctl -t ovsdb-server exit
PIDFile=/usr/local/var/run/ovn/ovnnb-server.pid
Restart=on-failure
RestartSec=5s

[Install]
WantedBy=multi-user.target
EOF

sudo tee /etc/systemd/system/ovn-sb.service > /dev/null << EOF
[Unit]
Description=OVN Southbound Database
After=network-online.target openvswitch.service ovn-nb.service
Requires=openvswitch.service ovn-nb.service

[Service]
Type=forking
ExecStart=/usr/local/sbin/ovsdb-server /usr/local/etc/ovn/ovnsb_db.db \
	    --remote=punix:/usr/local/var/run/ovn/ovnsb_db.sock \
	        --remote=ptcp:6642:0.0.0.0 \
		    --remote=db:OVN_Southbound,SB_Global,connections \
		        --pidfile=/usr/local/var/run/ovn/ovnsb-server.pid --detach \
			    --log-file=/usr/local/var/log/ovn/ovnsb-server.log
ExecStop=/usr/local/bin/ovs-appctl -t ovsdb-server exit
PIDFile=/usr/local/var/run/ovn/ovnsb-server.pid
Restart=on-failure
RestartSec=5s

[Install]
WantedBy=multi-user.target
EOF

sudo tee /etc/systemd/system/ovn-northd.service > /dev/null << EOF
[Unit]
Description=OVN Northd Daemon
After=ovn-nb.service ovn-sb.service
Requires=ovn-nb.service ovn-sb.service

[Service]
Type=forking
ExecStart=/usr/local/bin/ovn-northd --pidfile=/usr/local/var/run/ovn/ovn-northd.pid --detach --log-file=/usr/local/var/log/ovn/ovn-northd.log
ExecStop=/usr/local/bin/ovn-appctl -t ovn-northd exit
PIDFile=/usr/local/var/run/ovn/ovn-northd.pid
Restart=on-failure
RestartSec=5s

[Install]
WantedBy=multi-user.target
EOF

sudo tee /etc/systemd/system/ovn-controller.service > /dev/null << EOF
[Unit]
Description=OVN Controller
After=openvswitch.service ovn-sb.service
Requires=openvswitch.service ovn-sb.service

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
sleep 10 



sudo ovs-vsctl --may-exist del-br br-ext
sudo ovs-vsctl add-br br-ext
sudo ovs-vsctl add-port br-ext ens3


sudo systemctl restart systemd-networkd

sudo ovs-vsctl set open_vswitch . external-ids:ovn-remote="tcp:10.5.15.39:6642"
sudo ovs-vsctl set open_vswitch . external-ids:ovn-encap-type=geneve
sudo ovs-vsctl set open_vswitch . external-ids:ovn-nb="tcp:10.5.15.39:6641"
sudo ovs-vsctl set open_vswitch . external-ids:ovn-encap-ip="$IP_ADDRESS"
sudo ovs-vsctl set open_vswitch . external-ids:ovn-bridge-mappings=UPLINK:br-ext
sudo ovs-vsctl set open_vswitch . external-ids:system-id="$(hostname)"

sudo systemctl enable --now ovn-nb.service
sleep 1
sudo systemctl enable --now ovn-sb.service
sleep 1
sudo ovn-nbctl --no-wait init
sudo ovn-sbctl init

sudo systemctl enable --now ovn-northd.service
sleep 1
sudo systemctl enable --now ovn-controller.service


echo "Control Node Setup Complete."
echo "Check status with: sudo systemctl status openvswitch ovn-nb ovn-sb ovn-northd ovn-controller"
echo "Also check: sudo ovs-vsctl show, sudo ovn-sbctl show, sudo lsof -i :6641, sudo lsof -i :6642"
