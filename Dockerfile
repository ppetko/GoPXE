FROM centos:7.3.1611

RUN yum install tftp tftp-server* xinetd* dhcp* -y && \
   yum clean all

RUN mkdir /var/lib/tftpboot/pxelinux.cfg

EXPOSE 67 67/udp 69/udp 9090 9090/udp

ADD ./conf/dhcpd_template /etc/dhcp/dhcpd.conf
ADD ./conf/tftp /etc/xinetd.d/tftp
ADD ./pxe_conf /var/lib/tftpboot
ADD ./pxe/pxe_api /
ADD ./run_pxe.sh /

HEALTHCHECK --interval=4m --timeout=3s CMD curl --fail http://localhost:9090/status || exit 1

CMD ["/run_pxe.sh"]
