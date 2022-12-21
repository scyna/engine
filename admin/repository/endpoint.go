package repository

import (
	"github.com/scylladb/gocqlx/v2/qb"
	scyna "github.com/scyna/core"
)

type Endpoint struct {
	Context string `db:"context"`
	URL     string `db:"url"`
	Name    string `db:"name"`
}

func CreateEndpoint(LOG scyna.Logger, endpoint *Endpoint) *scyna.Error {
	if err := qb.Insert("scyna.endpoint").
		Columns("context", "url", "name").
		Query(scyna.DB).
		BindStruct(endpoint).
		ExecRelease(); err != nil {
		return scyna.SERVER_ERROR
	}
	return nil
}
