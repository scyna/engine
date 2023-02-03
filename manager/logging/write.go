package logging

import (
	"log"

	scyna "github.com/scyna/core"
	scyna_engine "github.com/scyna/core/engine"
)

func Write(signal *scyna_engine.LogCreatedSignal) {
	log.Print(signal.Text)
	scyna.AddLog(scyna.LogData{
		ID:       signal.ID,
		Sequence: signal.SEQ,
		Level:    scyna.LogLevel(signal.Level),
		Message:  signal.Text,
		Session:  signal.Session,
	})
}
