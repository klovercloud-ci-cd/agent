package v1

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Resource_Update(t *testing.T) {

	type TestCase struct {
		expected []string
		actual   []string
	}
	test:= TestCase{
		expected: []string{"abc"},
		actual:  []string{"abc"},
	}
	assert.ElementsMatch(t, test.expected, test.actual)
}