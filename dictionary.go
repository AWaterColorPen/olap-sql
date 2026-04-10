package olapsql

import (
	"github.com/awatercolorpen/olap-sql/api/types"
)

// Option wraps AdapterOption and is the configuration type for the OLAP dictionary.
type Option struct {
	AdapterOption
}

// Dictionary holds the schema adapter and provides query translation capabilities.
// It converts a high-level [types.Query] into a backend-specific [types.Clause].
type Dictionary struct {
	Adapter IAdapter
}

// Translator builds a [Translator] for the given query by looking up the
// target data set in the adapter and wiring up the appropriate translator options.
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

// Translate is a convenience method that builds a Translator and runs it in one call.
// It returns the resulting [types.Clause] ready to be executed against a database client.
func (d *Dictionary) Translate(query *types.Query) (types.Clause, error) {
	translator, err := d.Translator(query)
	if err != nil {
		return nil, err
	}
	return translator.Translate(query)
}

// NewDictionary creates a Dictionary from the provided Option.
// The option's DSN is used to locate and parse the schema file (e.g. a TOML config).
// Returns an error if the adapter cannot be initialised.
func NewDictionary(option *Option) (*Dictionary, error) {
	adapter, err := NewAdapter(&option.AdapterOption)
	if err != nil {
		return nil, err
	}
	return &Dictionary{Adapter: adapter}, nil
}
