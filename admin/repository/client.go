package repository

import (
	"github.com/scylladb/gocqlx/v2/qb"
	scyna "github.com/scyna/core"
)

type Client struct {
	Domain string `db:"domain"`
	ID     string `db:"id"`
	Secret string `db:"secret"`
}

func CreateClient(LOG scyna.Logger, client *Client) *scyna.Error {

	if err := qb.Insert("scyna.client").
		Columns("domain", "id", "secret").
		Query(scyna.DB).
		BindStruct(client).
		ExecRelease(); err != nil {
		return scyna.SERVER_ERROR
	}

	return nil
}
