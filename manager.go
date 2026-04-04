// Package olapsql provides a Go library for generating adapted SQL from OLAP queries.
// It supports metrics, dimensions, and filters to automatically produce SQL
// for multiple database backends (ClickHouse, MySQL, PostgreSQL, SQLite).
//
// Basic usage:
//
//	cfg := &olapsql.Configuration{
//	    ClientsOption: map[string]*olapsql.DBOption{
//	        "clickhouse": {DSN: "clickhouse://localhost:9000/default", Type: "clickhouse"},
//	    },
//	    DictionaryOption: &olapsql.Option{AdapterOption: olapsql.AdapterOption{Dsn: "olap-sql.toml"}},
//	}
//	manager, err := olapsql.NewManager(cfg)
//	result, err := manager.RunSync(query)
package olapsql

import (
	"fmt"

	"github.com/awatercolorpen/olap-sql/api/types"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Manager is the main entry point for olap-sql.
// It holds a set of database clients and an OLAP dictionary,
// and exposes methods to build or run OLAP queries.
type Manager struct {
	clients    Clients
	dictionary *Dictionary
}

// GetClients returns the registered database clients.
// Returns an error if the Manager was not initialised with a ClientsOption.
func (m *Manager) GetClients() (Clients, error) {
	if m.clients == nil {
		return nil, fmt.Errorf("it is no initialization")
	}
	return m.clients, nil
}

// GetDictionary returns the OLAP dictionary used for schema translation.
// Returns an error if the Manager was not initialised with a DictionaryOption.
func (m *Manager) GetDictionary() (*Dictionary, error) {
	if m.dictionary == nil {
		return nil, fmt.Errorf("it is no initialization")
	}
	return m.dictionary, nil
}

// SetLogger attaches a custom GORM logger to all registered database clients.
// Call this after creating the Manager to enable query logging.
func (m *Manager) SetLogger(log logger.Interface) {
	c, err := m.GetClients()
	if err == nil {
		c.SetLogger(log)
	}
}

// RunSync executes the OLAP query synchronously and returns the result.
// It translates the query into SQL, runs it against the target database,
// and returns a structured Result containing dimensions and row data.
func (m *Manager) RunSync(query *types.Query) (*types.Result, error) {
	db, err := m.BuildTransaction(query)
	if err != nil {
		return nil, err
	}
	rows, err := RunSync(db)
	if err != nil {
		return nil, err
	}
	return BuildResultSync(query, rows)
}

// RunChan executes the OLAP query and streams rows over a channel.
// This is useful for large result sets where you want to process rows
// as they arrive rather than buffering them all in memory.
func (m *Manager) RunChan(query *types.Query) (*types.Result, error) {
	db, err := m.BuildTransaction(query)
	if err != nil {
		return nil, err
	}
	rows, err := RunChan(db)
	if err != nil {
		return nil, err
	}
	return BuildResultChan(query, rows)
}

// BuildTransaction translates a Query into a *gorm.DB ready to execute.
// Use this when you need direct access to the GORM DB object,
// for example to add custom GORM hooks or inspect the generated SQL.
func (m *Manager) BuildTransaction(query *types.Query) (*gorm.DB, error) {
	clients, clause, err := m.build(query)
	if err != nil {
		return nil, err
	}
	return clients.BuildDB(clause)
}

// BuildSQL translates a Query into its SQL string without executing it.
// Useful for debugging, logging, or displaying the generated SQL to users.
func (m *Manager) BuildSQL(query *types.Query) (string, error) {
	clients, clause, err := m.build(query)
	if err != nil {
		return "", err
	}
	return clients.BuildSQL(clause)
}

// build is the internal helper that resolves the dictionary and clients
// needed to translate and execute a query.
func (m *Manager) build(query *types.Query) (Clients, types.Clause, error) {
	query.TranslateTimeIntervalToFilter()
	dict, err := m.GetDictionary()
	if err != nil {
		return nil, nil, err
	}
	clause, err := dict.Translate(query)
	if err != nil {
		return nil, nil, err
	}
	clients, err := m.GetClients()
	if err != nil {
		return nil, nil, err
	}
	return clients, clause, nil
}

// NewManager creates and initialises a Manager from the provided Configuration.
// At least one of ClientsOption or DictionaryOption should be set.
// Returns an error if any client DSN is invalid or the dictionary file cannot be parsed.
func NewManager(configuration *Configuration) (*Manager, error) {
	m := &Manager{}
	if configuration.ClientsOption != nil {
		clients, err := NewClients(configuration.ClientsOption)
		if err != nil {
			return nil, err
		}
		m.clients = clients
	}
	if configuration.DictionaryOption != nil {
		dict, err := NewDictionary(configuration.DictionaryOption)
		if err != nil {
			return nil, err
		}
		m.dictionary = dict
	}
	return m, nil
}
