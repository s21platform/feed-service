package model

type EntityInfo struct {
	Uuid     string `db:"external_uuid"`
	Metadata string `db:"metadata"`
}
