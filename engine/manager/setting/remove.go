package setting

import (
	"github.com/scylladb/gocqlx/v2/qb"
	scyna "github.com/scyna/core"
)

func Remove(s *scyna.Service, request *scyna.RemoveSettingRequest) {
	if err := qb.Delete("scyna.setting").
		Where(qb.Eq("context"), qb.Eq("key")).
		Query(scyna.DB).
		Bind(request.Context, request.Key).ExecRelease(); err != nil {
		s.Error(scyna.SERVER_ERROR)
		return
	}

	s.Done(scyna.OK)

	// s.EmitSignal(scyna.SETTING_REMOVE_CHANNEL+request.Module, &scyna.SettingUpdatedSignal{
	// 	Key: request.Key,
	// })
}
