package middleware

import "net/http"

type ResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewResponseWriter(
	w http.ResponseWriter,
) *ResponseWriter {
	return &ResponseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
	}
}

func (rw *ResponseWriter) WriteHeader(
	statusCode int,
) {
	rw.statusCode = statusCode

	rw.ResponseWriter.WriteHeader(statusCode)
}

func (rw *ResponseWriter) Write(
	data []byte,
) (int, error) {
	return rw.ResponseWriter.Write(data)
}
