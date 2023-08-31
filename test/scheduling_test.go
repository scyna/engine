package test

import (
	"testing"

	scyna "github.com/scyna/core"
)

func TestScheduling(t *testing.T) {
	scyna.RemoteInit(scyna.RemoteConfig{
		ManagerUrl: "http://127.0.0.1:8081",
		Name:       "scyna_test",
		Secret:     "123456",
	})
	scyna.UseRemoteLog(3)
}
