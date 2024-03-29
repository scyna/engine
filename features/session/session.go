package session

import (
	"log"
	"time"

	scyna "github.com/scyna/core"
	scyna_const "github.com/scyna/core/const"
)

func Init(module string, secret string) {
	if id, err := newSession(module, secret); err != scyna.OK {
		panic("Error in create session")
	} else {
		scyna.Session = scyna.NewSession(id)
	}
}

func newSession(module string, secret string) (uint64, scyna.Error) {
	log.Print("Creating session for module: ", module)
	var secret_ string
	if err := scyna.DB.QueryOne("SELECT secret FROM "+scyna_const.MODULE_TABLE+
		" WHERE code = ?", module).Scan(&secret_); err != nil {
		log.Print("Module not existed: ", err.Error())
		return 0, scyna.MODULE_NOT_EXISTS
	}

	if secret != secret_ {
		log.Print("Module secret is not correct")
		return 0, scyna.PERMISSION_ERROR
	}

	sid := scyna.ID.Next()
	now := time.Now()

	if err := scyna.DB.Execute("INSERT INTO "+scyna_const.SESSION_TABLE+
		" (id, module, started, updated) VALUES (?, ?, ?, ?)",
		sid, module, now, now); err != nil {
		log.Print("Can not save session to database:", err)
		return 0, scyna.SERVER_ERROR
	}

	log.Print("Session Created")
	return sid, scyna.OK
}
