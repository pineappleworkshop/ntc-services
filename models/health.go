package models

type Health struct {
	Service  string      `json:"service"`
	Status   int         `json:"status"`
	Version  string      `json:"version"`
	Database bool        `json:"database"`
	NodeType string      `json:"node_type"`
	State    interface{} `json:"state"`
}
