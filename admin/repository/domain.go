package repository

import (
	"github.com/scylladb/gocqlx/v2/qb"
	scyna "github.com/scyna/core"
)

type Domain struct {
	Code string `db:"code"`
	Name string `db:"name"`
}

func GetDomain(LOG scyna.Logger, code string) (*scyna.Error, *Domain) {
	var domain Domain

	if err := qb.Select("scyna.domain").
		Columns("code", "name").
		Where(qb.Eq("code")).
		Limit(1).
		Query(scyna.DB).Bind(code).GetRelease(&domain); err != nil {
		LOG.Error(err.Error())
		return scyna.REQUEST_INVALID, nil
	}

	return nil, &domain
}

func CreateDomain(LOG scyna.Logger, domain *Domain) *scyna.Error {

	if err := qb.Insert("scyna.domain").
		Columns("code", "name").
		Query(scyna.DB).
		BindStruct(domain).
		ExecRelease(); err != nil {
		return scyna.SERVER_ERROR
	}

	return nil
}
