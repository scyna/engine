package authentication

import (
	"log"
	"time"

	"github.com/scylladb/gocqlx/v2/qb"
	scyna "github.com/scyna/core"
	scyna_const "github.com/scyna/core/const"
	scyna_proto "github.com/scyna/core/proto/generated"
)

func Logout(ctx *scyna.Endpoint, request *scyna_proto.LogoutRequest) scyna.Error {
	log.Println("Receive LogoutRequest")

	if err := updateSession(request.Token, request.UserID); err != scyna.OK {
		return err
	}

	return scyna.OK
}

func updateSession(token string, userID string) scyna.Error {
	var userID_ string
	err := qb.Select(scyna_const.AUTHENTICATION_TABLE).
		Columns("uid").
		Where(qb.Eq("id")).
		Limit(1).
		Query(scyna.DB).
		Bind(token).
		GetRelease(&userID_)
	if err != nil {
		log.Print("Error:", err)
		return scyna.SERVER_ERROR
	}

	if userID != userID_ {
		return scyna.SERVER_ERROR
	}

	now := time.Now()

	err = qb.Update(scyna_const.AUTHENTICATION_TABLE).
		Set("expired").
		Where(qb.Eq("id")).
		Query(scyna.DB).
		Bind(now, token).
		ExecRelease()

	if err != nil {
		log.Print("Error:", err)
		return scyna.SERVER_ERROR
	}
	return nil
}
