package govn

import (
	"errors"
	"net/http"
)

type WrapperWriter struct {
	http.ResponseWriter
	Body   string
	Status int
}

func NewWrapperWriter(w http.ResponseWriter) *WrapperWriter {
	wrapper := new(WrapperWriter)
	wrapper.ResponseWriter = w
	return wrapper
}

func (ww *WrapperWriter) Write(p []byte) (n int, err error) {
	ww.Body += string(p)
	return len(p), nil
}

func (ww *WrapperWriter) Flush() (n int, err error) {
	if ww.Body == "" {
		return 0, errors.New("no body")
	}
	return ww.ResponseWriter.Write([]byte(ww.Body))
}

func (ww *WrapperWriter) WriteHeader(s int) {
	ww.Status = s
	ww.ResponseWriter.WriteHeader(s)
}
