package setting

import (
	scyna "github.com/scyna/go"

	"github.com/scylladb/gocqlx/v2/qb"
)

func Remove(s *scyna.Service) {
	var request scyna.RemoveSettingRequest
	if !s.Parse(&request) {
		return
	}

	if err := qb.Delete("scyna.setting").
		Where(qb.Eq("module_code"), qb.Eq("key")).
		Query(scyna.DB).
		Bind(request.Module, request.Key).ExecRelease(); err != nil {
		s.Error(scyna.SERVER_ERROR)
		return
	}

	s.Done(scyna.OK)

	scyna.EmitSignal(scyna.SETTING_REMOVE_CHANNEL+request.Module, &scyna.SettingUpdatedSignal{
		Key:   request.Key,
	})
}
