#!/bin/bash

# --- 0. 기존 프로세스 및 소켓/락 파일 정리 (매우 중요!) ---
# 스크립트 실행 전에 항상 이 부분을 추가하여 깨끗한 상태에서 시작하도록 합니다.
# 특히 디버깅 중에는 필수적입니다.
echo "--- Cleaning up existing OVS/OVN processes and files ---"
sudo pkill ovsdb-server
sudo pkill ovs-vswitchd
sudo pkill ovn-northd
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

# --- 2. OVN DB 파일 생성 ---
# 이 명령은 DB 파일이 없거나 손상되었을 때만 실행합니다.
# 이미 DB 파일이 존재하면 실행할 필요가 없으며, 기존 데이터를 덮어쓸 수 있습니다.
# 처음 설정할 때만 유용합니다.
echo "--- Creating OVN DB files (if not exist) ---"
sudo ovsdb-tool create /usr/local/etc/ovn/ovnnb_db.db /usr/local/share/ovn/ovn-nb.ovsschema || true
sudo ovsdb-tool create /usr/local/etc/ovn/ovnsb_db.db /usr/local/share/ovn/ovn-sb.ovsschema || true
echo "--- OVN DB files created ---"
echo ""

# --- 3. OVN NB/SB DB 서버 시작 ---
# 백슬래시 뒤에 공백이 없도록 주의하고, 한 줄로 붙여쓰는 것이 안전합니다.
echo "--- Starting OVN NB DB Server ---"
sudo ovsdb-server /usr/local/etc/ovn/ovnnb_db.db --remote=punix:/usr/local/var/run/ovn/ovnnb_db.sock --remote=ptcp:6641:0.0.0.0 --remote=db:OVN_Northbound,NB_Global,connections --pidfile=/usr/local/var/run/ovn/ovnnb-server.pid --detach --log-file=/usr/local/var/log/ovn/ovnnb-server.log
sleep 2 # DB 서버가 완전히 시작될 시간을 줍니다.
echo "--- OVN NB DB Server started ---"
echo ""

echo "--- Starting OVN SB DB Server ---"
sudo ovsdb-server /usr/local/etc/ovn/ovnsb_db.db --remote=punix:/usr/local/var/run/ovn/ovnsb_db.sock --remote=ptcp:6642:0.0.0.0 --remote=db:OVN_Southbound,SB_Global,connections --pidfile=/usr/local/var/run/ovn/ovnsb-server.pid --detach --log-file=/usr/local/var/log/ovn/ovnsb-server.log
sleep 2 # DB 서버가 완전히 시작될 시간을 줍니다.
echo "--- OVN SB DB Server started ---"
echo ""

# --- 4. OVN DB 초기화 ---
# DB 서버가 시작된 후에 실행해야 합니다.
echo "--- Initializing OVN DBs ---"
sudo ovn-nbctl --no-wait init
sudo ovn-sbctl init
echo "--- OVN DBs initialized ---"
echo ""

# --- 5. ovn-northd 데몬 시작 ---
# 로그 파일 경로를 정확히 지정해야 합니다.
echo "--- Starting ovn-northd daemon ---"
sudo ovn-northd --pidfile --detach --log-file=/usr/local/var/log/ovn/ovn-northd.log
sleep 2 # 데몬이 시작될 시간을 줍니다.
echo "--- ovn-northd started ---"
echo ""

# --- 6. OVS 코어 서비스 시작 ---
echo "--- Starting OVS core services ---"
sudo /usr/local/share/openvswitch/scripts/ovs-ctl start
sleep 5 # OVS가 완전히 시작될 시간을 줍니다.
echo "--- OVS core services started ---"
echo ""

# --- 7. 네트워크 설정 변수 ---
# 컨트롤러 노드의 IP와 DNS를 인자로 받습니다.
IP_ADDRESS=$1
DNS=$2
FILE_PATH="/etc/systemd/network/20-br-ext.network"

# --- 8. systemd-networkd 설정 파일 생성 ---
# Gateway를 $DNS로 설정하는 것은 일반적이지 않습니다. Gateway는 보통 라우터의 IP입니다.
# 만약 DNS 서버가 라우터와 동일하다면 괜찮지만, 아니라면 Gateway는 라우터 IP로 명시해야 합니다.
# 여기서는 Gateway를 $DNS로 유지하겠습니다.
echo "--- Creating systemd-networkd config for br-ext ---"
cat <<EOF | sudo tee "$FILE_PATH" > /dev/null
[Match]
Name=br-ext

[Network]
Address=$IP_ADDRESS/24
DNS=$DNS
Gateway=$DNS # Gateway는 보통 라우터 IP입니다. DNS와 다를 수 있습니다.
EOF
echo "--- systemd-networkd config created ---"
echo ""

# --- 9. OVS 브릿지 생성 및 포트 추가 ---
# br-ext가 이미 존재할 경우 오류가 나므로, 먼저 삭제 후 추가하는 것이 안전합니다.
echo "--- Configuring OVS bridges and ports ---"
sudo ovs-vsctl --may-exist del-br br-ext # 이미 존재하면 삭제
sudo ovs-vsctl add-br br-ext
sudo ovs-vsctl add-port br-ext ens3 # 물리 NIC 이름이 ens3인지 다시 확인하세요.
echo "--- OVS bridges and ports configured ---"
echo ""

# --- 10. systemd-networkd 서비스 재시작 ---
# OVS 브릿지 및 포트 설정 후 네트워크 설정을 적용합니다.
echo "--- Restarting systemd-networkd ---"
sudo systemctl restart systemd-networkd
sleep 5 # 네트워크 설정이 적용될 시간을 줍니다.
echo "--- systemd-networkd restarted ---"
echo ""

# --- 11. OVN external-ids 설정 ---
# **주의: 스마트 따옴표 “ ” 대신 일반 따옴표 " " 사용!**
echo "--- Setting OVN external-ids ---"
sudo ovs-vsctl set open_vswitch . external-ids:ovn-remote="tcp:10.5.15.39:6642"
sudo ovs-vsctl set open_vswitch . external-ids:ovn-encap-type=geneve
sudo ovs-vsctl set open_vswitch . external-ids:ovn-nb="tcp:10.5.15.39:6641" # 컨트롤러 노드이므로 NB DB 주소는 자신
sudo ovs-vsctl set open_vswitch . external-ids:ovn-encap-ip="$IP_ADDRESS" # 이 노드의 br-ext IP
sudo ovs-vsctl set open_vswitch . external-ids:ovn-bridge-mappings=UPLINK:br-ext
sudo ovs-vsctl set open_vswitch . external-ids:system-id="$(hostname)" # system-id 명시적 설정
echo "--- OVN external-ids set ---"
echo ""

# --- 12. ovn-controller 데몬 시작 ---
# ovs-ctl start_controller를 사용하면 ovs-ctl이 ovn-controller를 관리합니다.
echo "--- Starting ovn-controller daemon ---"
sudo ovn-ctl start_controller
sleep 5 # 데몬이 시작될 시간을 줍니다.
echo "--- ovn-controller started ---"
echo ""

echo "--- Controller Node Setup Complete ---"
echo "Please check status with: sudo ovs-vsctl show, sudo ovn-sbctl show, sudo lsof -i :6641, sudo lsof -i :6642"
