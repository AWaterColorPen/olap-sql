package olapsql_test

import (
	"os"
	"time"
)

type WikiStat struct {
	Date       time.Time `gorm:"column:date"       json:"date"`
	Time       time.Time `gorm:"column:time"       json:"time"`
	SubProject string    `gorm:"column:subproject" json:"subproject"`
	Path       string    `gorm:"column:path"       json:"path"`
	Hits       uint64    `gorm:"column:hits"       json:"hits"`
	Size       uint64    `gorm:"column:size"       json:"size"`
}

func (WikiStat) TableName() string {
	return "wikistat"
}



func DataWithClickhouse() bool {
	args := os.Args
	for _, arg := range args {
		if arg == "clickhouse" {
			return true
		}
	}
	return false
}


