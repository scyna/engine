package scheduler

import (
	"errors"
	"log"
	"math"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/gocql/gocql"
	scyna "github.com/scyna/core"
	scyna_const "github.com/scyna/core/const"
	scyna_proto "github.com/scyna/core/proto/generated"
)

func StartTask(ctx *scyna.Endpoint, request *scyna_proto.StartTaskRequest) scyna.Error {

	log.Println("Receive StartTaskRequest")

	if err := validateStartTaskRequest(request); err != nil {
		ctx.Error(err.Error())
		return scyna.REQUEST_INVALID
	}

	//TODO: Check module exits

	if request.Loop == 0 {
		request.Loop = math.MaxInt64
	}

	// Insert new task to scyna.task table
	taskID := scyna.ID.Next()
	start := time.Unix(request.Time, 0)
	qBatch := scyna.DB.Session.NewBatch(gocql.LoggedBatch)
	qBatch.Query("INSERT INTO "+scyna_const.TASK_TABLE+"(id, topic, data, start, next, interval, loop_count, loop_index, done) "+
		"VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);",
		taskID, request.Topic, request.Data, start, start, request.Interval, request.Loop, 0, false)

	qBatch.Query("INSERT INTO "+scyna_const.MODULE_HAS_TASK_TABLE+"(module, task_id) VALUES (?, ?);", request.Module, taskID)

	bucket := GetBucket(start) // Generate period id
	qBatch.Query("INSERT INTO "+scyna_const.TODO_TABLE+"(bucket, task_id) VALUES (?, ?);", bucket, taskID)
	if err := scyna.DB.Session.ExecuteBatch(qBatch); err != nil {
		ctx.Error(err.Error())
		return scyna.REQUEST_INVALID
	}

	return ctx.OK(&scyna_proto.StartTaskResponse{Id: taskID})
}

func validateStartTaskRequest(request *scyna_proto.StartTaskRequest) error {
	if int64(request.Time) < time.Now().Unix() {
		return errors.New("task time is less than now")
	}
	if request.Interval < 60 {
		return errors.New("interval must be greater than 60 second")
	}
	return validation.ValidateStruct(request,
		validation.Field(&request.Topic, validation.Required, validation.Length(1, 255)),
		validation.Field(&request.Module, validation.Required, validation.Length(1, 255)),
	)
}
