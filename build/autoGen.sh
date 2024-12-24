#!/bin/bash

METADATA=$1


mkdir /var/lib/kws/${METADATA}

#making user-data meta-data here 

qemu-img create -b /var/lib/kws/baseimg/ubuntu-cloud-24.04.img -f qcow2 -F qcow2 "/var/lib/kws/${METADATA}/${METADATA}.qcow2"  10G



genisoimage --output /var/lib/kws/${METADATA}/cidata.iso -V cidata -r -J /var/lib/kws/${METADATA}/user-data /var/lib/kws/${METADATA}/meta-data


