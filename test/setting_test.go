package test

import (
	"testing"

	scyna "github.com/scyna/core"
	scyna_const "github.com/scyna/core/const"
)

func TestWriteSetting(t *testing.T) {
	scyna.RemoteInit(scyna.RemoteConfig{
		ManagerUrl: "http://127.0.0.1:8081",
		Name:       "scyna_test",
		Secret:     "123456",
	})
	scyna.UseRemoteLog(3)
	value := "test"
	scyna.Settings.Write("test", value)
	if ok, val := scyna.Settings.ReadString("test"); !ok || val != value {
		t.Fatal("Can not write setting")
	}

	scyna.DB.Execute("DELETE FROM "+scyna_const.SETTING_TABLE+" WHERE module=? AND key = ?", "scyna_test", "test")
}
