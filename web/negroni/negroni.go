package negroni

import (
	"net/http"

	classico "github.com/codegangsta/negroni"
)

type Negroni struct {
	*classico.Negroni
}

func Sbagliato() *Negroni {
	n := classico.New()
	n.Use(NewRecovery())
	n.Use(NewLogger())
	n.Use(classico.NewStatic(http.Dir("public")))
	return &Negroni{n}
}
