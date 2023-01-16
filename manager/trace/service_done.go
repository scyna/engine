package trace

import (
	"github.com/scylladb/gocqlx/v2/qb"
	scyna "github.com/scyna/core"
	scyna_proto "github.com/scyna/core/proto/generated"
)

func ServiceDone(signal *scyna_proto.EndpointDoneSignal) {
	qb.Insert("scyna.tag").
		Columns("trace_id", "key", "value").
		Query(scyna.DB).
		Bind(signal.TraceID, "response", signal.Response).
		ExecRelease()
}
