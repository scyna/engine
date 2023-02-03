package manager

import (
	scyna_engine "github.com/scyna/core/engine"
)

const (
	MODULE_CODE   = "scyna.engine"
	MODULE_SECRET = "123456"
)

var DefaultConfig *scyna_engine.Configuration = &scyna_engine.Configuration{
	NatsUrl:      "127.0.0.1",
	NatsUsername: "",
	NatsPassword: "",
	DBHost:       "127.0.0.1",
	DBUsername:   "",
	DBPassword:   "",
	DBLocation:   "",
}
