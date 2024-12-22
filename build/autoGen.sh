#!/bin/bash

DATA=$1

echo "$DATA"

qemu-img create -b /var/lib/kws/baseimg/ubuntu-cloud-24.04.img -f qcow2 -F qcow2 "./user${DATA}.qcow2"  10G



genisoimage --output cidata.iso -V cidata -r -J user-data meta-data


