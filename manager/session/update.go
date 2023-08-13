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
		" SET last_update = ? WHERE id = ? AND module = ? IF EXISTS",
		time.Now(), signal.ID, signal.Module); err != nil {
		// if applied, err := qb.Update(scyna_const.SESSION_TABLE).
		// 	Set("last_update").
		// 	Where(qb.Eq("id"), qb.Eq("module")).
		// 	Existing().
		// 	Query(scyna.DB).
		// 	Bind(time.Now(), signal.ID, signal.Module).
		// 	ExecCASRelease(); !applied {
		log.Print("Can not update UpdateSessionSignal: ", err.Error())

	}
}
