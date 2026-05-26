package models

type Status string

const (
	draft     Status = "Draft"
	published Status = "Published"
	archived  Status = "Archived"
)
