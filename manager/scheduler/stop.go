package scheduler

import (
	"github.com/scylladb/gocqlx/v2/qb"
	scyna "github.com/scyna/core"
	scyna_const "github.com/scyna/core/const"
	scyna_proto "github.com/scyna/core/proto/generated"
)

func StopTask(ctx *scyna.Endpoint, request *scyna_proto.StopTaskRequest) scyna.Error {

	if err := qb.Update(scyna_const.TASK_TABLE).
		Set("done").
		Where(qb.Eq("id")).
		Query(scyna.DB).
		Bind(true, request.Id).
		ExecRelease(); err != nil {
		ctx.Error(err.Error())
		return scyna.REQUEST_INVALID
	}

	return scyna.OK
}
