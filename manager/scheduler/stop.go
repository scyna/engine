package scheduler

import (
	"github.com/scylladb/gocqlx/v2/qb"
	scyna "github.com/scyna/core"
	scyna_proto "github.com/scyna/core/proto/generated"
)

func StopTask(s *scyna.Endpoint, request *scyna_proto.StopTaskRequest) scyna.Error {

	if err := qb.Update("scyna.task").
		Set("done").
		Where(qb.Eq("id")).
		Query(scyna.DB).
		Bind(true, request.Id).
		ExecRelease(); err != nil {
		s.Error(err.Error())
		return scyna.REQUEST_INVALID
	}

	return scyna.OK
}
