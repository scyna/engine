package session

import (
	"log"
	"time"

	"github.com/scylladb/gocqlx/v2/qb"
	scyna "github.com/scyna/go/scyna"
	"google.golang.org/protobuf/proto"
)

func End(data []byte) {
	var signal scyna.EndSessionSignal
	if err := proto.Unmarshal(data, &signal); err != nil {
		scyna.LOG.Error("Can not parse EndSessionSignal")
		return
	}

	if applied, err := qb.Update("scyna.session").
		Set("end", "exit_code").
		Where(qb.Eq("id")).Existing().
		Query(scyna.DB).
		Bind(time.Now(), signal.Code, signal.ID).
		ExecCASRelease(); !applied {
		if err != nil {
			log.Print("Can not update EndSessionSignal:", err.Error())
		}
	}
}
