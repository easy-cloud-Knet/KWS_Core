#!/bin/bash
wget  https://go.dev/dl/go1.22.6.linux-amd64.tar.gz 
rm -rf /usr/local/go && tar -C /usr/local -xzf go1.22.6.linux-amd64.tar.gz

export PATH=$PATH:/usr/local/go/bin

go env -w CGO_ENABLED=1 

