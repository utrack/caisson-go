package httpbinding

import "net/http"

type wrappedWriter struct {
	blockWrites bool
	written     bool
	w           http.ResponseWriter
}

func (w wrappedWriter) Header() http.Header {
	return w.w.Header()
}

func (w *wrappedWriter) Write(b []byte) (int, error) {
	if w.blockWrites {
		panic("direct writing from handler is not allowed - RPC handler has marshalled return type")
	}
	w.written = true
	return w.w.Write(b)
}

func (w *wrappedWriter) WriteHeader(statusCode int) {
	if w.blockWrites {
		panic("direct writing from handler is not allowed - RPC handler has marshalled return type")
	}
	w.w.WriteHeader(statusCode)
}
