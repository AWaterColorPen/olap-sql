package olapsql_test

import (
	"testing"

	"github.com/awatercolorpen/olap-sql/api/types"
	"github.com/stretchr/testify/assert"
)

func TestMockQuery2(t *testing.T) {
	testMockQuery(t, MockQuery2(), MockQuery2ResultAssert)
}

func TestMockQuery3(t *testing.T) {
	testMockQuery(t, MockQuery3(), MockQuery3ResultAssert)
}

func testMockQuery(t *testing.T, query *types.Query, check func(t assert.TestingT, result *types.Result)) {
	m, err := newManager(t)
	assert.NoError(t, err)
	assert.NoError(t, MockLoad(m))
	result, err := m.RunChan(query)
	assert.NoError(t, err)
	check(t, result)
}
