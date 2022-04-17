package proxy

import (
	"bytes"
	"log"
	"net/http"
	"strings"
	"time"

	scyna "github.com/scyna/go"

	"google.golang.org/protobuf/proto"
)

type Proxy struct {
	Queries QueryPool
	Clients map[string]Client
}

func NewProxy() *Proxy {
	ret := &Proxy{Queries: NewQueryPool()}
	ret.initClients()
	return ret
}

func (proxy *Proxy) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	callID := scyna.ID.Next()
	start := time.Now()
	day := scyna.GetDayByTime(start)

	query := proxy.Queries.GetQuery()
	defer proxy.Queries.Put(query)

	/*authenticate*/
	url := req.URL.String()
	clientID := req.Header.Get("Client-Id")
	clientSecret := req.Header.Get("Client-Secret")
	client, ok := proxy.Clients[clientID]
	contentType := req.Header.Get("Content-Type")
	//https://descynaper.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Type
	for _, data := range strings.Split(contentType, ";") {
		value := strings.TrimSpace(strings.Trim(data, ";"))
		if strings.HasPrefix(value, "application/") {
			contentType = value
			continue
		}
	}

	/*CORS*/
	rw.Header().Set("Content-Type", contentType)
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if !ok || clientSecret != client.Secret {
		http.Error(rw, "Unauthorized", http.StatusUnauthorized)
		log.Print("Wrong client id or secret: ", clientID)
		proxy.SaveErrorCall(clientID, 401, callID, day, start, req.URL.Path)
		return
	}

	// if client.State != uint32(scyna.ClientState_ACTIVE) {
	// 	http.Error(rw, "Unauthorized", http.StatusUnauthorized)
	// 	log.Printf("Client is inactive: %s\n", clientID)
	// 	proxy.SaveErrorCall(clientID, 401, callID, day, start, req.URL.Path)
	// 	return
	// }

	if err := query.Authenticate.Bind(clientID, url).Get(&url); err != nil {
		http.Error(rw, "Unauthorized", http.StatusUnauthorized)
		log.Printf("Wrong url: %s, error = %s\n", url, err.Error())
		proxy.SaveErrorCall(clientID, 401, callID, day, start, req.URL.Path)
		return
	}

	service := scyna.Services.GetService()
	defer scyna.Services.PutService(service)

	if contentType == "application/json" {
		service.Request.JSON = true
	} else if contentType == "application/protobuf" {
		service.Request.JSON = false
	} else {
		http.Error(rw, "Content-Type must be JSON or PROTOBUF ", http.StatusNotAcceptable)
		proxy.SaveErrorCall(clientID, 401, callID, day, start, url)
		return
	}

	/*build request*/
	err := service.Request.Build(req)
	if err != nil {
		http.Error(rw, "Cannot process request", http.StatusInternalServerError)
		proxy.SaveErrorCall(clientID, 500, callID, day, start, url)
		return
	}

	service.Request.CallID = callID
	service.Request.Data = client.Type

	/*serialize the request */
	reqBytes, err := proto.Marshal(&service.Request)
	if err != nil {
		http.Error(rw, "Cannot process request", http.StatusInternalServerError)
		proxy.SaveErrorCall(clientID, 500, callID, day, start, url)
		return
	}

	/*post request to message queue*/
	msg, respErr := scyna.Connection.Request(scyna.PublishURL(url), reqBytes, 10*time.Second)
	if respErr != nil {
		http.Error(rw, "No response", http.StatusInternalServerError)
		log.Println("ServeHTTP: Nats: " + respErr.Error())
		proxy.SaveErrorCall(clientID, 500, callID, day, start, url)
		return
	}

	/*response*/
	err = service.Response.ReadFrom(msg.Data)
	if err != nil {
		log.Println("nats-proxy:" + err.Error())
		http.Error(rw, "Cannot deserialize response", http.StatusInternalServerError)
		proxy.SaveErrorCall(clientID, 500, callID, day, start, url)
		return
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
	proxy.SaveCall(clientID, callID, day, start, duration, req.URL.Path, service)
}
