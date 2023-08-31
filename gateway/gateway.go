package gateway

import (
	"bytes"
	"log"
	"net/http"
	"strings"
	"time"

	scyna "github.com/scyna/core"
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
}

type Gateway struct {
	Contexts        scyna_utils.HttpContextPool
	PublicEndpoints []string
}

func NewGateway() *Gateway {
	ret := &Gateway{Contexts: scyna_utils.NewContextPool()}
	ret.initPublicEndpoints()
	return ret
}

func (proxy *Gateway) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	callID := scyna.ID.Next()

	/*authenticate*/
	url := req.URL.String()
	contentType := req.Header.Get("Content-Type")

	//https://descynaper.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Type
	for _, data := range strings.Split(contentType, ";") {
		value := strings.TrimSpace(strings.Trim(data, ";"))
		if strings.HasPrefix(value, "application/") {
			contentType = value
			continue
		}
	}

	log.Println("Request: ", url, req.Header.Get("Origin"))

	//Accept preflight request
	if req.Method == "OPTIONS" {
		rw.Header().Set("Access-Control-Allow-Origin", "*")
		rw.Header().Set("Access-Control-Allow-Credentials", "true")
		rw.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		rw.Header().Set("Access-Control-Allow-Methods", "POST")
		rw.WriteHeader(http.StatusOK)
		return
	}

	/*CORS*/
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.Header().Set("Access-Control-Allow-Credentials", "true")
	rw.Header().Set("Access-Control-Allow-Headers", "Content-Type, access-control-allow-origin, access-control-allow-headers")
	rw.Header().Set("Access-Control-Allow-Methods", "POST")
	rw.Header().Set("Content-Type", req.Header.Get("Content-Type"))

	trace := trace{
		ID:        callID,
		ParentID:  0,
		Time:      time.Now(),
		Path:      url,
		Type:      scyna.TRACE_ENDPOINT,
		SessionID: scyna.Session.ID(),
	}
	defer saveTrace(&trace)

	ctx := proxy.Contexts.GetContext()
	defer proxy.Contexts.PutContext(ctx)

	if contentType == "application/json" {
		ctx.Request.JSON = true
	} else if contentType == "application/protobuf" {
		ctx.Request.JSON = false
	} else {
		http.Error(rw, "Content-Type must be JSON or PROTOBUF ", http.StatusNotAcceptable)
		trace.Status = http.StatusNotAcceptable
		return
	}

	/* check url within public endpoint */
	if !proxy.isPublicEndpoint(url) {
		http.Error(rw, "Unauthorized", http.StatusUnauthorized)
		trace.Status = http.StatusUnauthorized
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

	if f, ok := rw.(http.Flusher); ok {
		f.Flush()
	}
}
