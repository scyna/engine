package trace

import (
	"log"
	"time"

	scyna "github.com/scyna/core"
	scyna_const "github.com/scyna/core/const"
	scyna_proto "github.com/scyna/core/proto/generated"
	scyna_utils "github.com/scyna/core/utils"
)

func EndpointDoneHandler(signal *scyna_proto.EndpointDoneSignal) {
	now := time.Now()
	if err := scyna.DB.Execute("INSERT INTO "+scyna_const.ENDPOINT_TRACE_TABLE+
		"(day, time, id, request, response, status, session) VALUES (?,?,?,?,?,?,?)",
		scyna_utils.GetDayByTime(now), now, signal.TraceID, signal.Request, signal.Response, signal.Status, signal.SessionID); err != nil {
		log.Print("Can not save endpoint done" + err.Error())
	}
}
