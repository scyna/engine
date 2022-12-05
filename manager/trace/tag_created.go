package trace

import (
	"github.com/scylladb/gocqlx/v2/qb"
	scyna "github.com/scyna/core"
)

func TagCreated(signal *scyna.TagCreatedSignal) {
	qb.Insert("scyna.tag").
		Columns("trace_id", "key", "value").
		Query(scyna.DB).
		Bind(signal.TraceID, signal.Key, signal.Value).
		ExecRelease()
}
