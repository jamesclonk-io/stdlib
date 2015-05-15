package negroni

import (
	"net/http"

	classico "github.com/codegangsta/negroni"
	"github.com/jamesclonk-io/stdlib/web"
)

type Negroni struct {
	*classico.Negroni
}

func Sbagliato() *Negroni {
	n := classico.New()
	n.Use(web.NewRecovery())
	n.Use(web.NewLogger())
	n.Use(classico.NewStatic(http.Dir("public")))
	return &Negroni{n}
}
