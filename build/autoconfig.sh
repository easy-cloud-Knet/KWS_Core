#!/bin/bash

get_shell_type() {
    ps -p $$ -o cmd= | awk '{print $1}' | xargs basename
}

curr_dir=$(pwd)
shell_type=$(get_shell_type)

wget https://github.com/prometheus/node_exporter/releases/download/v1.9.1/node_exporter-1.9.1.linux-amd64.tar.gz
tar -xaf node_exporter-1.9.1.linux-amd64.tar.gz
mv node_exporter-1.9.1.linux-amd64 node_exporter

sudo tee /etc/systemd/system/node_exporter.service > /dev/null << EOF
[Unit]
Description=Node Exporter
Wants=network-online.target
After=network-online.target

[Service]
User=root
Group=root
Type=simple 
ExecStart=${curr_dir}/node_exporter/node_exporter

[Install]
WantedBy=multi-user.target
EOF

user=$(whoami)
sudo tee /etc/systemd/system/kws_core.service > /dev/null <<EOF
[Unit]
Description=kws daemon service for host computer
Wants=network-online.target
After=network.target
StartLimitIntervalSec=0
[Service]
Type=simple
Restart=always
RestartSec=1
User=${user}
ExecStart=${curr_dir}/KWS_Core

[Install]
WantedBy=multi-user.target
EOF


sudo rm node_exporter-1.9.1.linux-amd64.tar.gz
sudo systemctl daemon-reload
sudo systemctl enable node_exporter.service
sudo systemctl start node_exporter.service

if [[ "$shell_type" == "bash" ]]; then
    echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
    echo "셸을 다시 시작하거나 'source ~/.bashrc'를 실행하세요."
elif [[ "$shell_type" == "zsh" ]]; then
    echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.zshrc
    echo "셸을 다시 시작하거나 'source ~/.zshrc'를 실행하세요."
else
    echo "Unsupported shell: $shell_type"
fi

