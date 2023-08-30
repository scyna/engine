package trace

import (
	"log"

	scyna "github.com/scyna/core"
	scyna_proto "github.com/scyna/core/proto/generated"
)

func LogCreatedHandler(signal *scyna_proto.LogCreatedSignal) {
	log.Print(signal.Text)
	scyna.AddLog(scyna.LogData{
		ID:       signal.ID,
		Sequence: signal.SEQ,
		Level:    scyna.LogLevel(signal.Level),
		Message:  signal.Text,
	})
}
