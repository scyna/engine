package test

import (
	"testing"

	scyna "github.com/scyna/core"
)

func TestGenerateID(t *testing.T) {
	scyna.RemoteInit(scyna.RemoteConfig{
		ManagerUrl: "http://127.0.0.1:8081",
		Name:       "scyna_test",
		Secret:     "123456",
	})
	scyna.UseRemoteLog(3)
	if scyna.ID.Next() == 0 {
		t.Fatal("Can not generate id")
	}

	sn := scyna.InitSerialNumber("test_sn")
	if len(sn.Next()) == 0 {
		t.Fatal("Can not generate serial number")
	}
}
