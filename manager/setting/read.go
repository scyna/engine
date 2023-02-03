package setting

import (
	"log"

	"github.com/scylladb/gocqlx/v2/qb"
	scyna "github.com/scyna/core"
	scyna_engine "github.com/scyna/core/engine"
)

func Read(s *scyna.Endpoint, request *scyna_engine.ReadSettingRequest) scyna.Error {
	log.Println("Receive ReadSettingRequest")

	var value string
	if err := qb.Select("scyna.setting").
		Columns("value").
		Where(qb.Eq("module"), qb.Eq("key")).
		Limit(1).
		Query(scyna.DB).
		Bind(request.Module, request.Key).
		GetRelease(&value); err != nil {
		s.Logger.Info("Can not read setting - " + err.Error())
		return scyna.REQUEST_INVALID
	}

	s.Response(&scyna_engine.ReadSettingResponse{Value: value})
	return scyna.OK
}
