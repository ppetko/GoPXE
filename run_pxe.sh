#!/bin/bash 

## Start pxe_api
/pxe_api &

## Start dhcpd
/usr/sbin/dhcpd -4 -f -d --no-pid -cf /etc/dhcp/dhcpd.conf & 

## Start tftp
/usr/sbin/in.tftpd --foreground --address 0.0.0.0:69 --secure /var/lib/tftpboot
