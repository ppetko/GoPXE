package db

import (
	"testing"
)

var (
	tt BoltDB
	bb = "test-bucket"
	kk = "vm-boot-action"
	aa = "default linux label linux MENU LABEL CentOS 7 KERNEL centos74/vmlinuz APPEND ksdevice=link ip=dhcp load_ramdisk=1 initrd=centos74/initrd.img"
)

func TestCreateBucket(t *testing.T) {
	err := tt.CreateBucket(bb)
	if err != nil {
		t.Fatalf("Failed to create database bucket %s", err)
	}
}

func TestPutBootAction(t *testing.T) {
	err := tt.PutBootAction(bb, kk, aa)
	if err != nil {
		t.Fatalf("Failed to create boot action %s", err)
	}
}

func TestGetBootAction(t *testing.T) {
	actualResult := aa
	err, expectedResult := tt.GetBootAction(bb, kk)
	if err != nil {
		t.Fatalf("Failed to retrieve boot action %s", err)
	}

	if actualResult != expectedResult {
		t.Fatalf("Expected %s but got %s", expectedResult, actualResult)
	}
}

func TestGetAllBootActions(t *testing.T) {
	err, _ := tt.GetAllBootActions(bb)
	if err != nil {
		t.Fatalf("Failed to retrieve all boot action %s", err)
	}
}
