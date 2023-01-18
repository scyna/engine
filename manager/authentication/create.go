package authentication

import (
	"log"
	"time"

	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v2/qb"
	scyna "github.com/scyna/core"
	scyna_proto "github.com/scyna/core/proto/generated"
)

var serialNumber = scyna.InitSerialNumber("scyna.auth")

func Create(s *scyna.Endpoint, request *scyna_proto.CreateAuthRequest) scyna.Error {
	log.Println("Receive CreateAuthRequest")

	if !checkOrg(request.Organization, request.Secret) {
		s.Logger.Warning("Organization not exist")
		return scyna.REQUEST_INVALID
	}

	if len(request.Apps) == 0 {
		return scyna.REQUEST_INVALID
	}

	for _, app := range request.Apps {
		if !checkApp(app) {
			scyna.LOG.Warning("App not exist: " + app)
			return scyna.REQUEST_INVALID
		}
	}

	id := serialNumber.Next()
	if err := createAuth(id, request.Apps, request.UserID); err != scyna.OK {
		return err
	}

	now := time.Now()
	s.Response(&scyna_proto.CreateAuthResponse{
		Token:   id,
		Expired: uint64(now.Add(time.Hour * 8).UnixMicro()),
	})
	return scyna.OK
}

func createAuth(id string, apps []string, userID string) scyna.Error {
	now := time.Now()
	exp := now.Add(time.Hour * 8)

	session := scyna.DB.Session
	batch := session.NewBatch(gocql.LoggedBatch)
	batch.Query("INSERT INTO scyna.authentication (id, apps, expired, time, uid) VALUES (?,?,?,?,?);",
		id, apps, exp, now, userID)
	for _, app := range apps {
		batch.Query("INSERT INTO scyna.app_has_auth (auth_id, app_code, user_id) VALUES (?,?,?);",
			id, app, userID)
	}
	if err := session.ExecuteBatch(batch); err != nil {
		log.Print("Error:", err)
		return scyna.SERVER_ERROR
	}
	return scyna.OK
}

func checkOrg(code string, secret string) bool {

	var secret_ string
	if err := qb.Select("scyna.domain").
		Columns("password"). //FIXME: change to use secret
		Where(qb.Eq("code")).
		Query(scyna.DB).
		Bind(code).
		GetRelease(&secret_); err != nil {
		log.Println("Check OrgCode", err.Error())
		return false
	}

	if secret != secret_ {
		return false
	}

	return true
}

func checkApp(code string) bool {
	if err := qb.Select("scyna.application").
		Columns("code").
		Where(qb.Eq("code")).
		Query(scyna.DB).
		Bind(code).
		GetRelease(&code); err != nil {
		return false
	}
	return true
}
