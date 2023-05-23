package session

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/scylladb/gocqlx/v2/qb"
	scyna "github.com/scyna/core"
	scyna_const "github.com/scyna/core/const"
	scyna_proto "github.com/scyna/core/proto/generated"
	"github.com/scyna/go/engine/manager/manager"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func Create(w http.ResponseWriter, r *http.Request) {
	log.Println("Receive CreateSessionRequest")
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	var request scyna_proto.CreateSessionRequest
	if err := proto.Unmarshal(buf, &request); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if request.Module == manager.MODULE_CODE {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if sid, err := newSession(request.Module, request.Secret); err == scyna.OK {
		var response scyna_proto.CreateSessionResponse
		response.SessionID = sid

		var value string
		if err := qb.Select(scyna_const.SETTING_TABLE).
			Columns("value").
			Where(qb.Eq("module"), qb.Eq("key")).
			Limit(1).
			Query(scyna.DB).
			Bind(request.Module, scyna_const.SETTING_KEY).
			GetRelease(&value); err != nil {
			log.Println("Can not find module config for module " + request.Module + " - " + err.Error())
		}

		if len(value) > 0 {
			var config scyna_proto.Configuration
			err := protojson.Unmarshal([]byte(value), &config)
			if err != nil {
				response.Config = manager.DefaultConfig
			} else {
				response.Config = &config
			}
		} else {
			response.Config = manager.DefaultConfig
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
		return
	}

	http.Error(w, "Unauthorized", http.StatusUnauthorized)
}
