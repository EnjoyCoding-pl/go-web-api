package trip

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"go-web-api/features/trip/app/use_cases"
	"go-web-api/features/trip/domain/models"

	"github.com/gorilla/mux"
)

type Storage interface {
	Add(t models.Trip, ctx context.Context) error
	Update(t models.Trip, ctx context.Context) error
	Delete(id int, ctx context.Context) error
	Get(id int, ctx context.Context) (models.Trip, error)
}

type tripController struct {
	storage Storage
}

func NewTripController(s Storage) *tripController {
	return &tripController{storage: s}
}

func (t *tripController) MuxRegister(r *mux.Router) {
	r.HandleFunc("/trips", t.addHandler).Methods("POST")
	r.HandleFunc("/trips", t.updateHandler).Methods("PUT")
	r.HandleFunc("/trips/{id}", t.deleteHandler).Methods("DELETE")
	r.HandleFunc("/trips/{id}", t.getHandler).Methods("GET")
}

func (t *tripController) addHandler(w http.ResponseWriter, r *http.Request) {

	var trip models.Trip

	if err := json.NewDecoder(r.Body).Decode(&trip); err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}

	log.Println(trip)

	if err := use_cases.NewAddUseCase(t.storage).Execute(trip, r.Context()); err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(200)
}

func (t *tripController) updateHandler(w http.ResponseWriter, r *http.Request) {

	var trip models.Trip

	if err := json.NewDecoder(r.Body).Decode(&trip); err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}

	if err := use_cases.NewUpdatetUseCase(t.storage).Execute(trip, r.Context()); err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(200)
}

func (t *tripController) deleteHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}

	if err := use_cases.NewDeleteUseCase(t.storage).Execute(id, r.Context()); err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(200)
}

func (t *tripController) getHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}

	resp, err := use_cases.NewGetUseCase(t.storage).Execute(id, r.Context())

	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}

	json, err := json.Marshal(resp)

	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}

	w.Write(json)
}
