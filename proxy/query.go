package proxy

import (
	"sync"

	scyna "github.com/scyna/core"
	scyna_const "github.com/scyna/core/const"

	"github.com/scylladb/gocqlx/v2"
	"github.com/scylladb/gocqlx/v2/qb"
)

type Query struct {
	Authenticate *gocqlx.Queryx
}

type QueryPool struct {
	sync.Pool
}

func NewQuery() *Query {
	return &Query{
		Authenticate: qb.Select(scyna_const.CLIENT_USE_ENDPOINT_TABLE).
			Columns("url").
			Where(qb.Eq("client"), qb.Eq("url")).
			Limit(1).
			Query(scyna.DB),
	}
}

func (q *QueryPool) GetQuery() *Query {
	query, _ := q.Get().(*Query)
	return query
}

func (q *QueryPool) PutQuery(query *Query) {
	q.Put(query)
}

func NewQueryPool() QueryPool {
	return QueryPool{
		sync.Pool{
			New: func() interface{} { return NewQuery() },
		}}
}
