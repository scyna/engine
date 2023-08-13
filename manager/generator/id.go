package generator

import (
	"log"
	"time"

	scyna "github.com/scyna/core"
	scyna_const "github.com/scyna/core/const"
	scyna_proto "github.com/scyna/core/proto/generated"
)

const idPartitionSize = 1000
const tryCount = 10

func Init() {
	for i := 0; i < tryCount; i++ {
		if ok, prefix, start, end := allocate(); ok {
			scyna.ID.Reset(prefix, end, start)
			return
		}
	}
	panic("Can not init id generator")
}

func GetID(ctx *scyna.Endpoint, request *scyna_proto.EmptyRequest) scyna.Error {
	log.Print("Receive GetIDRequest")
	for i := 0; i < tryCount; i++ {
		if ok, prefix, start, end := allocate(); ok {
			return ctx.OK(&scyna_proto.GetIDResponse{
				Prefix: prefix,
				Start:  start,
				End:    end,
			})
		}
	}
	return scyna.SERVER_ERROR
}

func allocate() (ok bool, prefix uint32, start uint64, end uint64) {
	p := time.Now().Unix() / (60 * 60 * 24)
	ok = false

	seed := 0
	if err := scyna.DB.QueryOne("SELECT seed FROM "+scyna_const.GEN_ID_TABLE+
		" WHERE prefix = ?", p).Scan(&seed); err == nil {
		// if err := qb.Select(scyna_const.GEN_ID_TABLE).
		// 	Columns("seed").
		// 	Where(qb.Eq("prefix")).
		// 	Limit(1).
		// 	Query(scyna.DB).
		// 	Bind(p).
		// 	GetRelease(&seed); err == nil {
		seed += idPartitionSize
	} else {
		log.Println("generator.allocate: get seed: " + err.Error())
	}

	if err := scyna.DB.Execute("INSERT INTO "+scyna_const.GEN_ID_TABLE+
		" (prefix, seed) VALUES (?, ?) IF NOT EXISTS", p, seed); err == nil {
		// if applied, err := qb.Insert(scyna_const.GEN_ID_TABLE).
		// 	Columns("prefix", "seed").
		// 	Unique().
		// 	Query(scyna.DB).
		// 	Bind(p, seed).
		// 	ExecCASRelease(); applied {
		ok = true
		prefix = uint32(p)
		start = uint64(seed) + 1
		end = uint64(seed) + idPartitionSize
	} else {
		log.Println("generator.allocate: insert: " + err.Error())
	}
	return
}
