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
	adapter, err := d.Adapter.BuildDataSetAdapter(query.DataSetName)
	if err != nil {
		return nil, err
	}

	option := &TranslatorOption{
		Adapter: adapter,
		Query:   query,
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
