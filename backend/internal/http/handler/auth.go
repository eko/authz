package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"sync"

	"github.com/eko/authz/backend/internal/entity/manager"
	"github.com/eko/authz/backend/internal/entity/model"
	"github.com/eko/authz/backend/internal/entity/repository"
	"github.com/eko/authz/backend/internal/security/jwt"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type AuthRequest struct {
	Username string `json:"username" example:"john.doe"`
	Password string `json:"password" example:"mypassword"`
}

type AuthResponse struct {
	AccessToken string      `json:"access_token"`
	TokenType   string      `json:"token_type"`
	ExpiresIn   int64       `json:"expires_in"`
	User        *model.User `json:"user"`
}

type TokenRequest struct {
	GrantType    string `json:"grant_type" example:"client_credentials"`
	ClientID     string `json:"client_id" example:"0be4e0e0-6788-4b99-8e00-e0af5b4945b1"`
	ClientSecret string `json:"client_secret" example:"EXCAdNZjCz0qJ_8uYA2clkxVdp_f1tm7"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token,omitempty"`
	ExpiresIn    int    `json:"expires_in"`
}

// Authenticates a user
//
//	@security	Authentication
//	@Summary	Authenticates a user
//	@Tags		Auth
//	@Produce	json
//	@Param		default	body		AuthRequest	true	"Authentication request"
//	@Success	200		{object}	AuthResponse
//	@Failure	400		{object}	model.ErrorResponse
//	@Failure	500		{object}	model.ErrorResponse
//	@Router		/v1/auth [Post]
func Authenticate(
	validate *validator.Validate,
	userManager manager.User,
	tokenManager jwt.Manager,
) fiber.Handler {
	return func(c *fiber.Ctx) error {
		request := &AuthRequest{}

		// Parse request body
		if err := c.BodyParser(request); err != nil {
			return returnError(c, http.StatusBadRequest, err)
		}

		// Validate body
		if err := validateStruct(validate, request); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(err)
		}

		user, err := userManager.GetRepository().GetByFields(map[string]repository.FieldValue{
			"username": {Operator: "=", Value: request.Username},
		})
		if err != nil {
			return returnError(c, http.StatusBadRequest, err)
		}

		if err := bcrypt.CompareHashAndPassword(
			[]byte(user.PasswordHash),
			[]byte(request.Password),
		); err != nil {
			return returnError(c, http.StatusBadRequest, errors.New("invalid credentials"))
		}

		token, err := tokenManager.Generate(user.Username)
		if err != nil {
			return returnError(c, http.StatusInternalServerError, errors.New("unable to generate token"))
		}

		return c.JSON(&AuthResponse{
			AccessToken: token.Token,
			TokenType:   token.TokenType,
			ExpiresIn:   token.ExpiresIn,
			User:        user,
		})
	}
}

// Retrieve a client token
//
//	@security	Authentication
//	@Summary	Retrieve a client token
//	@Tags		Auth
//	@Produce	json
//	@Param		default	body		TokenRequest	true	"Token request"
//	@Success	200		{object}	TokenResponse
//	@Failure	400		{object}	model.ErrorResponse
//	@Failure	500		{object}	model.ErrorResponse
//	@Router		/v1/token [Post]
func TokenNew(
	server *server.Server,
) http.HandlerFunc {
	mutex := &sync.Mutex{}

	return func(w http.ResponseWriter, r *http.Request) {
		mutex.Lock()
		server.SetResponseTokenHandler(responseTokenHandlerFunc(w, r))
		mutex.Unlock()

		err := server.HandleTokenRequest(w, r)
		if err != nil {
			returnHTTPError(w, http.StatusInternalServerError, err)
			return
		}
	}
}

func responseTokenHandlerFunc(w http.ResponseWriter, r *http.Request) server.ResponseTokenHandler {
	return func(w http.ResponseWriter, data map[string]any, header http.Header, statusCode ...int) error {
		w.Header().Set("Content-Type", "application/json;charset=UTF-8")
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Pragma", "no-cache")

		for key := range header {
			w.Header().Set(key, header.Get(key))
		}

		status := http.StatusOK
		if len(statusCode) > 0 && statusCode[0] > 0 {
			status = statusCode[0]
		}

		w.WriteHeader(status)

		if status >= 200 && status < 300 {
			return json.NewEncoder(w).Encode(data)
		}

		returnHTTPError(w, http.StatusBadRequest, errors.New("unable to handle oauth request"))
		return nil
	}
}
