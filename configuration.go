package olapsql

// Configuration holds the top-level options used to initialise a Manager.
// Both fields are optional; omitting ClientsOption creates a dictionary-only
// instance (useful for SQL generation without execution), and omitting
// DictionaryOption creates a client-only instance.
type Configuration struct {
	// ClientsOption maps connection-key strings to database connection options.
	// Each key typically corresponds to a database type (e.g. "clickhouse").
	ClientsOption ClientsOption `json:"clients_option"`

	// DictionaryOption configures the OLAP schema adapter (e.g. a TOML file path).
	DictionaryOption *Option `json:"dictionary_option"`
}
