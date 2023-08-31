package proxy

import (
	"fmt"

	"github.com/nats-io/nats.go"
	scyna "github.com/scyna/core"
	scyna_const "github.com/scyna/core/const"
)

type Client struct {
	ID     string `db:"id"`
	Secret string `db:"secret"`
}

func (proxy *Proxy) initClients() {
	proxy.Clients = proxy.loadClients()
	_, err := scyna.Nats.Subscribe(scyna_const.CLIENT_UPDATE_CHANNEL, func(msg *nats.Msg) {
		scyna.Session.Info("Reload Clients")
		proxy.Clients = proxy.loadClients()
	})
	if err != nil {
		fmt.Println("initClients: " + err.Error())
	}
}

func (proxy *Proxy) loadClients() map[string]Client {
	ret := make(map[string]Client)

	scanner := scyna.DB.QueryMany("SELECT id, secret FROM client")
	for scanner.Next() {
		var client Client
		if err := scanner.Scan(&client.ID, &client.Secret); err != nil {
			scyna.Session.Error("Load Clients fail: " + err.Error())
		}
		ret[client.ID] = client
	}

	return ret
}
