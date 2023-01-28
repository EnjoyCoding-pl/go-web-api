package trip

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	dto "go-web-api/features/trip/app/models"
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
	Get(id int, ctx context.Context) (*models.Trip, error)
	GetAll(sp *models.TripSearchParam, ctx context.Context) (*[]models.Trip, error)
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
	r.HandleFunc("/trips", t.getAllHandler).Methods("GET")
}

func (t *tripController) addHandler(w http.ResponseWriter, r *http.Request) {

	spanCtx, span := otel.Tracer(globals.TracerAppName).Start(r.Context(), "trip-add")
	defer span.End()

	var trip dto.TripDto

	if err := json.NewDecoder(r.Body).Decode(&trip); err != nil {
		protocols.BadRequest(w, fmt.Errorf("trip-controller: unable to decode body; %v", err), span)
		return
	}

	if err := use_cases.NewAddUseCase(t.storage).Execute(trip, spanCtx); err != nil {
		protocols.InternalServerError(w, fmt.Errorf("trip-controller: unable to add new trip; %v", err), span)
		return
	}

	protocols.NoContent(w, span)
}

func (t *tripController) updateHandler(w http.ResponseWriter, r *http.Request) {

	spanCtx, span := otel.Tracer(globals.TracerAppName).Start(r.Context(), "trip-update")
	defer span.End()

	var trip dto.TripDto

	if err := json.NewDecoder(r.Body).Decode(&trip); err != nil {
		protocols.BadRequest(w, fmt.Errorf("trip-controller: unable to decode body; %v", err), span)
		return
	}

	if err := use_cases.NewUpdatetUseCase(t.storage).Execute(trip, spanCtx); err != nil {
		protocols.InternalServerError(w, fmt.Errorf("trip-controller: unable to update trip; %v", err), span)
		return
	}

	protocols.NoContent(w, span)
}

func (t *tripController) deleteHandler(w http.ResponseWriter, r *http.Request) {
	spanCtx, span := otel.Tracer(globals.TracerAppName).Start(r.Context(), "trip-delete")
	defer span.End()

	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		protocols.BadRequest(w, fmt.Errorf("trip-controller: unable to parse trip id; %v", err), span)
		return
	}

	if err := use_cases.NewDeleteUseCase(t.storage).Execute(id, spanCtx); err != nil {
		protocols.InternalServerError(w, fmt.Errorf("trip-controller: unable to delete trip; %v", err), span)
		return
	}

	protocols.NoContent(w, span)
}

func (t *tripController) getHandler(w http.ResponseWriter, r *http.Request) {
	spanCtx, span := otel.Tracer(globals.TracerAppName).Start(r.Context(), "trip-get")
	defer span.End()

	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		protocols.BadRequest(w, fmt.Errorf("trip-controller: unable to parse trip id; %v", err), span)
		return
	}

	resp, err := use_cases.NewGetUseCase(t.storage).Execute(id, spanCtx)

	if err != nil {
		protocols.InternalServerError(w, fmt.Errorf("trip-controller: unable to get trip %v; %v", id, err), span)
		return
	}

	protocols.Ok(w, &resp, span)
}

func (t *tripController) getAllHandler(w http.ResponseWriter, r *http.Request) {
	spanCtx, span := otel.Tracer(globals.TracerAppName).Start(r.Context(), "trip-get-all")
	defer span.End()

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		protocols.BadRequest(w, fmt.Errorf("trip-controller: unable to parse page; %v", err), span)
	}
	pageSize, err := strconv.Atoi(r.URL.Query().Get("pageSize"))
	if err != nil {
		protocols.BadRequest(w, fmt.Errorf("trip-controller: unable to parse page size; %v", err), span)
	}

	sp := models.NewTripSearchParam(page, pageSize)

	resp, err := use_cases.NewGetAllUseCase(t.storage).Execute(sp, spanCtx)

	if err != nil {
		protocols.InternalServerError(w, fmt.Errorf("trip-controller: unable to search trips; %v", err), span)
		return
	}
	protocols.Ok(w, resp, span)
}
