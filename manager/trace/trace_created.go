package trace

import (
	"log"
	"strconv"
	"time"

	"github.com/gocql/gocql"
	scyna "github.com/scyna/core"

	scyna_const "github.com/scyna/core/const"
	scyna_proto "github.com/scyna/core/proto/generated"
	scyna_utils "github.com/scyna/core/utils"
)

func TraceCreated(signal *scyna_proto.TraceCreatedSignal) {
	day := scyna_utils.GetDayByTime(time.Now())
	var source *string = nil
	if len(signal.Source) > 0 {
		source = &signal.Source
	}

	if signal.ParentID == 0 {
		if err := scyna.DB.Execute("INSERT INTO "+scyna_const.TRACE_TABLE+
			"(type, path, day, id, time, duration, session_id, source, status)"+" VALUES (?,?,?,?,?,?,?,?,?)",
			signal.Type, signal.Path, day, signal.ID, time.UnixMicro(int64(signal.Time)),
			signal.Duration, signal.SessionID, source, signal.Status); err != nil {
			// if err := qb.Insert(scyna_const.TRACE_TABLE).
			// 	Columns("type", "path", "day", "id", "time", "duration", "session_id", "source", "status").
			// 	Query(scyna.DB).
			// 	Bind(
			// 		signal.Type,
			// 		signal.Path,
			// 		day,
			// 		signal.ID,
			// 		time.UnixMicro(int64(signal.Time)),
			// 		signal.Duration,
			// 		signal.SessionID,
			// 		source,
			// 		signal.Status).
			// 	ExecRelease(); err != nil {
			log.Print("Can not save trace created " + strconv.FormatUint(signal.ID, 10) + " / " + err.Error())
		}
	} else {
		qBatch := scyna.DB.Session.NewBatch(gocql.LoggedBatch)
		qBatch.Query("INSERT INTO "+scyna_const.TRACE_TABLE+"(type, path, day, id, time, duration, session_id, parent_id, source, status)"+
			" VALUES (?,?,?,?,?,?,?,?,?,?)",
			signal.Type,
			signal.Path,
			day,
			signal.ID,
			time.UnixMicro(int64(signal.Time)),
			signal.Duration,
			signal.SessionID,
			signal.ParentID,
			source,
			signal.Status)
		qBatch.Query("INSERT INTO "+scyna_const.SPAN_TABLE+"(parent_id, child_id) VALUES (?,?)", signal.ParentID, signal.ID)

		if err := scyna.DB.Session.ExecuteBatch(qBatch); err != nil {
			log.Print("Can not save trace created " + strconv.FormatUint(signal.ID, 10) + " / " + strconv.FormatUint(signal.ParentID, 10) + " / " + err.Error())
		}
	}

}
