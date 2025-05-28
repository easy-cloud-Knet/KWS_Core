#!/bin/bash
wget  https://go.dev/dl/go1.22.6.linux-amd64.tar.gz 
rm -rf /usr/local/go && tar -C /usr/local -xzf go1.22.6.linux-amd64.tar.gz


sudo ln -s /usr/local/go/bin/go /usr/local/bin/go

rm go1.22.6.linux-amd64.tar.gz
