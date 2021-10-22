package olapsql

import (
	"github.com/awatercolorpen/olap-sql/api/types"
	"github.com/awatercolorpen/olap-sql/dictionary"
)

type Option struct {
	dictionary.AdapterOption
}
type Dictionary struct {
	Adapter dictionary.Adapter
}

func (d *Dictionary) Translator(query *types.Query) (Translator, error) {
	set, err := d.Adapter.GetDataSetByName(query.DataSetName)
	if err != nil {
		return nil, err
	}

	id := set.Schema.DataSourceID()
	sources, err := d.Adapter.GetSourcesByIds(id)
	if err != nil {
		return nil, err
	}

	metrics, err := d.Adapter.GetMetricsByIds(id)
	if err != nil {
		return nil, err
	}

	dimensions, err := d.Adapter.GetDimensionsByIds(id)
	if err != nil {
		return nil, err
	}

	t := &DictionaryTranslator{
		set:        set,
		sources:    sources,
		metrics:    metrics,
		dimensions: dimensions,
	}
	return t, nil
}

func (d *Dictionary) Translate(query *types.Query) (Clause, error) {
	translator, err := d.Translator(query)
	if err != nil {
		return nil, err
	}
	return translator.Translate(query)
}

func NewDictionary(option *Option) (*Dictionary, error) {
	adapter, err := dictionary.NewAdapter(&option.AdapterOption)
	if err != nil {
		return nil, err
	}
	return &Dictionary{Adapter: adapter}, nil
}
