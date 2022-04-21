package logging

import (
	"log"

	scyna "github.com/scyna/go"

	"google.golang.org/protobuf/proto"
)

func Write(data []byte) {
	var signal scyna.WriteLogSignal
	if err := proto.Unmarshal(data, &signal); err != nil {
		scyna.LOG.Error("Can not parse log info")
		return
	}

	log.Print(signal.Text)

	scyna.AddLog(scyna.LogData{
		ID:       signal.Id,
		Sequence: signal.Seq,
		Level:    scyna.LogLevel(signal.Level),
		Message:  signal.Text,
		Session:  signal.Session,
	})
}
