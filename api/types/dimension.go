package types

import (
	"fmt"
	"regexp"

	"github.com/awatercolorpen/olap-sql/api/proto"
)

var (
	reg = regexp.MustCompile(`^[0-9A-Za-z]+$`)
)

type Dimension struct {
	Table     string `json:"table"`
	Name      string `json:"name"`
	FieldName string `json:"field_name"`
}

func (d *Dimension) Statement() (string, error) {
	if d.expression() {
		return fmt.Sprintf("%v AS %v", d.FieldName, d.Name), nil
	}
	return fmt.Sprintf("%v.%v AS %v", d.Table, d.FieldName, d.Name), nil
}

func (d *Dimension) expression () bool {
	return !reg.MatchString(d.FieldName)
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
