package olapsql

import (
	"fmt"

	"github.com/awatercolorpen/olap-sql/api/types"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// ClientsOption is a map from connection-key to DBOption.
// The key is used to look up the correct database connection when running a query.
// Typical keys follow the pattern "<dbtype>" (e.g. "clickhouse") or
// "<dbtype>/<dataset>" for dataset-scoped connections.
type ClientsOption = map[string]*DBOption

// Clients is a registry of open *gorm.DB connections, keyed by "<dbtype>" or "<dbtype>/<dataset>".
type Clients map[string]*gorm.DB

// RegisterByKV registers a *gorm.DB connection under the composite key derived from dbType and dataset.
func (c Clients) RegisterByKV(dbType types.DBType, dataset string, db *gorm.DB) {
	key := c.key(dbType, dataset)
	c[key] = db
}

// RegisterByOption opens database connections for each entry in option
// and registers them in the Clients map.
// Returns an error if any connection cannot be established.
func (c Clients) RegisterByOption(option ClientsOption) error {
	for k, v := range option {
		db, err := v.NewDB()
		if err != nil {
			return err
		}
		c[k] = db
	}
	return nil
}

// SetLogger replaces the GORM logger on every registered connection.
// Call this to enable SQL statement logging or to plug in a custom logger.
func (c Clients) SetLogger(log logger.Interface) {
	for _, v := range c {
		v.Config.Logger = log
	}
}

// Get returns the *gorm.DB for the given dbType and dataset.
// If no dataset-specific connection is registered, it falls back to the
// type-level connection (dataset == "").
// Returns an error if neither key is found.
func (c Clients) Get(dbType types.DBType, dataset string) (*gorm.DB, error) {
	key1 := c.key(dbType, dataset)
	if v, ok := c[key1]; ok {
		return v, nil
	}
	key2 := c.key(dbType, "")
	if v, ok := c[key2]; ok {
		return v, nil
	}
	return nil, fmt.Errorf("not found client %v %v", dbType, dataset)
}

// key builds the internal lookup key for a (dbType, dataset) pair.
func (c Clients) key(dbType types.DBType, dataset string) string {
	if dataset == "" {
		return fmt.Sprintf("%v", dbType)
	}
	return fmt.Sprintf("%v/%v", dbType, dataset)
}

// BuildDB selects the correct client for the clause and constructs
// a *gorm.DB with the translated query applied.
func (c Clients) BuildDB(clause types.Clause) (*gorm.DB, error) {
	client, err := c.Get(clause.GetDBType(), clause.GetDataset())
	if err != nil {
		return nil, err
	}
	return clause.BuildDB(client)
}

// BuildSQL selects the correct client for the clause and returns
// the SQL string that would be executed, without actually running it.
func (c Clients) BuildSQL(clause types.Clause) (string, error) {
	client, err := c.Get(clause.GetDBType(), clause.GetDataset())
	if err != nil {
		return "", err
	}
	return clause.BuildSQL(client)
}

// NewClients creates a Clients registry by opening connections for each DBOption in option.
// Returns an error if any connection fails to open.
func NewClients(option ClientsOption) (Clients, error) {
	c := Clients{}
	if err := c.RegisterByOption(option); err != nil {
		return nil, err
	}
	return c, nil
}
