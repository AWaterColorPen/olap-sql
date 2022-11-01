package olapsql

type Configuration struct {
	// configurations for clients, data_dictionary
	ClientsOption    ClientsOption `json:"clients_option"`
	DictionaryOption *Option       `json:"dictionary_option"`
}
