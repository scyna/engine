package module

import (
	"log"

	scyna "github.com/scyna/core"
	scyna_const "github.com/scyna/core/const"
	PROTO "github.com/scyna/engine/proto/generated"
)

func ListSettingHandler(ctx *scyna.Endpoint, request *PROTO.ListSettingRequest) scyna.Error {
	log.Println("Receive ListSettingRequest")

	if err := scyna.DB.AssureExists("SELECT code FROM "+scyna_const.MODULE_TABLE+" WHERE code = ?",
		request.Module); err != nil {
		ctx.Error(err.Error())
		return scyna.REQUEST_INVALID
	}

	rs := scyna.DB.QueryMany("SELECT key,value FROM "+scyna_const.SETTING_TABLE+" WHERE module_code = ?", request.Module)
	response := &PROTO.ListSettingResponse{}
	for rs.Next() {
		var key string
		var value string
		if err := rs.Scan(&key, &value); err != nil {
			ctx.Error(err.Error())
			return scyna.SERVER_ERROR
		}

		response.Items = append(response.Items, &PROTO.SettingItem{Key: key, Value: value})
	}

	return ctx.OK(response)
}
