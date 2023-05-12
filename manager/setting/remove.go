package setting

import (
	"github.com/scylladb/gocqlx/v2/qb"
	scyna "github.com/scyna/core"
	scyna_const "github.com/scyna/core/const"
	scyna_proto "github.com/scyna/core/proto/generated"
)

func Remove(ctx *scyna.Endpoint, request *scyna_proto.RemoveSettingRequest) scyna.Error {
	if err := qb.Delete(scyna_const.SETTING_TABLE).
		Where(qb.Eq("module"), qb.Eq("key")).
		Query(scyna.DB).
		Bind(request.Module, request.Key).ExecRelease(); err != nil {
		return scyna.SERVER_ERROR
	}

	// s.EmitSignal(scyna.SETTING_REMOVE_CHANNEL+request.Module, &scyna.SettingUpdatedSignal{
	// 	Key: request.Key,
	// })

	return scyna.OK
}
