package gateway

import (
	"bytes"
	"net/http"
	"strings"
	"time"

	scyna "github.com/scyna/core"
	scyna_utils "github.com/scyna/core/utils"
	"google.golang.org/protobuf/proto"
)

type Gateway struct {
	Queries  QueryPool
	Contexts scyna_utils.HttpContextPool
}

func NewGateway() *Gateway {
	ret := &Gateway{
		Queries:  NewQueryPool(),
		Contexts: scyna_utils.NewContextPool(),
	}
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

	/*CORS*/
	rw.Header().Set("Content-Type", contentType)
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	trace := scyna.Trace{
		ID:        callID,
		ParentID:  0,
		Time:      time.Now(),
		Path:      url,
		Type:      scyna.TRACE_ENDPOINT,
		SessionID: scyna.Session.ID(),
	}
	defer saveTrace(trace)

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

	/*build request*/
	err := ctx.Request.Build(req)
	if err != nil {
		http.Error(rw, "Cannot process request", http.StatusInternalServerError)
		trace.Status = http.StatusInternalServerError
		return
	}
	trace.RequestBody = string(ctx.Request.Body)
	ctx.Request.TraceID = callID

	/*serialize the request */
	reqBytes, err := proto.Marshal(&ctx.Request)
	if err != nil {
		http.Error(rw, "Cannot process request", http.StatusInternalServerError)
		trace.Status = http.StatusInternalServerError
		return
	}

	/*post request to message queue*/
	msg, respErr := scyna.Connection.Request(scyna_utils.PublishURL(url), reqBytes, 60*time.Second)
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

	trace.SessionID = ctx.Response.SessionID
	trace.Status = ctx.Response.Code

	if f, ok := rw.(http.Flusher); ok {
		f.Flush()
	}
}
