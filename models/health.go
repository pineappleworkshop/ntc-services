package models

type Health struct {
	Service          string      `json:"service"`
	Status           int         `json:"status"`
	Version          string      `json:"version"`
	Database         bool        `json:"database"`
	NodeType         string      `json:"node_type"`
	BestInSlotStatus int         `json:"bis_status"`
	OrdexStatus      int         `json:"ordex_status"`
	State            interface{} `json:"state"`
}
