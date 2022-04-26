package session

import (
	"log"
	"time"

	"github.com/scylladb/gocqlx/v2/qb"
	scyna "github.com/scyna/go/scyna"
)

func End(signal *scyna.EndSessionSignal) {
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
