package olapsql

import "github.com/awatercolorpen/olap-sql/dictionary"

var (
	DefaultParallelNumber = 3
)

type Configuration struct {
	// default olap-sql option
	DefaultParallelNumber int `json:"default_parallel_number"`

	// configurations for clients, data_dictionary
	ClientsOption        ClientsOption                `json:"clients_option"`
	DataDictionaryOption *dictionary.DictionaryOption `json:"data_dictionary_option"`
}
