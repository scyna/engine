package setting

import (
	"log"

	scyna "github.com/scyna/core"
	scyna_const "github.com/scyna/core/const"
	scyna_proto "github.com/scyna/core/proto/generated"
)

func Read(ctx *scyna.Endpoint, request *scyna_proto.ReadSettingRequest) scyna.Error {
	log.Println("Receive ReadSettingRequest")

	var value string
	if err := scyna.DB.QueryOne("SELECT value FROM "+scyna_const.SETTING_TABLE+
		" WHERE module = ? AND key = ? LIMIT 1",
		request.Module, request.Key).Scan(&value); err != nil {
		ctx.Info("Can not read setting - " + err.Error())
		return scyna.REQUEST_INVALID
	}

	return ctx.OK(&scyna_proto.ReadSettingResponse{Value: value})
}
