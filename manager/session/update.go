package session

import (
	"github.com/scylladb/gocqlx/v2/qb"
	scyna "github.com/scyna/go"
	"google.golang.org/protobuf/proto"
	"log"
	"time"
)

func Update(data []byte) {
	var signal scyna.UpdateSessionSignal
	if err := proto.Unmarshal(data, &signal); err != nil {
		scyna.LOG.Error("Can not parse UpdateSessionSignal")
		return
	}

	if applied, err := qb.Update("scyna.session").
		Set("last_update").
		Where(qb.Eq("id")).
		Existing().
		Query(scyna.DB).
		Bind(time.Now(), signal.ID).
		ExecCASRelease(); !applied {
		if err != nil {
			log.Print("Can not update UpdateSessionSignal: ", err.Error())
		}
	}
}
