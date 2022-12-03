package repository

import (
	"github.com/scylladb/gocqlx/v2/qb"
	scyna "github.com/scyna/core"
	"github.com/scyna/go/manager/model"
)

func RemoveService(LOG scyna.Logger, client string, service string) *scyna.Error {
	if applied, err := qb.Delete("scyna.client_use_service").
		Where(qb.Eq("client_id"), qb.Eq("service_url")).
		Existing().Query(scyna.DB).
		Bind(client, service).
		ExecCASRelease(); !applied {
		if err == nil {
			return model.CLIENT_EXISTED
		} else {
			LOG.Info("Can not delete client use service " + client + " : " + err.Error())
			return scyna.SERVER_ERROR
		}
	}
	return nil
}
