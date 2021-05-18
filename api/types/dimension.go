package types

import (
	"fmt"
	"regexp"

	"github.com/awatercolorpen/olap-sql/api/proto"
)

var (
	reg = regexp.MustCompile(`^[0-9A-Za-z_]+$`)
)

type Dimension struct {
	Table     string `json:"table"`
	Name      string `json:"name"`
	FieldName string `json:"field_name"`
}

func (d *Dimension) Expression() (string, error) {
	if !reg.MatchString(d.FieldName) {
		return d.FieldName, nil
	}
	return fmt.Sprintf("%v.%v", d.Table, d.FieldName), nil
}

func (d *Dimension) Alias() (string, error) {
	return d.Name, nil
}

func (d *Dimension) Statement() (string, error) {
	expression, err := d.Expression()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%v AS %v", expression, d.Name), nil
}

func (d *Dimension) expression() bool {
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
