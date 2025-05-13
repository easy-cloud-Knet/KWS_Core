#!/bin/bash

get_shell_type() {
    ps -p $$ -o cmd= | awk '{print $1}' | xargs basename
}

curr_dir=$(pwd)
shell_type=$(get_shell_type)

wget https://github.com/prometheus/node_exporter/releases/download/v1.9.1/node_exporter-1.9.1.linux-amd64.tar.gz
tar -xaf node_exporter-1.9.1.linux-amd64.tar.gz
mv node_exporter-1.9.1.linux-amd64 node_exporter

wget https://github.com/grafana/loki/releases/latest/download/promtail-linux-amd64.zip
unzip promtail-linux-amd64.zip
sudo mv promtail-linux-amd64 /usr/local/bin/promtail
sudo chmod +x /usr/local/bin/promtail

sudo mkdir -p /etc/promtail

sudo tee /etc/promtail/config.yaml > /dev/null << EOF
server:
  http_listen_port: 9080
  grpc_listen_port: 0

positions:
  filename: /tmp/positions.yaml

clients:
  - url: http://10.5.15.3:3100/loki/api/v1/push

scrape_configs:
  - job_name: core1_system
    static_configs:
      - targets:
          - localhost
        labels:
          job: core1_varlogs
          __path__: /var/log/kws/*

  - job_name: system_logs
    static_configs:
      - targets:
          - localhost
        labels:
          job: core1_system_logs
          __path__: /var/log/{syslog,messages,auth.log,kern.log}
EOF

sudo tee /etc/systemd/system/promtail.service > /dev/null << EOF
[Unit]
Description=Promtail service
After=network.target

[Service]
User=root
Group=root
Type=simple
ExecStart=/usr/local/bin/promtail -config.file /etc/promtail/config.yaml

[Install]
WantedBy=multi-user.target
EOF

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

sudo rm promtail-linux-amd64.zip
sudo rm node_exporter-1.9.1.linux-amd64.tar.gz
sudo systemctl daemon-reload
sudo systemctl enable promtail.service
sudo systemctl start promtail.service
sudo systemctl enable node_exporter.service
sudo systemctl start node_exporter.service

# 셸 설정
if [[ "$shell_type" == "bash" ]]; then
    echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
    echo "셸을 다시 시작하거나 'source ~/.bashrc'를 실행하세요."
elif [[ "$shell_type" == "zsh" ]]; then
    echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.zshrc
    echo "셸을 다시 시작하거나 'source ~/.zshrc'를 실행하세요."
else
    echo "Unsupported shell: $shell_type"
fi

