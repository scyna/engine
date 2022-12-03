package repository

import (
	"github.com/scylladb/gocqlx/v2/qb"
	scyna "github.com/scyna/core"
)

func CheckApplication(org string, code string) bool {
	if err := qb.Select("scyna.application").
		Columns("code").
		Where(qb.Eq("org"), qb.Eq("code")).
		Query(scyna.DB).
		Bind(org, code).
		GetRelease(&code); err != nil {
		return false
	}
	return true
}
