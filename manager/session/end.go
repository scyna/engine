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
		" SET ended = ?, exit_code = ? WHERE id = ? AND module = ? IF EXISTS",
		time.Now(), signal.Code, signal.ID, signal.Module); err != nil {
		log.Print("Can not update EndSessionSignal:", err.Error())
	}
}
