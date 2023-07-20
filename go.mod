module github.com/scyna/engine

go 1.19

//replace github.com/scyna/core => ../core

require (
	github.com/go-ozzo/ozzo-validation v3.6.0+incompatible
	github.com/gocql/gocql v1.0.0
	github.com/nats-io/nats.go v1.23.0
	github.com/scylladb/gocqlx/v2 v2.7.0
	github.com/scyna/core v1.2.8
	google.golang.org/protobuf v1.28.0
)

require github.com/nats-io/nats-server/v2 v2.9.14 // indirect

require (
	github.com/asaskevich/govalidator v0.0.0-20210307081110-f21760c49a8d // indirect
	github.com/aws/aws-sigv4-auth-cassandra-gocql-driver-plugin v0.0.0-20220331165046-e4d000c0d6a6 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/hailocab/go-hostpool v0.0.0-20160125115350-e80d13ce29ed // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/nats-io/nkeys v0.3.0 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	github.com/scylladb/go-reflectx v1.0.1 // indirect
	github.com/stretchr/testify v1.7.0 // indirect
	golang.org/x/crypto v0.5.0 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
)
