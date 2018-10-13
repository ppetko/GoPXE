package conf

import (
        "flag"
)

var (
        tftpPath   string
        ksURL      string
        port       string
        bucket     string
        dbName     string
)

func Setup() {
        // Define flags 
        flag.StringVar(&tftpPath, "tftpPath", "/var/lib/tftpboot/pxelinux.cfg/", "tftp conf path e.g /var/lib/tftpboot/pxelinux.cfg/")
        flag.StringVar(&port, "port", "9090", "tcp port")
        flag.StringVar(&bucket, "bucket", "bootactions", "db bucket")
        flag.StringVar(&dbName, "dbName", "gopxe.db", "database file name")
        flag.StringVar(&ksURL, "ksURL", "localhost", "kickstart url")

        // Parsing flags
        flag.Parse()
}

