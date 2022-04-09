package protocols

import (
	"encoding/json"
	"log"
	"net/http"
)

func Ok(w http.ResponseWriter, data interface{}) {
	json, err := json.Marshal(data)

	if err != nil {
		log.Println(err)
		InternalServerError(w)
	}
	w.Header().Add("content-type", "application/json")
	w.Write(json)
}

func NoContent(w http.ResponseWriter) {
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

func BadRequest(w http.ResponseWriter) {
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func NotFound(w http.ResponseWriter) {
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
}

func InternalServerError(w http.ResponseWriter) {
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
}
