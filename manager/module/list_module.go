package module

import (
	"log"

	scyna "github.com/scyna/core"
	PROTO "github.com/scyna/engine/proto/generated"
)

func ListModuleHandler(ctx *scyna.Endpoint, request *PROTO.ListModuleRequest) scyna.Error {
	log.Println("Receive ListModuleRequest")

	rs := scyna.DB.QueryMany("SELECT code FROM scyna.module")

	response := &PROTO.ListModuleResponse{}
	for rs.Next() {
		var code string
		if err := rs.Scan(&code); err != nil {
			ctx.Error(err.Error())
			return scyna.REQUEST_INVALID
		}

		response.Items = append(response.Items, code)
	}

	return ctx.OK(response)
}
