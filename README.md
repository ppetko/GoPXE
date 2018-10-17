# GoPXE Project

## Introduction
GoPXE is a dynamic build system for installing various operating systems on virtual machines and bare metal servers using pxe boot, tftp and dhcp wrapped into docker container orchestrated by APIs. 

## Quickstart

Download our docker image 

```
docker pull ppetko/gopxe
```

### Configuration - edit the configuration files accordingly. Sample configs has been provided in the repo. 

```
$ vi ./conf/dhcpd.conf 
$ vi ./conf/tftpd.conf 
```

### Start GoPXE using Docker image 

```
sudo docker run --rm --net=host --name goPXE -td \
            --mount type=bind,source="$(pwd)"/conf/dhcpd.conf,target=/etc/dhcp/dhcpd.conf \
            --mount type=bind,source="$(pwd)"/conf/tftpd.conf,target=/etc/xinetd.d/tftp \
            ppetko/gopxe
```

## Build your own docker image 

```
$ go get github.com/ppetko/gopxe
$ cd $GOPATH/src/github.com/ppetko/gopxe
$ sudo make docker-build

```

## How PXE works?

* The target host (the PXE client) is booted.
* The target host makes a DHCP request.
* The DHCP server responds with the IP information and provides information about the location of a TFTP server.
* When the client receives the information, it contacts the TFTP server requesting the file that the DHCP server specified (in this case, the network boot loader).
* The TFTP server sends the network boot loader, and the client executes it.
* PXELINUX searches for a configuration file on the TFTP server, and boots a kernel according to that configuration file (centos7/vmlinux and centos7/vmlizimage). In this case, the configuration file instructs PXE to load the kernel (vmlinuz) and a ramdisk (initrd.img).
* The client downloads the files it needs and then loads them.
* The system boots the OS installer using our Kickstart service.
* The installer runs interactively or scripted, as directed by the PXE configuration file.
* The installer uses remote repository, or locally content from ISO file.
* OS is installed.

## APIs Reference Examples

Note: Parameters are specific to your environment.

### Create BootAction

* PXEBoot images are already configured for centos7 in pxebootImages/centos7, so these options are valid kernel:centos7/vmlinux and initrd:centos7/vmlizimage. You can add different specific version or different OS images in pxebootImages. 
* myfirstbootaction is the name of your boot action. You can create many bootactions with specific parameters as long as the names are different. 

```
curl -vv -H "Content-Type: application/json" -d \
'{
  "default": "linux",
  "label": "linux",
  "menu": "centos7",
  "kernel": "centos7/vmlinuz",
  "ksdevice": "link",
  "ip": "dhcp",
  "load_ramdisk": "1",
  "initrd":"centos7/initrd.img"
}' -X POST localhost:9090/bootaction/myfirstbootaction

```

* bootaction: myfirstbootaction will reference your bootaction created earlier. 
* ksfile - this will reference the default kickstart which is preload. You can add you own kickstart file in ksTempl/
* os and version - will be used to build repository in the default kickstart. If you are not using the default kickstart file, leave these options empty.  
* If you want to PXEBoot using mac address instead of UUID, use "uuid": "01-YOUR-MAC-ADDRESS" option.

### Create basic PXEBoot record 

```
curl -vv  -H "Content-Type: application/json" -d \
'{
  "bootaction": "myfirstbootaction",
  "ksfile": "default",
  "os": "centos-7",
  "version": "7.5.1804",
  "hostname": "test-myvm.local",
  "uuid": "42330d5a-0ead-f7fa-4e3a-ae3bdcb08c69"
}' http://localhost:9090/pxeboot

```

* hostname, ip, mask, ns1, ns2, gw, will be used to build network configuration in the default kickstart

### Create PXEBoot record with network options 

* You can view all bxeboot at http://localhost:9090/pxelinux/

```
curl -vv  -H "Content-Type: application/json" -d \
'{
  "bootaction": "vmtemplate",
  "ksfile": "default",
  "os": "centos-7",
  "version": "7.5.1804",
  "hostname": "test-myvm1.local",
  "ip": "10.1.20.50",
  "mask": "255.255.255.0",
  "ns1": "8.8.8.8",
  "ns2": "8.8.4.4",
  "gw": "10.1.20.1",
  "uuid": "42330d5a-0ead-f7fa-4e3a-ae3bdcb08c69"
}' http://localhost:9090/pxeboot

```

#### Results

```
default linux
 label linux
 MENU LABEL centos7
 KERNEL centos7/vmlinuz
 APPEND ksdevice=link ip=dhcp load_ramdisk=1 initrd=centos7/initrd.img ks=http://localhost:9090/kickstart/?name=default&os=centos-7&version=7.5.1804&fqdn=test-myvm1.local&ip=10.1.20.50&mask=255.255.255.0&gw=255.255.255.0&ns1=8.8.8.8&ns2=8.8.4.4

```

## Setup local install images 

```
# wget http://mirror.cc.columbia.edu/pub/linux/centos/7/isos/x86_64/CentOS-7-x86_64-Minimal-1804.iso
# mkdir /mnt/iso
# mount -t iso9660 -o loop CentOS-7-x86_64-Minimal-1804.iso  /mnt/iso/

```

#### Once you have all files in /mnt/iso/, start and mount /mnt/iso/ inside the docker container. 

```
sudo docker run --rm --net=host --name goPXE -td \
            --mount type=bind,source="$(pwd)"/conf/dhcpd.conf,target=/etc/dhcp/dhcpd.conf \
            --mount type=bind,source="$(pwd)"/conf/tftpd.conf,target=/etc/xinetd.d/tftp \
            --mount type=bind,source="/mnt/iso/",target=/opt/localrepo \
            docker-pxe:latest
```

## TODO
- [ ] Add function documentation and godocs generation.
- [ ] Add test code coverage for the rest of the code base.
- [ ] Add piplene build for the project.

## RoadMap
- [ ] Add pxelinux configuration page https://golangcode.com/download-a-file-from-a-url/
- [ ] Create Ansible  hook - API endpoint that accepts ansible run configs per specific host. 
- [ ] Create status output of the job and perhaps synch the results back to the db. 
- [ ] Add Status dashboard
- [ ] Manage(start, stop, restart) dhcpd and tftpd processes from the main application instead of bash 
- [ ] Output logs to html web page from dhcpd and the rest of the app so the user could easily troubleshoot.
- [ ] Add Notifications/Status API (calling home with errors type)
- [ ] Add multi OS support. Add paramerts for pxelinux.0. Currenly we can configure only one OS at the time. 

## Notes
- [ ] Currently suported OS is CentOS/RedHat by the default installation. But you could reconfigure GoPXE for any other OS of choice. 

## Pull requests welcome!
Spotted an error? You have good idea on improvment? Send me a [pull request](github.com/ppetko/gopxe/pulls)! Thanks. 

## License
[Creative Commons Attribution License](http://creativecommons.org/licenses/by/2.0/)
