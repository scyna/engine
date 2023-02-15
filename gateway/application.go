package gateway

import (
	"fmt"
	"log"

	"github.com/nats-io/nats.go"
	"github.com/scylladb/gocqlx/v2/qb"
	scyna "github.com/scyna/core"
	scyna_proto "github.com/scyna/core/proto/generated"
)

type Application struct {
	Code    string `db:"code"`
	AuthURL string `db:"auth_url"`
}

func (g *Gateway) initApplications() {
	g.Applications = loadApplications()
	_, err := scyna.Connection.Subscribe(scyna_proto.APP_UPDATE_CHANNEL, func(msg *nats.Msg) {
		scyna.LOG.Info("Reload Application")
		g.Applications = loadApplications()
	})
	if err != nil {
		fmt.Println("initClients: " + err.Error())
	}
}

func loadApplications() map[string]Application {
	ret := make(map[string]Application)
	var apps []Application

	if err := qb.Select("scyna.application").
		Columns("code", "auth_url").
		Query(scyna.DB).
		SelectRelease(&apps); err == nil {
		for _, c := range apps {
			log.Print(c.Code, "/", c.AuthURL)
			ret[c.Code] = c
		}
	} else {
		scyna.LOG.Error("Load Clients fail: " + err.Error())
	}
	return ret
}
