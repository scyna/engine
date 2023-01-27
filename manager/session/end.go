package session

import (
	"log"
	"time"

	"github.com/scylladb/gocqlx/v2/qb"
	scyna "github.com/scyna/core"
	scyna_proto "github.com/scyna/core/proto/generated"
)

//https://tldp.org/LDP/abs/html/exitcodes.html

func End(signal *scyna_proto.EndSessionSignal) {
	if applied, err := qb.Update("scyna.session").
		Set("end", "exit_code").
		Where(qb.Eq("id"), qb.Eq("module")).Existing().
		Query(scyna.DB).
		Bind(time.Now(), signal.Code, signal.ID, signal.Module).
		ExecCASRelease(); !applied {
		if err != nil {
			log.Print("Can not update EndSessionSignal:", err.Error())
		}
	}
}
