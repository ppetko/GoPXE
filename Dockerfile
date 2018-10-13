FROM golang:latest as builder 
RUN mkdir -p /go/src/github.com/ppetko/gopxe
ADD . /go/src/github.com/ppetko/gopxe
WORKDIR /go/src/github.com/ppetko/gopxe
RUN go test ./...
RUN go build -o main .

FROM centos:7.5.1804
RUN yum install -y tftp tftp-server* xinetd* dhcp* epel-release syslinux  && yum clean all 
EXPOSE 67 67/udp 69/udp 9090 9090/udp
RUN mkdir /var/lib/tftpboot/pxelinux.cfg /opt/localrepo
RUN cp -r /usr/share/syslinux/pxelinux.0 /var/lib/tftpboot
ADD ./pxebootImages /var/lib/tftpboot
RUN mkdir -p /gopxe/public ; mkdir /gopxe/ksTempl
WORKDIR /gopxe
COPY --from=builder /go/src/bitbucket.org/ppetkov85/docker-pxe/main /gopxe/
ADD ./public /gopxe/public
ADD ./ksTempl /gopxe/ksTempl
ADD ./start-gopxe.sh /gopxe/
HEALTHCHECK --interval=4m --timeout=60s CMD curl --fail http://localhost:9090/health || exit 1
ENTRYPOINT ["/gopxe/start-gopxe.sh"]
