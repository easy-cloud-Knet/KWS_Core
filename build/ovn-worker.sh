#!/bin/bash

# --- 0. 기존 프로세스 및 소켓/락 파일 정리 (매우 중요!) ---
echo "--- Cleaning up existing OVS/OVN processes and files ---"
sudo pkill ovsdb-server
sudo pkill ovs-vswitchd
sudo pkill ovn-northd # 워커 노드에서는 실행되지 않아야 하지만, 안전을 위해 포함
sudo pkill ovn-controller
sudo rm -f /usr/local/etc/ovn/*.lock
sudo rm -f /usr/local/var/run/ovn/*.sock
sudo rm -f /usr/local/var/run/openvswitch/*.sock
echo "--- Cleanup complete ---"
echo ""

# --- 1. 필수 디렉토리 생성 ---
echo "--- Creating necessary directories ---"
sudo mkdir -p /usr/local/etc/ovn
sudo mkdir -p /usr/local/var/run/ovn
sudo mkdir -p /usr/local/var/log/ovn
sudo mkdir -p /usr/local/var/log/openvswitch/
echo "--- Directories created ---"
echo ""

# --- 2. OVS 코어 서비스 시작 ---
echo "--- Starting OVS core services ---"
sudo /usr/local/share/openvswitch/scripts/ovs-ctl start
sleep 5 # OVS가 완전히 시작될 시간을 줍니다.
echo "--- OVS core services started ---"
echo ""

# --- 3. 네트워크 설정 변수 ---
# 워커 노드의 IP와 DNS를 인자로 받습니다.
IP_ADDRESS=$1
DNS=$2
FILE_PATH="/etc/systemd/network/20-br-ext.network"

# --- 4. systemd-networkd 설정 파일 생성 ---
echo "--- Creating systemd-networkd config for br-ext ---"
cat <<EOF | sudo tee "$FILE_PATH" > /dev/null
[Match]
Name=br-ext

[Network]
Address=$IP_ADDRESS/24
DNS=$DNS
Gateway=$DNS
EOF
echo "--- systemd-networkd config created ---"
echo ""

# --- 5. OVS 브릿지 생성 및 포트 추가 ---
# br-ext가 이미 존재할 경우 오류가 나므로, 먼저 삭제 후 추가하는 것이 안전합니다.
echo "--- Configuring OVS bridges and ports ---"
sudo ovs-vsctl  del-br br-ext # 이미 존재하면 삭제
sudo ovs-vsctl add-br br-ext
sudo ovs-vsctl add-port br-ext ens3 # 물리 NIC 이름이 ens3인지 다시 확인하세요.
echo "--- OVS bridges and ports configured ---"
echo ""

# --- 6. systemd-networkd 서비스 재시작 ---
echo "--- Restarting systemd-networkd ---"
sudo systemctl restart systemd-networkd
sleep 5 # 네트워크 설정이 적용될 시간을 줍니다.
echo "--- systemd-networkd restarted ---"
echo ""

# --- 7. OVN external-ids 설정 ---
# **주의: 스마트 따옴표 “ ” 대신 일반 따옴표 " " 사용!**
# 워커 노드이므로 ovn-nb 설정은 필요 없습니다.
echo "--- Setting OVN external-ids ---"
sudo ovs-vsctl set open_vswitch . external-ids:ovn-remote="tcp:10.5.15.39:6642" # SB DB는 컨트롤러 노드(10.5.15.39)
sudo ovs-vsctl set open_vswitch . external-ids:ovn-encap-type=geneve
sudo ovs-vsctl set open_vswitch . external-ids:ovn-encap-ip="$IP_ADDRESS" # 이 노드의 br-ext IP
sudo ovs-vsctl set open_vswitch . external-ids:ovn-bridge-mappings=UPLINK:br-ext
sudo ovs-vsctl set open_vswitch . external-ids:system-id="$(hostname)" # system-id 명시적 설정
echo "--- OVN external-ids set ---"
echo ""

# --- 8. ovn-controller 데몬 시작 ---
echo "--- Starting ovn-controller daemon ---"
sudo ovn-ctl start_controller
sleep 5 # 데몬이 시작될 시간을 줍니다.
echo "--- ovn-controller started ---"
echo ""

echo "--- Worker Node Setup Complete ---"
echo "Please check status with: sudo ovs-vsctl show, sudo ovn-sbctl show (from control node), tail -f /usr/local/var/log/ovn/ovn-controller.log"
