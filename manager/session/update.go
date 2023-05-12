package session

import (
	"log"
	"time"

	"github.com/scylladb/gocqlx/v2/qb"
	scyna "github.com/scyna/core"
	scyna_const "github.com/scyna/core/const"
	scyna_proto "github.com/scyna/core/proto/generated"
)

func Update(signal *scyna_proto.UpdateSessionSignal) {
	if applied, err := qb.Update(scyna_const.SESSION_TABLE).
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
