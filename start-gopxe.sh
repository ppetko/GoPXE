#!/bin/bash 

DATE=`date +%Y-%m-%d-%H:%M:%S`
DHCPD_CONF="/etc/dhcp/dhcpd.conf"
TFTPD_CONF="/etc/xinetd.d/tftp"
TFTPD_BOOT_PATH="/var/lib/tftpboot"

function log(){
    echo "$DATE $@" 
    return 0
}

function warn(){
    echo "$DATE $@" 
    return 1
}

function panic(){
    echo "$DATE $@" 
    exit 1
}

## Starting goPXE
log "gopxe is starting..."
/gopxe/main -ksURL $(hostname -I | awk '{print $1}') & 

## Start dhcpd
log "starting dhcpd"
/usr/sbin/dhcpd -4 -f -d --no-pid -cf ${DHCPD_CONF} & 

## Start tftp
log "starting tftpd"
/usr/sbin/in.tftpd --foreground --address 0.0.0.0:69 --secure ${TFTPD_BOOT_PATH}
