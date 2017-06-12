package main

import (
    "fmt"
    "log"
    "net/http"
    "os"

    "github.com/gorilla/mux"
)

var (
    tftp_conf string = "/var/lib/tftpboot/pxelinux.cfg/"
    ks_server string = "Kickstart_File" // https://github.com/mycodinglab/kickstart-generator
)

func main() {
    router := mux.NewRouter().StrictSlash(true)
    router.HandleFunc("/", Index)
    router.HandleFunc("/status", Status)
    router.HandleFunc("/pxe/{hostname}/{uuid}", configure_pxe)
    log.Fatal(http.ListenAndServe(":9090", router))
}

func Status(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "API is up and running")
}

func Index(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Welcome to PXEaaS")
}

func configure_pxe(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    hostname := vars["hostname"]
    uuid := vars["uuid"]
    file_path := tftp_conf + uuid

    file, err := os.OpenFile(file_path, os.O_RDWR|os.O_CREATE, 0644)
    if err != nil {
        fmt.Printf("Error %v", err)
        w.WriteHeader(http.StatusBadRequest)
    }
    defer file.Close()

    append := "default linux\nlabel linux\nMENU LABEL CentOS 7\nKERNEL vmlinuz\nAPPEND ksdevice=bootif load_ramdisk=1 initrd=initrd.img network ks=" + ks_server + hostname + " " + "ksdevice=bootif biosdevname=0"

    _, err = file.WriteString(append)
    if err != nil {
        fmt.Printf("Error %v", err)
        w.WriteHeader(http.StatusBadRequest)
    }
    defer file.Close()

    err = file.Sync()
    if err != nil {
        fmt.Fprintln(w, err)
        w.WriteHeader(http.StatusBadRequest)
    } else {
        fmt.Fprintln(w, "success")
        w.WriteHeader(http.StatusOK)
    }
}
