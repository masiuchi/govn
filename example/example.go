package main

import (
	"fmt"
	"github.com/masiuchi/govn"
	"net/http"
)

func viewHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", r.URL.Path[1:])
}

func main() {
	view := http.HandlerFunc(viewHandler)

	settings := govn.NewSettings()
	settings.UserToken = "IRb6-"
	settings.SecretKey = "secret"

	interceptor := govn.NewInterceptor(settings)

	http.Handle("/", interceptor.Call(view))
	http.ListenAndServe(":5000", nil)
}
