package gateway

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/nats-io/nats.go"
	"github.com/scylladb/gocqlx/v2/qb"
	scyna "github.com/scyna/core"
	scyna_const "github.com/scyna/core/const"
	scyna_proto "github.com/scyna/core/proto/generated"
	"google.golang.org/protobuf/proto"
)

const PUBLIC_ENDPOINT_UPDATE_CHANNEL = scyna_const.KEYSPACE + ".public_endpoint.update"
const PUBLIC_ENDPOINT_TABLE = scyna_const.KEYSPACE + ".public_endpoint"
const ADD_PUBLIC_ENDPOINT_URL = scyna_const.BASEPATH + "/gateway/public-endpoint/add"
const REMOVE_PUBLIC_ENDPOINT_URL = scyna_const.BASEPATH + "/gateway/public-endpoint/remove"

func (gateway *Gateway) initPublicEndpoints() {
	gateway.PublicEndpoints = gateway.loadPublicEndPoints()
	_, err := scyna.Connection.Subscribe(PUBLIC_ENDPOINT_UPDATE_CHANNEL, func(msg *nats.Msg) {
		scyna.Session.Info("Reload Publics Endpoints")
		gateway.PublicEndpoints = gateway.loadPublicEndPoints()
	})
	if err != nil {
		fmt.Println("initPublicEndpoints: " + err.Error())
	}
}

func (gateway *Gateway) loadPublicEndPoints() []string {
	var ret []string
	if err := qb.Select(PUBLIC_ENDPOINT_TABLE).
		Columns("url").
		Query(scyna.DB).
		SelectRelease(&ret); err == nil {
	} else {
		scyna.Session.Error("Load Public Endpoints fail: " + err.Error())
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

type RequestPublicEndpoint struct {
	Url string `json:"url"`
}

func AddPublicEndpoint(w http.ResponseWriter, r *http.Request) {
	log.Println("Receive AddPublicEndpoint")

	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	var request RequestPublicEndpoint
	if err := json.Unmarshal(buf, &request); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if err := qb.Insert(PUBLIC_ENDPOINT_TABLE).
		Columns("url").
		Query(scyna.DB).
		Bind(request.Url).
		ExecRelease(); err != nil {
		log.Println(err)
		http.Error(w, "Server Error", 400)
	}
	scyna.Connection.Publish(PUBLIC_ENDPOINT_UPDATE_CHANNEL, nil)
	response := scyna_proto.Response{
		Code: 200,
	}
	if data, err := proto.Marshal(&response); err == nil {
		w.WriteHeader(200)
		_, err = bytes.NewBuffer(data).WriteTo(w)
		if err != nil {
			log.Println("Proxy write data error: " + err.Error())
		}
	} else {
		http.Error(w, "Server Error", 400)
	}
}

func RemovePublicEndpoint(w http.ResponseWriter, r *http.Request) {
	log.Println("Receive RemovePublicEndpoint")

	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	var request RequestPublicEndpoint
	if err := json.Unmarshal(buf, &request); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if err := qb.Delete(PUBLIC_ENDPOINT_TABLE).
		Where(qb.Eq("url")).
		Query(scyna.DB).
		Bind(request.Url).
		ExecRelease(); err != nil {
		log.Println(err)
		http.Error(w, "Server Error", 400)
	}

	scyna.Connection.Publish(PUBLIC_ENDPOINT_UPDATE_CHANNEL, nil)

	response := scyna_proto.Response{
		Code: 200,
	}
	if data, err := proto.Marshal(&response); err == nil {
		w.WriteHeader(200)
		_, err = bytes.NewBuffer(data).WriteTo(w)
		if err != nil {
			log.Println("Proxy write data error: " + err.Error())
		}
	} else {
		http.Error(w, "Server Error", 400)
	}

}
