package call

import (
	"log"
	"time"

	"github.com/gocql/gocql"
	"github.com/scyna/go/scyna"
)

func Write(signal *scyna.WriteCallSignal) {
	qBatch := scyna.DB.NewBatch(gocql.LoggedBatch)
	qBatch.Query("INSERT INTO scyna.call(id, day, time, duration, request, response, source, status, session_id, caller_id)"+
		" VALUES (?,?,?,?,?,?,?,?,?,?)",
		signal.Id,
		signal.Day,
		time.UnixMicro(int64(signal.Time)),
		signal.Duration,
		signal.Request,
		signal.Response,
		signal.Source,
		signal.Status,
		signal.SessionId,
		signal.CallerId)
	if len(signal.CallerId) > 0 {
		qBatch.Query("INSERT INTO scyna.client_has_call(client_id, call_id, day) VALUES (?,?,?)",
			signal.CallerId, signal.Id, signal.Day)
	}
	if err := scyna.DB.ExecuteBatch(qBatch); err != nil {
		log.Print(err)
		scyna.LOG.Error("Can not save call: " + err.Error())
	}
}
