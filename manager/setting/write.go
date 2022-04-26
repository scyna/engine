package setting

import (
	"log"

	scyna "github.com/scyna/go/scyna"

	"github.com/scylladb/gocqlx/v2/qb"
)

func Write(s *scyna.Context, request *scyna.WriteSettingRequest) {
	log.Println("Receive WriteSettingRequest")

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
