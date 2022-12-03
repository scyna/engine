package repository

import (
	"github.com/scylladb/gocqlx/v2/qb"
	scyna "github.com/scyna/core"
)

func CheckOrganization(code string) bool {
	if err := qb.Select("scyna.organization").
		Columns("code").
		Where(qb.Eq("code")).
		Query(scyna.DB).
		Bind(code).
		GetRelease(&code); err != nil {
		return false
	}
	return true
}
