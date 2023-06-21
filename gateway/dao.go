package gateway

import (
	"log"

	"time"

	"github.com/scylladb/gocqlx/v2/qb"
	scyna "github.com/scyna/core"
	scyna_const "github.com/scyna/core/const"
)

func updateSession(token string, exp time.Time) bool {
	err := qb.Update(scyna_const.AUTHENTICATION_TABLE).
		Set("expired").
		Where(qb.Eq("id")).
		Query(scyna.DB).
		Bind(exp, token).
		ExecRelease()
	return err == nil
}

func checkAuthentication(token string, app string, url string) *time.Time {
	/*check authentication*/
	var auth struct {
		Expired time.Time `db:"expired"`
		Apps    []string  `db:"apps"`
	}

	if err := qb.Select(scyna_const.AUTHENTICATION_TABLE).
		Columns("expired", "apps").
		Where(qb.Eq("id")).
		Limit(1).
		Query(scyna.DB).
		Bind(token).
		GetRelease(&auth); err != nil {
		log.Println("authentication", err.Error())
		return nil
	}

	hasApp := false
	for _, a := range auth.Apps {
		if a == app {
			hasApp = true
			break
		}
	}

	if !hasApp {
		log.Print("No app in auth" + app)
		return nil
	}

	/*check app_use_service*/
	if err := qb.Select(scyna_const.APP_USE_ENDPOINT_TABLE).
		Columns("application").
		Where(qb.Eq("application"), qb.Eq("url")).
		Limit(1).
		Query(scyna.DB).
		Bind(app, url).
		GetRelease(&app); err != nil {
		log.Println("application_use_endpoint", err.Error())
		return nil
	}
	ret := auth.Expired
	return &ret
}
