package model

type Entity struct {
	ExternalUUID string `db:"external_uuid"`
	Metadata     string `db:"metadata"`
}
