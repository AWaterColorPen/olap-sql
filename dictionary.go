package olapsql

import (
	"github.com/awatercolorpen/olap-sql/api/types"
)

type Option struct {
	AdapterOption
}
type Dictionary struct {
	Adapter IAdapter
}

func (d *Dictionary) Translator(query *types.Query) (Translator, error) {
	set, err := d.Adapter.GetDataSetByKey(query.DataSetName)
	if err != nil {
		return nil, err
	}
	option := &TranslatorOption{
		Adapter: d.Adapter,
		Query:   query,
		DBType:  set.DBType,
		Current: set.GetCurrent(),
	}
	return NewTranslator(option)
}

func (d *Dictionary) Translate(query *types.Query) (types.Clause, error) {
	translator, err := d.Translator(query)
	if err != nil {
		return nil, err
	}
	return translator.Translate(query)
}

func NewDictionary(option *Option) (*Dictionary, error) {
	adapter, err := NewAdapter(&option.AdapterOption)
	if err != nil {
		return nil, err
	}
	return &Dictionary{Adapter: adapter}, nil
}
