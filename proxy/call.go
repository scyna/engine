package proxy

import (
	"log"
	"time"

	scyna "github.com/scyna/go"

	"github.com/gocql/gocql"
)

func (proxy *Proxy) SaveErrorCall(client string, status int, id uint64, day int, start time.Time, url string) {
	qBatch := scyna.DB.NewBatch(gocql.LoggedBatch)
	qBatch.Query("INSERT INTO scyna.call(id, day, time, source, status, caller_id)"+
		" VALUES (?,?,?,?,?,?)", id, day, start, url, status, client)
	qBatch.Query("INSERT INTO scyna.client_has_call(client_id, call_id, day) VALUES (?,?,?)", client, id, day)
	if err := scyna.DB.ExecuteBatch(qBatch); err != nil {
		log.Print(err)
		scyna.LOG.Error("Error in save call")
	}
}

func (proxy *Proxy) SaveCall(client string, id uint64, day int, start time.Time, duration int64, url string, service *scyna.Service) {
	qBatch := scyna.DB.NewBatch(gocql.LoggedBatch)
	qBatch.Query("INSERT INTO scyna.call(id, day, time, duration, source,request, response, status, session_id, caller_id)"+
		" VALUES (?,?,?,?,?,?,?,?,?,?)", id, day, start, duration, url, service.Request.Body, service.Response.Body,
		service.Response.Code, service.Response.SessionID, client,
	)
	qBatch.Query("INSERT INTO scyna.client_has_call(client_id, call_id, day) VALUES (?,?,?)", client, id, day)
	if err := scyna.DB.ExecuteBatch(qBatch); err != nil {
		log.Print(err)
		scyna.LOG.Error("Error in save call")
	}
}
