package setting

import (
	"log"

	scyna "github.com/scyna/go"

	"github.com/scylladb/gocqlx/v2/qb"
)

func Read(s *scyna.Service) {
	log.Println("Receive ReadSettingRequest")
	var request scyna.ReadSettingRequest
	if !s.Parse(&request) {
		return
	}

	var value string
	if err := qb.Select("scyna.setting").
		Columns("value").
		Where(qb.Eq("module_code"), qb.Eq("key")).
		Limit(1).
		Query(scyna.DB).
		BindStruct(&request).
		GetRelease(&value); err != nil {
		log.Println(err.Error())
		s.Error(scyna.REQUEST_INVALID)
		return
	}

	s.Done(&scyna.ReadSettingResponse{Value: value})
}
