package session

import (
	"log"
	"time"

	scyna "github.com/scyna/core"
	scyna_const "github.com/scyna/core/const"
	scyna_proto "github.com/scyna/core/proto/generated"
)

//https://tldp.org/LDP/abs/html/exitcodes.html

func End(signal *scyna_proto.EndSessionSignal) {

	if err := scyna.DB.Execute("UPDATE "+scyna_const.SESSION_TABLE+
		" SET end = ?, exit_code = ? WHERE id = ? AND module = ?",
		time.Now(), signal.Code, signal.ID, signal.Module); err != nil {
		// if applied, err := qb.Update(scyna_const.SESSION_TABLE).
		// 	Set("end", "exit_code").
		// 	Where(qb.Eq("id"), qb.Eq("module")).Existing().
		// 	Query(scyna.DB).
		// 	Bind(time.Now(), signal.Code, signal.ID, signal.Module).
		// 	ExecCASRelease(); !applied {
		log.Print("Can not update EndSessionSignal:", err.Error())

	}
}
