package generator

import (
	"log"
	"time"

	"github.com/scylladb/gocqlx/v2/qb"
	scyna "github.com/scyna/core"
	scyna_proto "github.com/scyna/core/proto/generated"
)

const snPartitionSize = 500

func GetSN(s scyna.Context, request *scyna_proto.GetSNRequest) scyna.Error {
	log.Print("Receive GetSNRequest")

	for i := 0; i < tryCount; i++ {
		if bucket := nextBucket(request.Key); bucket != nil {
			return s.OK(bucket)
		}
	}

	return scyna.SERVER_ERROR
}

func nextBucket(key string) *scyna_proto.GetSNResponse {
	prefix := time.Now().Unix() / (60 * 60 * 24)
	seed := 0
	if err := qb.Select("scyna.gen_sn").
		Columns("seed").
		Where(qb.Eq("key"), qb.Eq("prefix")).
		Limit(1).
		Query(scyna.DB).
		Bind(key, prefix).
		GetRelease(&seed); err == nil {
		seed += snPartitionSize
	} else {
		log.Print("OneID:", err)
	}

	if applied, err := qb.Insert("scyna.gen_sn").
		Columns("key", "prefix", "seed").
		Unique().
		Query(scyna.DB).
		Bind(key, prefix, seed).
		ExecCASRelease(); applied {
		return &scyna_proto.GetSNResponse{
			Prefix: uint32(prefix),
			Start:  uint64(seed) + 1,
			End:    uint64(seed) + snPartitionSize,
		}
	} else {
		if err != nil {
			log.Print("nextBucket: insert seed: ", err.Error())
		}
	}
	return nil
}
