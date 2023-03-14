package setting

import (
	"log"

	"github.com/scylladb/gocqlx/v2/qb"
	scyna "github.com/scyna/core"
	scyna_proto "github.com/scyna/core/proto/generated"
)

func Write(ctx *scyna.Endpoint, request *scyna_proto.WriteSettingRequest) scyna.Error {
	log.Println("Receive WriteSettingRequest")

	if err := qb.Insert("scyna.setting").
		Columns("module", "key", "value").
		Query(scyna.DB).
		Bind(request.Module, request.Key, request.Value).
		ExecRelease(); err != nil {
		ctx.Error("WriteSetting: " + err.Error())
		return scyna.REQUEST_INVALID
	}

	// s.EmitSignal(scyna.SETTING_UPDATE_CHANNEL+request.Module, &scyna.SettingUpdatedSignal{
	// 	Key:   request.Key,
	// 	Value: request.Value,
	// })

	return scyna.OK
}
