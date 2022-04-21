package manager

import scyna "github.com/scyna/go"

const (
	MODULE_CODE   = "scyna.engine"
	MODULE_SECRET = "123456"
)

var DefaultConfig *scyna.Configuration = &scyna.Configuration{
	NatsUrl:      "127.0.0.1",
	NatsUsername: "",
	NatsPassword: "",
	DBHost:       "127.0.0.1",
	DBUsername:   "",
	DBPassword:   "",
	DBLocation:   "",
}
