package olapsql

type Manager struct {
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
	return m, nil
}