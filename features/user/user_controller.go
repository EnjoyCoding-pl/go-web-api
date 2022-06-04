package user

import (
	"context"
	"encoding/json"
	dto "go-web-api/features/user/app/models"
	"go-web-api/features/user/app/use_cases"
	"go-web-api/features/user/domain/models"
	"go-web-api/internal/globals"
	"go-web-api/internal/protocols"
	"net/http"

	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel"
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
		protocols.BadRequest(w)
		return
	}

	err := use_cases.NewRegisterUserUseCase(c.storage).Execute(dto, spanCtx)

	if err != nil {
		protocols.InternalServerError(w)
		return
	}

	protocols.NoContent(w)
}

func (c *userController) loginUser(w http.ResponseWriter, r *http.Request) {
	spanCtx, span := otel.Tracer(globals.TracerAppName).Start(r.Context(), "user-login")
	defer span.End()

	var dto dto.LoginUserDto

	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		protocols.BadRequest(w)
		return
	}

	token, err := use_cases.NewLoginUserUseCase(c.storage, c.tokenProvider).Execute(&dto, spanCtx)

	if err != nil {
		protocols.BadRequest(w)
		return
	}

	protocols.Ok(w, token)
}
