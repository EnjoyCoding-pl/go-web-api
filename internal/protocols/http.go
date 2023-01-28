package protocols

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

func Ok(w http.ResponseWriter, data interface{}, span trace.Span) {
	json, err := json.Marshal(data)

	if err != nil {
		InternalServerError(w, err, span)
	}
	span.SetStatus(codes.Ok, "")
	w.Header().Add("content-type", "application/json")
	w.Write(json)
}

func NoContent(w http.ResponseWriter, span trace.Span) {
	span.SetStatus(codes.Ok, "")
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

func BadRequest(w http.ResponseWriter, err error, span trace.Span) {
	log.Error(err)
	span.SetStatus(codes.Error, err.Error())
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
}

func NotFound(w http.ResponseWriter) {
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusNotFound)
}

func InternalServerError(w http.ResponseWriter, err error, span trace.Span) {
	log.Error(err)
	span.SetStatus(codes.Error, err.Error())
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
}

func Unauthorized(w http.ResponseWriter) {
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
}
