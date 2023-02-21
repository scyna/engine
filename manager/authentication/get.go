package authentication

import (
	"log"
	"time"

	"github.com/scylladb/gocqlx/v2/qb"
	scyna "github.com/scyna/core"
	scyna_proto "github.com/scyna/core/proto/generated"
)

func Get(s *scyna.Endpoint, request *scyna_proto.GetAuthRequest) scyna.Error {
	log.Println("Receive GetAuthRequest")
	if expired, userID := getAuthentication(request.Token, request.App); expired != nil {
		return s.OK(&scyna_proto.GetAuthResponse{
			Token:   request.Token,
			UserID:  userID,
			Expired: uint64(expired.UnixMicro()),
		})
	} else {
		s.Warning("Not exists Token, App")
		return scyna.REQUEST_INVALID
	}
}

func getAuthentication(token string, app string) (*time.Time, string) {
	/*check authentication*/
	var auth struct {
		Expired time.Time `db:"expired"`
		Apps    []string  `db:"apps"`
		UserID  string    `db:"uid"`
	}

	if err := qb.Select("scyna.authentication").
		Columns("expired", "apps", "uid").
		Where(qb.Eq("id")).
		Limit(1).
		Query(scyna.DB).
		Bind(token).
		GetRelease(&auth); err != nil {
		log.Println("authentication", err.Error())
		return nil, ""
	}

	hasApp := false
	for _, a := range auth.Apps {
		if a == app {
			hasApp = true
			break
		}
	}

	if !hasApp {
		log.Print("No app")
		return nil, ""
	}

	return &auth.Expired, auth.UserID
}
