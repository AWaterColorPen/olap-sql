package models

import (
	"strings"
)

type DataSource struct {
	ID          uint64               `toml:"id"          json:"id,omitempty"`
	Name        string               `toml:"name"        json:"name"`
	Description string               `toml:"description" json:"description"`
}

func (d *DataSource) GetTableName() string {
	out := strings.Split(d.Name, ".")
	return out[len(out)-1]
}

func (d *DataSource) GetDatabaseName() string {
	if !strings.Contains(d.Name, ".") {
		return ""
	}
	out := strings.Split(d.Name, ".")
	return out[0]
}
