package setting

import (
	"log"

	scyna "github.com/scyna/go"

	"github.com/scylladb/gocqlx/v2/qb"
)

func Write(s *scyna.Service) {
	log.Println("Receive WriteSettingRequest")
	var request scyna.WriteSettingRequest
	if !s.Parse(&request) {
		return
	}

	if applied, err := qb.Insert("scyna.setting").
		Columns("module_code", "key", "value").
		Unique().
		Query(scyna.DB).
		Bind(request.Module, request.Key, request.Value).
		ExecCASRelease(); !applied {
		if err != nil {
			s.LOG.Error("WriteSetting: " + err.Error())
			s.Error(scyna.REQUEST_INVALID)
			return
		}
		return
	}

	s.Done(scyna.OK)

	scyna.EmitSignal(scyna.SETTING_UPDATE_CHANNEL+request.Module, &scyna.SettingUpdatedSignal{
		Key:   request.Key,
		Value: request.Value,
	})
}
