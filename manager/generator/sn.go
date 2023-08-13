package generator

import (
	"log"
	"time"

	scyna "github.com/scyna/core"
	scyna_const "github.com/scyna/core/const"
	scyna_proto "github.com/scyna/core/proto/generated"
)

const snPartitionSize = 500

func GetSN(ctx *scyna.Endpoint, request *scyna_proto.GetSNRequest) scyna.Error {
	log.Print("Receive GetSNRequest")

	for i := 0; i < tryCount; i++ {
		if bucket := nextBucket(request.Key); bucket != nil {
			return ctx.OK(bucket)
		}
	}

	return scyna.SERVER_ERROR
}

func nextBucket(key string) *scyna_proto.GetSNResponse {
	prefix := time.Now().Unix() / (60 * 60 * 24)
	seed := 0

	if err := scyna.DB.QueryOne("SELECT seed FROM "+scyna_const.GEN_SN_TABLE+
		" WHERE key = ? AND prefix = ?", key, prefix).Scan(&seed); err == nil {
		// if err := qb.Select(scyna_const.GEN_SN_TABLE).
		// 	Columns("seed").
		// 	Where(qb.Eq("key"), qb.Eq("prefix")).
		// 	Limit(1).
		// 	Query(scyna.DB).
		// 	Bind(key, prefix).
		// 	GetRelease(&seed); err == nil {
		seed += snPartitionSize
	} else {
		log.Print("nextBucket:", err)
		return nil
	}

	if err := scyna.DB.Execute("UPDATE "+scyna_const.GEN_SN_TABLE+
		" SET key = ?, seed = ?, prefix = ? IF NOT EXISTS", key, seed, prefix); err == nil {
		// if applied, err := qb.Insert(scyna_const.GEN_SN_TABLE).
		// 	Columns("key", "prefix", "seed").
		// 	Unique().
		// 	Query(scyna.DB).
		// 	Bind(key, prefix, seed).
		// 	ExecCASRelease(); applied {
		return &scyna_proto.GetSNResponse{
			Prefix: uint32(prefix),
			Start:  uint64(seed) + 1,
			End:    uint64(seed) + snPartitionSize,
		}
	} else {
		log.Print("nextBucket: insert seed: ", err.Error())
		return nil
	}
}
