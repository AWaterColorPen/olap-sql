package olapsql

var (
	DefaultParallelNumber = 3
)

type Configuration struct {
	// default olap-sql option
	DefaultParallelNumber int                   `json:"default_parallel_number"`

	// configurations for data_dictionary
	DataDictionaryOption  *DataDictionaryOption `json:"data_dictionary_option"`
}
