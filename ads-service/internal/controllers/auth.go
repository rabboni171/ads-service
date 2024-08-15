package controllers

import (
	authv1 "ads-service/api/auth-service/gen/proto"
	"ads-service/internal/clients/grpc/auth"
	"ads-service/internal/config"
	"ads-service/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

type authController struct {
	ctx        context.Context
	mux        *http.ServeMux
	authClient *auth.Client
}

func InitAuthController(
	ctx context.Context,
	mux *http.ServeMux,
	authClient *auth.Client,
) *authController {
	authController := &authController{
		ctx:        ctx,
		mux:        mux,
		authClient: authClient,
	}
	mux.HandleFunc("POST /user/login/", authController.Login)
	mux.HandleFunc("POST /user/register/", authController.Register)

	return authController
}

// @Summary Аутентификация пользователя
// @Tags auth
// @Accept json
// @Produce json
// @Param user body models.User true "Почта и пароль пользователя"
// @Success 200 {string} string "token" 
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router	/user/login/ [post]
func (c *authController) Login(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		slog.Error("unable to decode ad", "err", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	loginIn := &authv1.LoginRequest{
		AppId:    config.Cfg.AuthGPRC.AppId,
		Email:    user.Email,
		Password: user.Password,
	}
	res, err := c.authClient.Api.Login(c.ctx, loginIn)
	if err != nil {
		slog.Error("unable to login", "err", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{"token": "%s"}`, res.GetToken())))
}

// @Summary Регистрация пользователя
// @Tags auth
// @Accept json
// @Produce json
// @Param user body models.User true "Почта и пароль пользователя"
// @Success 201 {string} int "id" 
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router	/user/register/ [post]
func (c *authController) Register(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		slog.Error("unable to decode ad", "err", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	registerIn := &authv1.RegisterRequest{
		Email:    user.Email,
		Password: user.Password,
	}
	res, err := c.authClient.Api.Register(c.ctx, registerIn)
	if err != nil {
		slog.Error("unable to register", "err", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf(`{"id": %d}`, res.GetUserId())))
}
