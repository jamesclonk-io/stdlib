package env

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Env_Get(t *testing.T) {
	assert.NotEqual(t, "foobar", Get("PATH", "foobar"))
	assert.Equal(t, "bar", Get("foo", "bar"))
}
