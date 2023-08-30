package trace

import (
	"log"
	"strconv"
	"time"

	scyna "github.com/scyna/core"

	scyna_const "github.com/scyna/core/const"
	scyna_proto "github.com/scyna/core/proto/generated"
	scyna_utils "github.com/scyna/core/utils"
)

func TraceCreated(signal *scyna_proto.TraceCreatedSignal) {
	day := scyna_utils.GetDayByTime(time.Now())

	if err := scyna.DB.Execute("INSERT INTO "+scyna_const.TRACE_TABLE+
		"(type, path, day, id, time, duration, session_id)"+" VALUES (?,?,?,?,?,?,?,?)",
		signal.Type, signal.Path, day, signal.ID, time.UnixMicro(int64(signal.Time)),
		signal.Duration, signal.SessionID); err != nil {
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
}
