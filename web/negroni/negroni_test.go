package negroni

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Negroni_Sbagliato(t *testing.T) {
	n := Sbagliato()
	assert.Equal(t, 3, len(n.Handlers()))
	assert.Equal(t, "*negroni.Recovery", reflect.TypeOf(n.Handlers()[0]).String())
	assert.Equal(t, "*negroni.Logger", reflect.TypeOf(n.Handlers()[1]).String())
	assert.Equal(t, "*negroni.Static", reflect.TypeOf(n.Handlers()[2]).String())
}
