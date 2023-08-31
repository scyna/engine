package gateway

import (
	"fmt"
	"log"

	"github.com/nats-io/nats.go"
	scyna "github.com/scyna/core"
	scyna_const "github.com/scyna/core/const"
)

const PUBLIC_ENDPOINT_UPDATE_CHANNEL = scyna_const.KEYSPACE + ".public_endpoint.update"
const PUBLIC_ENDPOINT_TABLE = scyna_const.KEYSPACE + ".public_endpoint"
const ADD_PUBLIC_ENDPOINT_URL = scyna_const.BASEPATH + "/gateway/public-endpoint/add"
const REMOVE_PUBLIC_ENDPOINT_URL = scyna_const.BASEPATH + "/gateway/public-endpoint/remove"

func (gateway *Gateway) initPublicEndpoints() {
	gateway.PublicEndpoints = gateway.loadPublicEndPoints()
	_, err := scyna.Nats.Subscribe(PUBLIC_ENDPOINT_UPDATE_CHANNEL, func(msg *nats.Msg) {
		scyna.Session.Info("Reload Publics Endpoints")
		gateway.PublicEndpoints = gateway.loadPublicEndPoints()
	})
	if err != nil {
		fmt.Println("initPublicEndpoints: " + err.Error())
	}
}

func (gateway *Gateway) loadPublicEndPoints() []string {
	var ret []string

	scanners := scyna.DB.QueryMany("SELECT url FROM " + PUBLIC_ENDPOINT_TABLE)

	for scanners.Next() {
		var url string
		if err := scanners.Scan(&url); err != nil {
			scyna.Session.Error("Load Public Endpoints fail: " + err.Error())
		} else {
			ret = append(ret, url)
		}
	}

	return ret
}

func (gateway *Gateway) isPublicEndpoint(url string) bool {
	for _, publicEndpoint := range gateway.PublicEndpoints {
		if publicEndpoint == url {
			return true
		}
	}
	return false
}

func AddPublicEndpoint(ctx *scyna.Endpoint, request *AddPublicEndpointRequest) scyna.Error {
	log.Println("Receive AddPublicEndpoint")

	if err := scyna.DB.Execute("INSERT INTO "+PUBLIC_ENDPOINT_TABLE+" (url) VALUES (?)", request.Url); err != nil {
		// if err := qb.Insert(PUBLIC_ENDPOINT_TABLE).
		// 	Columns("url").
		// 	Query(scyna.DB).
		// 	Bind(request.Url).
		// 	ExecRelease(); err != nil {
		log.Println(err)
		return scyna.SERVER_ERROR
	}

	scyna.Nats.Publish(PUBLIC_ENDPOINT_UPDATE_CHANNEL, nil)
	return scyna.OK
}

func RemovePublicEndpoint(ctx *scyna.Endpoint, request *AddPublicEndpointRequest) scyna.Error {
	log.Println("Receive RemovePublicEndpoint")

	if err := scyna.DB.Execute("DELETE FROM "+PUBLIC_ENDPOINT_TABLE+" WHERE url = ?", request.Url); err != nil {
		// if err := qb.Delete(PUBLIC_ENDPOINT_TABLE).
		// 	Where(qb.Eq("url")).
		// 	Query(scyna.DB).
		// 	Bind(request.Url).
		// 	ExecRelease(); err != nil {
		log.Println(err)
		return scyna.SERVER_ERROR
	}

	scyna.Nats.Publish(PUBLIC_ENDPOINT_UPDATE_CHANNEL, nil)
	return scyna.OK
}
