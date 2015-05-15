package env

import (
	"testing"
)

func Test_Env_Get(t *testing.T) {
	if Get("PATH", "foobar") == "foobar" {
		t.Error("PATH exists, should not get back 'foobar' value")
	}

	// nvl test
	if Get("foo", "bar") != "bar" {
		t.Error("foo env var does not exist, should get back 'bar' value")
	}
}
