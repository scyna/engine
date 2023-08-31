package test

import (
	"testing"

	scyna "github.com/scyna/core"
)

func TestCreateLog(t *testing.T) {
	scyna.RemoteInit(scyna.RemoteConfig{
		ManagerUrl: "http://127.0.0.1:8081",
		Name:       "scyna_test",
		Secret:     "123456",
	})

	scyna.UseRemoteLog(3)
	scyna.Session.Info("Test Info log")
	scyna.Session.Error("Test Error log")
	scyna.Session.Warning("Test Warning log2")
}
