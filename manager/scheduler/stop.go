package scheduler

import (
	"github.com/scylladb/gocqlx/v2/qb"
	scyna "github.com/scyna/core"
)

func StopTask(s *scyna.Endpoint, request *scyna.StopTaskRequest) {

	if err := qb.Update("scyna.task").
		Set("done").
		Where(qb.Eq("id")).
		Query(scyna.DB).
		Bind(true, request.Id).
		ExecRelease(); err != nil {
		s.Error(scyna.REQUEST_INVALID)
		s.Logger.Error(err.Error())
		return
	}

	s.Done(scyna.OK)
}
