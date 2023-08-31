package setting

import (
	scyna "github.com/scyna/core"
	scyna_const "github.com/scyna/core/const"
	scyna_proto "github.com/scyna/core/proto/generated"
)

func Remove(ctx *scyna.Endpoint, request *scyna_proto.RemoveSettingRequest) scyna.Error {

	if err := scyna.DB.Execute("DELETE FROM "+scyna_const.SETTING_TABLE+
		" WHERE module = ? AND key = ?", request.Module, request.Key); err != nil {
		return scyna.SERVER_ERROR
	}

	scyna.EmitSignal(scyna_const.SETTING_REMOVE_CHANNEL+request.Module, &scyna_proto.SettingUpdatedSignal{
		Key: request.Key,
	})

	return scyna.OK
}
