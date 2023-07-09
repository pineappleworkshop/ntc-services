package config

const (
	VERSION      = "0.0.13"
	PORT         = 3444
	SERVICE_NAME = "ntc-services"

	TEST        = "test"
	WORKSTATION = "workstation"
	DEV         = "dev"
	STAGE       = "stage"
	PROD        = "prod"

	CONSUL_KV           = "ntc"
	CONSUL_HOST_DEV     = "localhost"
	CONSUL_PORT_DEV     = "8500"
	CONSUL_HOST_CLUSTER = "consul-server"
	CONSUL_PORT_CLUSTER = "8500"
)

var (
	MONGOHOSTS_WORKSTATION = []string{"localhost:27017"}
)
