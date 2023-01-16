package repository

import (
	"github.com/scylladb/gocqlx/v2/qb"
	scyna "github.com/scyna/core"
)

type Context struct {
	Domain string `db:"domain"`
	Code   string `db:"code"`
	Name   string `db:"name"`
}

func GetContext(LOG scyna.Logger, code string) (scyna.Error, *Context) {
	var context Context

	if err := qb.Select("scyna.context").
		Columns("code", "name").
		Where(qb.Eq("code")).
		Limit(1).
		Query(scyna.DB).Bind(code).GetRelease(&context); err != nil {
		LOG.Error(err.Error())
		return scyna.REQUEST_INVALID, nil
	}

	return nil, &context
}

func CreateContext(LOG scyna.Logger, context *Context) scyna.Error {

	if err := qb.Insert("scyna.context").
		Columns("domain", "code", "name").
		Query(scyna.DB).
		BindStruct(context).
		ExecRelease(); err != nil {
		return scyna.SERVER_ERROR
	}

	return nil
}
