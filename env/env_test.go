package env

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Env_Get(t *testing.T) {
	assert.NotEqual(t, "foobar", Get("PATH", "foobar"), "PATH env var exists, value should not be 'foobar'")
	assert.Equal(t, "bar", Get("foo", "bar"), "foo env var does not exist, should give back 'bar' value")
}
