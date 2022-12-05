package trace

import (
	"github.com/scylladb/gocqlx/v2/qb"
	scyna "github.com/scyna/core"
)

func ServiceDone(signal *scyna.EnpointDoneSignal) {
	qb.Insert("scyna.tag").
		Columns("trace_id", "key", "value").
		Query(scyna.DB).
		Bind(signal.TraceID, "response", signal.Response).
		ExecRelease()
}
