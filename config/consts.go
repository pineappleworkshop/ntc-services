package config

const (
	VERSION      = "0.0.18"
	PORT         = 3444
	SERVICE_NAME = "omnisat-api"

	TEST        = "test"
	WORKSTATION = "workstation"
	DEV         = "dev"
	STAGE       = "stage"
	PROD        = "prod"

	CONSUL_KV           = "omnisat-api"
	CONSUL_HOST_DEV     = "localhost"
	CONSUL_PORT_DEV     = "8500"
	CONSUL_HOST_CLUSTER = "consul-server"
	CONSUL_PORT_CLUSTER = "8500"
)

var (
	MONGOHOSTS_WORKSTATION = []string{"localhost:27017"}
)
