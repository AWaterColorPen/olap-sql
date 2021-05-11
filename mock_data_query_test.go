package olapsql_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMockQuery2(t *testing.T) {
	m, err := newManager(t)
	assert.NoError(t, err)
	assert.NoError(t, MockLoad(m))

	query := MockQuery2()
	result, err := m.RunChan(query)
	assert.NoError(t, err)
	MockQuery2ResultAssert(t, result)
}

