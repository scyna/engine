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

func TraceCreatedHandler(signal *scyna_proto.TraceCreatedSignal) {
	day := scyna_utils.GetDayByTime(time.Now())

	if err := scyna.DB.Execute("INSERT INTO "+scyna_const.TRACE_TABLE+
		"(type, path, day, id, time, duration, session, status, source)"+" VALUES (?,?,?,?,?,?,?,?,?)",
		signal.Type, signal.Path, day, signal.ID, time.UnixMicro(int64(signal.Time)),
		signal.Duration, signal.SessionID, signal.Status, signal.Source); err != nil {
		log.Print("Can not save trace created " + strconv.FormatUint(signal.ID, 10) + " / " + err.Error())
	}
}
