package olapsql

import "github.com/awatercolorpen/olap-sql/api/types"

type Translator interface {
	Translate(interface{}) (interface{}, error)
}

type translator struct {

}

func (t *translator) Translate(in interface{}) (interface{}, error) {
	_, ok := in.(*types.Query)
	if !ok {
		return nil, nil
	}



	request := &types.Request{}
	return request, nil
}