package main

import (
	"fmt"
	"github.com/masiuchi/govn"
	"net/http"
)

func viewHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	fmt.Fprintf(w, "<html><body><div>Hello, %s</div></body></html>", r.URL.Path[1:])
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
