package trace

import (
	"github.com/scylladb/gocqlx/v2/qb"
	scyna "github.com/scyna/core"
	scyna_engine "github.com/scyna/core/engine"
)

func TagCreated(signal *scyna_engine.TagCreatedSignal) {
	qb.Insert("scyna.tag").
		Columns("trace_id", "key", "value").
		Query(scyna.DB).
		Bind(signal.TraceID, signal.Key, signal.Value).
		ExecRelease()
}
