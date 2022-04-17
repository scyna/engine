package gateway

import (
	"bytes"
	"log"
	"net/http"
	"time"

	scyna "github.com/scyna/go"

	"google.golang.org/protobuf/proto"
)

type Gateway struct {
	Queries      QueryPool
	Applications map[string]Application
}

func NewGateway() *Gateway {
	ret := &Gateway{Queries: NewQueryPool()}
	ret.initApplications()
	return ret
}

func (gateway *Gateway) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	var app Application
	callID := scyna.ID.Next()
	start := time.Now()
	day := scyna.GetDayByTime(start)
	auth := false
	ok, appID, json, url := parseUrl(req.URL.String())

	if !ok {
		log.Print("Path not ok")
		http.Error(rw, "NotFound", http.StatusNotFound)
		return
	}

	log.Print(url)

	query := gateway.Queries.GetQuery()
	defer gateway.Queries.Put(query)

	service := scyna.Services.GetService()
	defer scyna.Services.Put(service)

	/*headers*/
	rw.Header().Set("Access-Control-Allow-Origin", req.Header.Get("Origin"))
	rw.Header().Set("Access-Control-Allow-Credentials", "true")
	rw.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	rw.Header().Set("Access-Control-Allow-Methods", "POST")
	if json {
		rw.Header().Set("Content-Type", "application/json")
	} else {
		rw.Header().Set("Content-Type", "application/octet-stream")
	}

	service.Request.JSON = json
	service.Request.CallID = callID

	if a, ok := gateway.Applications[appID]; !ok {
		http.Error(rw, "Forbidden", http.StatusForbidden)
		return
	} else {
		app = a
	}

	if url == "/auth" {
		auth = true
		url = app.AuthURL
		log.Print(url)
	} else {
		if cookie, err := req.Cookie("session"); err == nil {
			token := cookie.Value
			service.Request.Data = token
			if exp := checkService(token, appID, url); exp == nil {
				http.Error(rw, "Unauthorized", http.StatusUnauthorized)
				return
			} else {
				log.Print(exp)
				now := time.Now()
				if exp.Before(now) {
					log.Print("Session expired")
					http.Error(rw, "Unauthorized", http.StatusUnauthorized)
					return
				} else {
					if exp.After(now.Add(time.Minute * 10)) {
						/*auto extend expire*/
						if updateSesion(token, now.Add(time.Hour*8)) {
							cookie.Expires = now
							http.SetCookie(rw, cookie)
						}
					}
				}
			}
		} else {
			log.Print("Can not get cookie")
			http.Error(rw, "Unauthorized", http.StatusUnauthorized)
			return
		}
	}

	/*build request*/
	err := service.Request.Build(req)
	if err != nil {
		http.Error(rw, "Cannot process request", http.StatusInternalServerError)
		gateway.saveErrorCall(appID, 500, callID, day, start, url, "app")
		return
	}

	/*serialize the request */
	reqBytes, err := proto.Marshal(&service.Request)
	if err != nil {
		http.Error(rw, "Cannot process request", http.StatusInternalServerError)
		gateway.saveErrorCall(appID, 500, callID, day, start, url, "app")
		return
	}

	/*post request to message queue*/
	msg, respErr := scyna.Connection.Request(scyna.PublishURL(url), reqBytes, 10*time.Second)
	if respErr != nil {
		http.Error(rw, "No response", http.StatusInternalServerError)
		log.Println("ServeHTTP: Nats: " + respErr.Error())
		gateway.saveErrorCall(appID, 500, callID, day, start, url, "app")
		return
	}

	/*response*/
	err = service.Response.ReadFrom(msg.Data)
	if err != nil {
		log.Println("nats-proxy:" + err.Error())
		http.Error(rw, "Cannot deserialize response", http.StatusInternalServerError)
		gateway.saveErrorCall(appID, 500, callID, day, start, url, "app")
		return
	}

	if auth {
		if service.Response.Code == 200 {
			cookie := &http.Cookie{
				Name:     "session",
				Value:    service.Response.Token,
				Path:     "/",
				Expires:  time.Unix(0, int64(service.Response.Expired*1000)),
				HttpOnly: true,
				SameSite: http.SameSiteNoneMode,
				Secure:   true,
			}
			http.SetCookie(rw, cookie)
			log.Print("Set cookie:", service.Response.Token)
		} else {
			c := &http.Cookie{
				Name:     "session",
				Value:    "",
				Path:     "/",
				Expires:  time.Unix(0, 0),
				HttpOnly: true,
				SameSite: http.SameSiteNoneMode,
				Secure:   true,
			}
			http.SetCookie(rw, c) /*clear cookie*/
			/*TODO: make Authentication inactive here or delete from database*/
		}
	}

	rw.WriteHeader(int(service.Response.Code))
	_, err = bytes.NewBuffer(service.Response.Body).WriteTo(rw)
	if err != nil {
		log.Println("Proxy write data error: " + err.Error())
	}

	if f, ok := rw.(http.Flusher); ok {
		f.Flush()
	}

	duration := time.Now().UnixMicro() - start.UnixMicro()
	gateway.saveCall(appID, callID, day, start, duration, url, "app", service)
}
