package mwencoding

import (
	"bufio"
	"github.com/urakozz/go-middleware-encoding/encoding"
	"github.com/ant0ine/go-json-rest/rest"
	"net"
	"net/http"
	"strings"
)

type Flusher interface {
	Flush() error
}

var _encNone = &encoding.EncoderNone{}

var encoders []encoding.Encoder

func RegisterEncoder(e encoding.Encoder) {
	encoders = append(encoders, e)
}

func SelectEncoder(accept string) (enc encoding.Encoder) {
	for _, e := range encoders {
		if strings.Contains(accept, e.GetName()) {
			enc = e
		}
	}
	if enc == nil {
		enc = _encNone
	}
	return
}

func init() {
	RegisterEncoder(&encoding.EncoderGzip{})
	RegisterEncoder(&encoding.EncoderDeflate{})
	//RegisterEncoder(&encoding.EncoderCBrotli{})
}

// GzipMiddleware is responsible for compressing the payload with gzip and setting the proper
// headers when supported by the client. It must be wrapped by TimerMiddleware for the
// compression time to be captured. And It must be wrapped by RecorderMiddleware for the
// compressed BYTES_WRITTEN to be captured.
type GzipMiddleware struct {
	Level int
}

// MiddlewareFunc makes GzipMiddleware implement the Middleware interface.
func (mw *GzipMiddleware) MiddlewareFunc(h rest.HandlerFunc) rest.HandlerFunc {
	return func(w rest.ResponseWriter, r *rest.Request) {

		encoder := SelectEncoder(r.Header.Get("Accept-Encoding"))
		if mw.Level == 0 {
			mw.Level = 1
		}
		// client accepts gzip ?
		writer := &gzipResponseWriter{w, false, encoder, mw.Level}
		// call the handler with the wrapped writer
		h(writer, r)
	}
}

// Private responseWriter intantiated by the gzip middleware.
// It encodes the payload with gzip and set the proper headers.
// It implements the following interfaces:
// ResponseWriter
// http.ResponseWriter
// http.Flusher
// http.CloseNotifier
// http.Hijacker
type gzipResponseWriter struct {
	rest.ResponseWriter
	wroteHeader   bool
	encoderWriter encoding.Encoder
	level         int
}

// Set the right headers for gzip encoded responses.
func (w *gzipResponseWriter) WriteHeader(code int) {

	// Always set the Vary header, even if this particular request
	// is not gzipped.
	w.Header().Add("Vary", "Accept-Encoding")

	if w.encoderWriter.GetName() != "" {
		w.Header().Set("Content-Encoding", w.encoderWriter.GetName())
	}

	w.ResponseWriter.WriteHeader(code)
	w.wroteHeader = true
}

// Make sure the local Write is called.
func (w *gzipResponseWriter) WriteJson(v interface{}) error {
	b, err := w.EncodeJson(v)
	if err != nil {
		return err
	}
	_, err = w.Write(b)
	if err != nil {
		return err
	}
	return nil
}

// Make sure the local WriteHeader is called, and call the parent Flush.
// Provided in order to implement the http.Flusher interface.
func (w *gzipResponseWriter) Flush() {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	flusher := w.ResponseWriter.(http.Flusher)
	flusher.Flush()
}

// Call the parent CloseNotify.
// Provided in order to implement the http.CloseNotifier interface.
func (w *gzipResponseWriter) CloseNotify() <-chan bool {
	notifier := w.ResponseWriter.(http.CloseNotifier)
	return notifier.CloseNotify()
}

// Provided in order to implement the http.Hijacker interface.
func (w *gzipResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hijacker := w.ResponseWriter.(http.Hijacker)
	return hijacker.Hijack()
}

// Make sure the local WriteHeader is called, and encode the payload if necessary.
// Provided in order to implement the http.ResponseWriter interface.
func (w *gzipResponseWriter) Write(b []byte) (int, error) {

	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}

	writer := w.ResponseWriter.(http.ResponseWriter)

	count, errW := w.encoderWriter.NewWriter(w, w.level).Write(b)
	var errF error
	if f, ok := w.encoderWriter.(Flusher); ok {
		errF = f.Flush()
	}
	if errW != nil {
		return count, errW
	}
	if errF != nil {
		return count, errF
	}
	return count, nil

	return writer.Write(b)
}

var _ rest.Middleware = (*GzipMiddleware)(nil)
