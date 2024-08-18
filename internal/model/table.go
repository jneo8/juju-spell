package model

import "github.com/jneo8/jujuspell/internal/client"

type Table struct {
	resource client.Resource
}

func NewTable(resource client.Resource) *Table {
	return &Table{
		resource: resource,
	}
}
