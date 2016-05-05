package main

import (
	"fmt"
	"net/http"

	"github.com/masiuchi/govn"
)

func viewHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	fmt.Fprintf(w, "<html><body><div>Hello, %s</div></body></html>", r.URL.Path[1:])
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	fmt.Fprintf(w, "<html><body><div>This is a test page.</div></body></html>")
}

func main() {
	view := http.HandlerFunc(viewHandler)
	test := http.HandlerFunc(testHandler)

	settings := govn.NewSettings()
	settings.UserToken = "IRb6-"
	settings.SecretKey = "secret"

	interceptor := govn.NewInterceptor(settings)

	http.Handle("/test/", interceptor.Call(test))
	http.Handle("/", interceptor.Call(view))
	http.ListenAndServe(":5000", nil)
}
