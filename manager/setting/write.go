package setting

import (
	"log"

	"github.com/scylladb/gocqlx/v2/qb"
	scyna "github.com/scyna/core"
)

func Write(s *scyna.Endpoint, request *scyna.WriteSettingRequest) {
	log.Println("Receive WriteSettingRequest")

	if err := qb.Insert("scyna.setting").
		Columns("context", "key", "value").
		Query(scyna.DB).
		Bind(request.Context, request.Key, request.Value).
		ExecRelease(); err != nil {
		s.Logger.Error("WriteSetting: " + err.Error())
		s.Error(scyna.REQUEST_INVALID)
		return
	}

	s.Done(scyna.OK)

	// s.EmitSignal(scyna.SETTING_UPDATE_CHANNEL+request.Module, &scyna.SettingUpdatedSignal{
	// 	Key:   request.Key,
	// 	Value: request.Value,
	// })
}
