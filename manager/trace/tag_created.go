package trace

import (
	scyna "github.com/scyna/core"
	scyna_const "github.com/scyna/core/const"
	scyna_proto "github.com/scyna/core/proto/generated"
)

func TagCreated(signal *scyna_proto.TagCreatedSignal) {
	scyna.DB.Execute("INSERT INTO "+scyna_const.TAG_TABLE+
		"(trace_id, key, value) VALUES(?,?,?)", signal.TraceID, signal.Key, signal.Value)
	// qb.Insert(scyna_const.TAG_TABLE).
	// 	Columns("trace_id", "key", "value").
	// 	Query(scyna.DB).
	// 	Bind(signal.TraceID, signal.Key, signal.Value).
	// 	ExecRelease()
}
