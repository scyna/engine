package authentication

import (
	"log"
	"time"

	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v2/qb"
	scyna "github.com/scyna/core"
	scyna_const "github.com/scyna/core/const"
	scyna_proto "github.com/scyna/core/proto/generated"
)

var serialNumber = scyna.InitSerialNumber("scyna.auth")

func Create(s *scyna.Endpoint, request *scyna_proto.CreateAuthRequest) scyna.Error {
	log.Println("Receive CreateAuthRequest")

	if len(request.Apps) == 0 {
		return scyna.REQUEST_INVALID
	}

	for _, app := range request.Apps {
		if !checkApp(app) {
			scyna.Session.Warning("App not exist: " + app)
			return scyna.REQUEST_INVALID
		}
	}

	id := serialNumber.Next()
	if err := createAuth(id, request.Apps, request.UID); err != scyna.OK {
		return err
	}

	now := time.Now()

	return s.OK(&scyna_proto.CreateAuthResponse{Token: id, Expired: uint64(now.Add(time.Hour * 8).UnixMicro())})
}

func createAuth(id string, apps []string, userID string) scyna.Error {
	now := time.Now()
	exp := now.Add(time.Hour * 8)

	session := scyna.DB.Session
	batch := session.NewBatch(gocql.LoggedBatch)
	batch.Query("INSERT INTO "+scyna_const.AUTHENTICATION_TABLE+" (id, apps, expired, time, uid) VALUES (?,?,?,?,?);",
		id, apps, exp, now, userID)
	for _, app := range apps {
		batch.Query("INSERT INTO "+scyna_const.APP_HAS_AUTH_TABLE+" (app, auth, uid) VALUES (?,?,?);",
			app, id, userID)
	}
	if err := session.ExecuteBatch(batch); err != nil {
		log.Print("Error:", err)
		return scyna.SERVER_ERROR
	}
	return scyna.OK
}

func checkApp(code string) bool {
	if err := qb.Select(scyna_const.APPLICATION_TABLE).
		Columns("code").
		Where(qb.Eq("code")).
		Query(scyna.DB).
		Bind(code).
		GetRelease(&code); err != nil {
		return false
	}
	return true
}
