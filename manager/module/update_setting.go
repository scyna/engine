package module

import (
	"log"

	scyna "github.com/scyna/core"
	scyna_const "github.com/scyna/core/const"
	scyna_proto "github.com/scyna/core/proto/generated"
	PROTO "github.com/scyna/engine/proto/generated"
)

func UpdateSettingHandler(ctx *scyna.Endpoint, request *PROTO.UpdateSettingRequest) scyna.Error {
	log.Println("Receive UpdateSettingRequest")

	if err := scyna.DB.Execute("INSERT INTO "+scyna_const.SETTING_TABLE+
		" (module, key, value) VALUES (?, ?, ?)",
		request.Module, request.Key, request.Value); err != nil {
		ctx.Error("WriteSetting: " + err.Error())
		return scyna.REQUEST_INVALID
	}

	scyna.EmitSignal(scyna_const.SETTING_UPDATE_CHANNEL+request.Module, &scyna_proto.SettingUpdatedSignal{
		Key:   request.Key,
		Value: request.Value,
	})

	return scyna.OK
}
