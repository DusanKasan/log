package log

import (
	"bufio"
	"context"
	"io"
	"net/http"
	"net/url"
)

type response struct {
	Header http.Header
	Body   string
	Status int
}

type bufferedResponseWriter struct {
	w        http.ResponseWriter
	response response
}

func (b *bufferedResponseWriter) Header() http.Header {
	return b.w.Header()
}

func (b *bufferedResponseWriter) Write(bytes []byte) (int, error) {
	if len(b.response.Body) < 2048 {
		sliced := bytes
		if len(sliced) > 2048 {
			sliced = sliced[:2048]
		}

		b.response.Body = b.response.Body + string(sliced)
	}

	return b.w.Write(bytes)
}

func (b *bufferedResponseWriter) WriteHeader(status int) {
	b.response.Status = status
	b.w.WriteHeader(status)
}

func newBufferedResponse(w http.ResponseWriter) *bufferedResponseWriter {
	return &bufferedResponseWriter{w, response{Header: w.Header(), Status: 200}}
}

type request struct {
	Method      string
	Header      http.Header
	URL         url.URL
	Proto       string
	BodyPreview string
}

type bufferedReadCloser struct {
	closer io.Closer
	*bufio.Reader
}

func (b bufferedReadCloser) Close() error {
	return b.closer.Close()
}

func Middleware(next func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		requestID := generateRequestID()
		r = r.WithContext(context.WithValue(r.Context(), keyId, string(requestID)))
		w.Header().Set("Request-Id", string(requestID))

		buff := bufferedReadCloser{r.Body, bufio.NewReader(r.Body)}
		r.Body = buff
		b, err := buff.Peek(2048)
		if err != nil && err != io.EOF {
			Error(r.Context(), "unable to peek into buffered request body, cancelling request", request{r.Method, r.Header, *r.URL, r.Proto, ""})
		}

		Info(r.Context(), "request started", request{r.Method, r.Header, *r.URL, r.Proto, string(b)})

		br := newBufferedResponse(w)
		next(br, r)

		Info(r.Context(), "request finished", br.response)
	}
}
