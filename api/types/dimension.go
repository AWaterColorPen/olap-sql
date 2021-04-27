package types

import (
	"fmt"
	
	"github.com/awatercolorpen/olap-sql/api/proto"
)

type Dimension struct {
	FieldName string `json:"field_name"`
	Name      string `json:"name"`
}

func (d *Dimension) Statement() (string, error) {
	if d.Name == "" {
		return d.FieldName, nil
	}
	return fmt.Sprintf("%v AS %v", d.FieldName, d.Name), nil
}

func (d *Dimension) ToProto() *proto.Dimension {
	return &proto.Dimension{
		FieldName: d.FieldName,
		Name:      d.Name,
	}
}

func ProtoToDimension(d *proto.Dimension) *Dimension {
	return &Dimension{
		FieldName: d.FieldName,
		Name:      d.Name,
	}
}
