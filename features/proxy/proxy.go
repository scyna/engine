package proxy

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
	"time"

	scyna "github.com/scyna/core"
	scyna_const "github.com/scyna/core/const"
	scyna_utils "github.com/scyna/core/utils"
	"google.golang.org/protobuf/proto"
)

type trace struct {
	ParentID  uint64
	ID        uint64
	Type      scyna.TraceType
	Time      time.Time
	Duration  uint64
	Path      string
	SessionID uint64
	Status    uint32
	Source    string
}

type Proxy struct {
	Clients  map[string]Client
	Contexts scyna_utils.HttpContextPool
}

func NewProxy() *Proxy {
	ret := &Proxy{Contexts: scyna_utils.NewContextPool()}
	ret.initClients()
	return ret
}

func (proxy *Proxy) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	callID := scyna.ID.Next()

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

	trace := trace{
		ID:        callID,
		ParentID:  0,
		Time:      time.Now(),
		Path:      url,
		Type:      scyna.TRACE_ENDPOINT,
		SessionID: scyna.Session.ID(),
		Source:    clientID,
	}
	defer saveTrace(&trace)

	ctx := proxy.Contexts.GetContext()
	defer proxy.Contexts.PutContext(ctx)

	if !ok || clientSecret != client.Secret {
		http.Error(rw, "Unauthorized", http.StatusUnauthorized)
		scyna.Session.Info("Wrong client id or secret: " + clientID + ", secret: " + clientSecret)
		trace.Status = http.StatusUnauthorized
		return
	}

	if err := scyna.DB.QueryOne("SELECT url FROM "+scyna_const.CLIENT_USE_ENDPOINT_TABLE+" WHERE client=? AND url=?", 
			clientID, url).Scan(&url); err != nil {
		http.Error(rw, "Unauthorized", http.StatusUnauthorized)
		scyna.Session.Info(fmt.Sprintf("Wrong url: %s, error = %s\n", url, err.Error()))
		trace.Status = http.StatusUnauthorized
		return
	}

	if contentType == "application/json" {
		ctx.Request.JSON = true
	} else if contentType == "application/protobuf" {
		ctx.Request.JSON = false
	} else {
		http.Error(rw, "Content-Type must be JSON or PROTOBUF ", http.StatusNotAcceptable)
		trace.Status = http.StatusNotAcceptable
		return
	}

	/*build request*/
	err := ctx.Request.Build(req)
	if err != nil {
		http.Error(rw, "Cannot process request", http.StatusInternalServerError)
		trace.Status = http.StatusInternalServerError
		return
	}
	ctx.Request.TraceID = callID

	/*serialize the request */
	reqBytes, err := proto.Marshal(&ctx.Request)
	if err != nil {
		http.Error(rw, "Cannot process request", http.StatusInternalServerError)
		trace.Status = http.StatusInternalServerError
		return
	}

	/*post request to message queue*/
	msg, respErr := scyna.Nats.Request(scyna_utils.PublishURL(url), reqBytes, 60*time.Second)
	if respErr != nil {
		http.Error(rw, "No response", http.StatusInternalServerError)
		trace.Status = http.StatusInternalServerError
		scyna.Session.Error("ServeHTTP: Nats: " + respErr.Error())
		return
	}

	/*response*/
	if err := proto.Unmarshal(msg.Data, &ctx.Response); err != nil {
		http.Error(rw, "Cannot deserialize response", http.StatusInternalServerError)
		scyna.Session.Error("nats-proxy:" + err.Error())
		trace.Status = http.StatusInternalServerError
		return
	}

	rw.WriteHeader(int(ctx.Response.Code))
	_, err = bytes.NewBuffer(ctx.Response.Body).WriteTo(rw)
	if err != nil {
		scyna.Session.Error("Proxy write data error: " + err.Error())
		trace.Status = 0
	}

	trace.Status = uint32(ctx.Response.Code)

	if f, ok := rw.(http.Flusher); ok {
		f.Flush()
	}
}
