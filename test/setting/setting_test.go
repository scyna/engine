package setting

import (
	"testing"

	scyna "github.com/scyna/core"
)

func TestWriteSetting(t *testing.T) {
	scyna.RemoteInit(scyna.RemoteConfig{
		ManagerUrl: "http://127.0.0.1:8081",
		Name:       "scyna_test",
		Secret:     "123456",
	})
	scyna.UseRemoteLog(3)
	scyna.Settings.Write("test", "test")
	if ok, val := scyna.Settings.ReadString("test"); !ok || val != "test" {
		t.Fatal("Can not write setting")
	}
}
