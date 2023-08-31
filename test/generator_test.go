package test

import (
	"log"
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
	log.Print(scyna.ID.Next())

	sn := scyna.InitSerialNumber("test_sn")

	log.Print(sn.Next())
}
