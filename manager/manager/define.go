package manager

import (
	scyna_proto "github.com/scyna/core/proto/generated"
)

const (
	CONTEXT_CODE  = "scyna.engine"
	MODULE_SECRET = "123456"
)

var DefaultConfig *scyna_proto.Configuration = &scyna_proto.Configuration{
	NatsUrl:      "127.0.0.1",
	NatsUsername: "",
	NatsPassword: "",
	DBHost:       "127.0.0.1",
	DBUsername:   "",
	DBPassword:   "",
	DBLocation:   "",
}
