package olapsql

type Manager struct {
	clients     Clients
	dictionary *DataDictionary
}

func (m *Manager) Get() (interface{}, error) {
	return nil, nil
}

func NewManager(configuration *Configuration) (*Manager, error) {
	// set default olap-sql option
	if configuration.DefaultParallelNumber != 0 {
		DefaultParallelNumber = configuration.DefaultParallelNumber
	}

	m := &Manager{}
	if configuration.ClientsOption != nil {
		clients, err := NewClients(configuration.ClientsOption)
		if err != nil {
			return nil, err
		}
		m.clients = clients
	}
	if configuration.DataDictionaryOption != nil {
		dictionary, err := NewDataDictionary(configuration.DataDictionaryOption)
		if err != nil {
			return nil, err
		}
		m.dictionary = dictionary
	}
	return m, nil
}