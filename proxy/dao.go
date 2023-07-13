package proxy

import (
	"time"

	"github.com/gocql/gocql"
	scyna "github.com/scyna/core"
	scyna_const "github.com/scyna/core/const"
	scyna_utils "github.com/scyna/core/utils"
)

func saveTrace(trace scyna.Trace) {

	day := scyna_utils.GetDayByTime(time.Now())
	trace.Duration = uint64(time.Now().UnixNano() - trace.Time.UnixNano())
	qBatch := scyna.DB.NewBatch(gocql.LoggedBatch)
	qBatch.Query("INSERT INTO "+scyna_const.TRACE_TABLE+"(type, path, day, id, time, duration, session_id, source, status) VALUES (?,?,?,?,?,?,?,?,?)",
		trace.Type,
		trace.Path,
		day,
		trace.ID,
		trace.Time,
		trace.Duration,
		trace.SessionID,
		trace.Source,
		trace.Status,
	)
	qBatch.Query("INSERT INTO "+scyna_const.CLIENT_HAS_TRACE_TABLE+"(client_id, trace_id, day) VALUES (?,?,?)",
		trace.Source,
		trace.ID,
		day,
	)
	if err := scyna.DB.ExecuteBatch(qBatch); err != nil {
		scyna.Session.Error("Can not save trace - " + err.Error())
	}
}
