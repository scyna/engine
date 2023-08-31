package proxy

import (
	"time"

	scyna "github.com/scyna/core"
	scyna_const "github.com/scyna/core/const"
	scyna_utils "github.com/scyna/core/utils"
)

func saveTrace(trace *trace) {
	day := scyna_utils.GetDayByTime(time.Now())
	trace.Duration = uint64(time.Now().UnixNano() - trace.Time.UnixNano())

	if err := scyna.DB.Execute("INSERT INTO "+scyna_const.TRACE_TABLE+
		"(type, path, day, id, time, duration, session, status, source) VALUES (?,?,?,?,?,?,?,?,?)",
		trace.Type, trace.Path, day, trace.ID, trace.Time, trace.Duration, trace.SessionID, trace.Status, trace.Source); err != nil {
		scyna.Session.Error("Can not save trace - " + err.Error())
	}
}
