package trip

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"go-web-api/features/trip/app/use_cases"
	"go-web-api/features/trip/domain/models"
	"go-web-api/internal/globals"
	"go-web-api/internal/protocols"

	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel"
)

type Storage interface {
	Add(t models.Trip, ctx context.Context) error
	Update(t models.Trip, ctx context.Context) error
	Delete(id int, ctx context.Context) error
	Get(id int, ctx context.Context) (models.Trip, error)
	GetAll(sp *models.TripSearchParam, ctx context.Context) (*[]models.Trip, error)
}

type tripController struct {
	storage Storage
	log     log.Logger
}

func NewTripController(s Storage, log log.Logger) *tripController {
	return &tripController{storage: s, log: log}
}

func (t *tripController) MuxRegister(r *mux.Router) {
	r.HandleFunc("/trips", t.addHandler).Methods("POST")
	r.HandleFunc("/trips", t.updateHandler).Methods("PUT")
	r.HandleFunc("/trips/{id}", t.deleteHandler).Methods("DELETE")
	r.HandleFunc("/trips/{id}", t.getHandler).Methods("GET")
	r.HandleFunc("/trips", t.getAllHandler).Methods("GET")
}

func (t *tripController) addHandler(w http.ResponseWriter, r *http.Request) {

	spanCtx, span := otel.Tracer(globals.TracerAppName).Start(r.Context(), "trip-add")
	defer span.End()

	var trip models.Trip

	if err := json.NewDecoder(r.Body).Decode(&trip); err != nil {
		t.log.Println(err)
		protocols.BadRequest(w)
		return
	}

	if err := use_cases.NewAddUseCase(t.storage).Execute(trip, spanCtx); err != nil {
		t.log.Println(err)
		protocols.InternalServerError(w)
		return
	}

	protocols.NoContent(w)
}

func (t *tripController) updateHandler(w http.ResponseWriter, r *http.Request) {

	spanCtx, span := otel.Tracer(globals.TracerAppName).Start(r.Context(), "trip-update")
	defer span.End()

	var trip models.Trip

	if err := json.NewDecoder(r.Body).Decode(&trip); err != nil {
		t.log.Println(err)
		protocols.BadRequest(w)
		return
	}

	if err := use_cases.NewUpdatetUseCase(t.storage).Execute(trip, spanCtx); err != nil {
		t.log.Println(err)
		protocols.InternalServerError(w)
		return
	}

	protocols.NoContent(w)
}

func (t *tripController) deleteHandler(w http.ResponseWriter, r *http.Request) {
	spanCtx, span := otel.Tracer(globals.TracerAppName).Start(r.Context(), "trip-delete")
	defer span.End()

	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		t.log.Println(err)
		protocols.BadRequest(w)
		return
	}

	if err := use_cases.NewDeleteUseCase(t.storage).Execute(id, spanCtx); err != nil {
		t.log.Println(err)
		protocols.InternalServerError(w)
		return
	}

	protocols.NoContent(w)
}

func (t *tripController) getHandler(w http.ResponseWriter, r *http.Request) {
	spanCtx, span := otel.Tracer(globals.TracerAppName).Start(r.Context(), "trip-get")
	defer span.End()

	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		t.log.Println(err)
		protocols.BadRequest(w)
		return
	}

	resp, err := use_cases.NewGetUseCase(t.storage).Execute(id, spanCtx)

	if err != nil {
		t.log.Println(err)
		protocols.InternalServerError(w)
		return
	}

	protocols.Ok(w, &resp)
}

func (t *tripController) getAllHandler(w http.ResponseWriter, r *http.Request) {
	spanCtx, span := otel.Tracer(globals.TracerAppName).Start(r.Context(), "trip-get-all")
	defer span.End()

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		protocols.BadRequest(w)
	}
	pageSize, err := strconv.Atoi(r.URL.Query().Get("pageSize"))
	if err != nil {
		protocols.BadRequest(w)
	}

	sp := models.NewTripSearchParam(page, pageSize)

	resp, err := use_cases.NewGetAllUseCase(t.storage).Execute(sp, spanCtx)

	if err != nil {
		t.log.Println(err)
		protocols.InternalServerError(w)
		return
	}
	protocols.Ok(w, resp)
}
