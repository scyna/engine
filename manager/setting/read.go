package setting

import (
	"log"

	"github.com/scylladb/gocqlx/v2/qb"
	scyna "github.com/scyna/core"
	scyna_proto "github.com/scyna/core/proto/generated"
)

func Read(s *scyna.Endpoint, request *scyna_proto.ReadSettingRequest) scyna.Error {
	log.Println("Receive ReadSettingRequest")

	var value string
	if err := qb.Select("scyna.setting").
		Columns("value").
		Where(qb.Eq("context"), qb.Eq("key")).
		Limit(1).
		Query(scyna.DB).
		Bind(request.Context, request.Key).
		GetRelease(&value); err != nil {
		s.Logger.Info("Can not read setting - " + err.Error())
		return scyna.REQUEST_INVALID
	}

	s.Response(&scyna_proto.ReadSettingResponse{Value: value})
	return scyna.OK
}
