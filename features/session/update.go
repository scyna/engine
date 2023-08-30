package session

import (
	"log"
	"time"

	scyna "github.com/scyna/core"
	scyna_const "github.com/scyna/core/const"
	scyna_proto "github.com/scyna/core/proto/generated"
)

func Update(signal *scyna_proto.UpdateSessionSignal) {
	if err := scyna.DB.Execute("UPDATE "+scyna_const.SESSION_TABLE+
		" SET updated = ? WHERE id = ? AND module = ? IF EXISTS",
		time.Now(), signal.ID, signal.Module); err != nil {
		log.Print("Can not update UpdateSessionSignal: ", err.Error())
	}
}
