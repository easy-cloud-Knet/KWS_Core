#!/bin/bash

sudo apt install make libssl-dev unbound -y autoconf automake libtool  
git clone https://github.com/ovn-org/ovn.git
cd ovn
touch /dev/urandom
touch /dev/net/tun


sudo ./boot.sh

#updating submodule from ovn github, can be replaced if advance feature of ovs are needed

git submodule update --init
cd ovs
sudo ./boot.sh
sudo ./configure
sudo make
sudo make install

cd ..
sudo ./configure
sudo make
sudo make install 


echo 'PATH="$PATH:/usr/local/share/ovn/scripts"'>> ~/.bashrc  
echo 'PATH="$PATH:/usr/local/share/openvswitch/scripts"'>> ~/.bashrc  


ovn-ctl start_controller


