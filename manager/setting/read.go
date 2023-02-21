package setting

import (
	"log"

	"github.com/scylladb/gocqlx/v2/qb"
	scyna "github.com/scyna/core"
	scyna_proto "github.com/scyna/core/proto/generated"
)

func Read(s *scyna.Context, request *scyna_proto.ReadSettingRequest) scyna.Error {
	log.Println("Receive ReadSettingRequest")

	var value string
	if err := qb.Select("scyna.setting").
		Columns("value").
		Where(qb.Eq("module"), qb.Eq("key")).
		Limit(1).
		Query(scyna.DB).
		Bind(request.Module, request.Key).
		GetRelease(&value); err != nil {
		s.Info("Can not read setting - " + err.Error())
		return scyna.REQUEST_INVALID
	}

	return s.OK(&scyna_proto.ReadSettingResponse{Value: value})
}
