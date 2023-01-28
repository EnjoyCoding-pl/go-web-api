package user

import (
	"context"
	"encoding/json"
	"fmt"
	dto "go-web-api/features/user/app/models"
	"go-web-api/features/user/app/use_cases"
	"go-web-api/features/user/domain/models"
	"go-web-api/internal/globals"
	"go-web-api/internal/protocols"
	"net/http"

	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
)

type UserStorage interface {
	Add(u *models.User, ctx context.Context) error
	GetByLogin(login string, ctx context.Context) (*models.User, error)
}

type TokenProvider interface {
	Generate(userId int) (*string, error)
}

type userController struct {
	storage       UserStorage
	tokenProvider TokenProvider
}

func NewUserController(s UserStorage, p TokenProvider) *userController {
	return &userController{storage: s, tokenProvider: p}
}

func (c *userController) MuxRegister(r *mux.Router) {
	r.HandleFunc("/register", c.registerUser).Methods("POST")
	r.HandleFunc("/login", c.loginUser).Methods("POST")
}

func (c *userController) registerUser(w http.ResponseWriter, r *http.Request) {

	spanCtx, span := otel.Tracer(globals.TracerAppName).Start(r.Context(), "user-register")
	defer span.End()

	var dto dto.RegisterUserDto

	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		errMsg := fmt.Errorf("user-controller: unable to decode body; %v", err)
		protocols.BadRequest(w, errMsg, span)
		return
	}

	err := use_cases.NewRegisterUserUseCase(c.storage).Execute(dto, spanCtx)

	if err != nil {
		errMsg := fmt.Errorf("user-controller: unable to register user; %v", err)
		protocols.InternalServerError(w, errMsg, span)
		return
	}
	span.SetStatus(codes.Ok, "")
	protocols.NoContent(w, span)
}

func (c *userController) loginUser(w http.ResponseWriter, r *http.Request) {
	spanCtx, span := otel.Tracer(globals.TracerAppName).Start(r.Context(), "user-login")
	defer span.End()

	var dto dto.LoginUserDto

	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		errMsg := fmt.Errorf("user-controller: unable to decode body; %v", err)
		protocols.BadRequest(w, errMsg, span)
		return
	}

	token, err := use_cases.NewLoginUserUseCase(c.storage, c.tokenProvider).Execute(&dto, spanCtx)

	if err != nil {
		errMsg := fmt.Errorf("user-controller: unable to login user; %v", err)
		protocols.BadRequest(w, errMsg, span)
		return
	}
	span.SetStatus(codes.Ok, "")
	protocols.Ok(w, token, span)
}
