package trace

import (
	"github.com/gocql/gocql"
	scyna "github.com/scyna/core"
	scyna_proto "github.com/scyna/core/proto/generated"
)

func ServiceDone(signal *scyna_proto.EndpointDoneSignal) {
	batch := scyna.DB.NewBatch(gocql.UnloggedBatch)
	batch.Query("INSERT INTO scyna.tag(trace_id, key, value) VALUES(?,?,?)", signal.TraceID, "request", signal.Request)
	batch.Query("INSERT INTO scyna.tag(trace_id, key, value) VALUES(?,?,?)", signal.TraceID, "response", signal.Response)
	scyna.DB.ExecuteBatch(batch)
}
