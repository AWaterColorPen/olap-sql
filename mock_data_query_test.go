package olapsql_test

import (
	"testing"

	"github.com/awatercolorpen/olap-sql/api/types"
	"github.com/stretchr/testify/assert"
)

func TestMockQuery1(t *testing.T) {
	testMockQuery(t, MockQuery1(), MockQuery1ResultAssert)
}

func TestMockQuery2(t *testing.T) {
	testMockQuery(t, MockQuery2(), MockQuery2ResultAssert)
}

func TestMockQuery3(t *testing.T) {
	testMockQuery(t, MockQuery3(), MockQuery3ResultAssert)
}

func TestMockQuery4(t *testing.T) {
	testMockQuery(t, MockQuery4(), MockQuery4ResultAssert)
}

func TestMockQuery5(t *testing.T) {
	testMockQuery(t, MockQuery5(), MockQuery5ResultAssert)
}

func TestMockQuery6(t *testing.T) {
	testMockQuery(t, MockQuery6(), MockQuery6ResultAssert)
}

func TestMockQuery7(t *testing.T) {
	testMockQuery(t, MockQuery7(), MockQuery7ResultAssert)
}

func TestMockQuery8(t *testing.T) {
	testMockQuery(t, MockQuery8(), MockQuery8ResultAssert)
}

func TestMockQuery9(t *testing.T) {
	testMockQuery(t, MockQuery9(), MockQuery9ResultAssert)
}

func testMockQuery(t *testing.T, query *types.Query, check func(t assert.TestingT, result *types.Result)) {
	m, err := newManager(t)
	assert.NoError(t, err)
	assert.NoError(t, MockLoad(m))
	result, err := m.RunSync(query)
	assert.NoError(t, err)
	check(t, result)
}
