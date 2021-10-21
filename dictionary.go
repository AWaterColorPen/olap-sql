package olapsql

import (
	"github.com/awatercolorpen/olap-sql/api/types"
	"github.com/awatercolorpen/olap-sql/dictionary"
)

type Option struct {
	dictionary.AdapterOption
}
type Dictionary struct {
	adapter dictionary.Adapter
}

func (d *Dictionary) GetAdapter() dictionary.Adapter {
	return d.adapter
}

func (d *Dictionary) Create(item interface{}) error {
	return nil
}

func (d *Dictionary) Translator(query *types.Query) (Translator, error) {
	set, err := d.adapter.GetDataSetByName(query.DataSetName)
	if err != nil {
		return nil, err
	}

	id := set.Schema.DataSourceID()
	sources, err := d.adapter.GetSourcesByIds(id)
	if err != nil {
		return nil, err
	}

	metrics, err := d.adapter.GetMetricsByIds(id)
	if err != nil {
		return nil, err
	}

	dimensions, err := d.adapter.GetDimensionsByIds(id)
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
	return &Dictionary{adapter: adapter}, nil
}
