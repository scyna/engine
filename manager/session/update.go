package session

import (
	"log"
	"time"

	"github.com/scylladb/gocqlx/v2/qb"
	scyna "github.com/scyna/core"
	scyna_engine "github.com/scyna/core/engine"
)

func Update(signal *scyna_engine.UpdateSessionSignal) {
	if applied, err := qb.Update("scyna.session").
		Set("last_update").
		Where(qb.Eq("id"), qb.Eq("module")).
		Existing().
		Query(scyna.DB).
		Bind(time.Now(), signal.ID, signal.Module).
		ExecCASRelease(); !applied {
		if err != nil {
			log.Print("Can not update UpdateSessionSignal: ", err.Error())
		}
	}
}
