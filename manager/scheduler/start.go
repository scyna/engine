package scheduler

import (
	"errors"
	"math"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/gocql/gocql"
	scyna "github.com/scyna/core"
	scyna_proto "github.com/scyna/core/proto/generated"
)

func StartTask(s *scyna.Endpoint, request *scyna_proto.StartTaskRequest) scyna.Error {
	if err := validateStartTaskRequest(request); err != nil {
		s.Logger.Error(err.Error())
		return scyna.REQUEST_INVALID
	}

	//TODO: Check module exits

	if request.Loop == 0 {
		request.Loop = math.MaxInt64
	}

	// Insert new task to scyna.task table
	taskID := scyna.ID.Next()
	start := time.Unix(request.Time, 0)
	qBatch := scyna.DB.NewBatch(gocql.LoggedBatch)
	qBatch.Query("INSERT INTO scyna.task(id, topic, data, start, next, interval, loop_count, loop_index, done) "+
		"VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);",
		taskID, request.Topic, request.Data, start, start, request.Interval, request.Loop, 0, false)

	qBatch.Query("INSERT INTO scyna.module_has_task(module, task_id) VALUES (?, ?);", request.Module, taskID)

	bucket := GetBucket(start) // Generate period id
	qBatch.Query("INSERT INTO scyna.todo(bucket, task_id) VALUES (?, ?);", bucket, taskID)
	if err := scyna.DB.ExecuteBatch(qBatch); err != nil {
		s.Logger.Error(err.Error())
		return scyna.REQUEST_INVALID
	}

	return s.OK(&scyna_proto.StartTaskResponse{Id: taskID})
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
