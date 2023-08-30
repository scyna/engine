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
		seed += snPartitionSize
	} else {
		log.Print("nextBucket:", err)
	}

	if err := scyna.DB.Execute("INSERT INTO "+scyna_const.GEN_SN_TABLE+"(key, prefix, seed) VALUES (?, ?, ?) IF NOT EXISTS",
		key, prefix, seed); err == nil {
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
