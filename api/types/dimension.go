package types

import (
	"fmt"

	"github.com/awatercolorpen/olap-sql/api/proto"
)

type Dimension struct {
	Table     string `json:"table"`
	Name      string `json:"name"`
	FieldName string `json:"field_name"`
}

func (d *Dimension) Statement() (string, error) {
	if d.Name == "" {
		return d.FieldName, nil
	}
	return fmt.Sprintf("%v.%v AS %v", d.Table, d.FieldName, d.Name), nil
}

func (d *Dimension) ToProto() *proto.Dimension {
	return &proto.Dimension{
		Table:     d.Table,
		Name:      d.Name,
		FieldName: d.FieldName,
	}
}

func ProtoToDimension(d *proto.Dimension) *Dimension {
	return &Dimension{
		Table:     d.Table,
		Name:      d.Name,
		FieldName: d.FieldName,
	}
}
