package logging

import (
	"log"

	"github.com/scyna/go/scyna"
)

func Write(signal *scyna.WriteLogSignal) {
	log.Print(signal.Text)
	scyna.AddLog(scyna.LogData{
		ID:       signal.Id,
		Sequence: signal.Seq,
		Level:    scyna.LogLevel(signal.Level),
		Message:  signal.Text,
		Session:  signal.Session,
	})
}
