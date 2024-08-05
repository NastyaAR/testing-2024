package handlers

import (
	"avito-test-task/internal/domain"
	"context"
	"encoding/json"
	"go.uber.org/zap"
	"io"
	"net/http"
	"time"
)

type UserHandler struct {
	uc        domain.UserUsecase
	lg        *zap.Logger
	dbTimeout time.Duration
}

func NewUserHandler(uc domain.UserUsecase, timeout time.Duration, lg *zap.Logger) *UserHandler {
	return &UserHandler{
		uc:        uc,
		lg:        lg,
		dbTimeout: timeout,
	}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var (
		respBody         []byte
		registerRequest  domain.RegisterUserRequest
		registerResponse domain.RegisterUserResponse
	)
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.lg.Warn("user handler: register error", zap.Error(err))
		respBody = CreateErrorResponse(r.Context(), ReadHTTPBodyError, ReadHTTPBodyMsg)
		w.Write(respBody)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(body, &registerRequest)
	if err != nil {
		h.lg.Warn("user handler: register error", zap.Error(err))
		respBody = CreateErrorResponse(r.Context(), UnmarshalHTTPBodyError, UnmarshalHTTPBodyMsg)
		w.Write(respBody)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	registerResponse, err = h.uc.Register(ctx, &registerRequest, h.lg)
	if err != nil {
		h.lg.Warn("user handler: register error", zap.Error(err))
		respBody = CreateErrorResponse(r.Context(), RegisterUserError, RegisterUserErrorMsg)
		w.Write(respBody)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	respBody, err = json.Marshal(registerResponse)
	if err != nil {
		h.lg.Warn("user handler: register error", zap.Error(err))
		respBody = CreateErrorResponse(r.Context(), MarshalHTTPBodyError, MarshalHTTPBodyErrorMsg)
		w.Write(respBody)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(respBody)
	w.WriteHeader(http.StatusOK)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var (
		respBody      []byte
		loginRequest  domain.LoginUserRequest
		loginResponse domain.LoginUserResponse
	)
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.lg.Warn("user handler: login error", zap.Error(err))
		respBody = CreateErrorResponse(r.Context(), ReadHTTPBodyError, ReadHTTPBodyMsg)
		w.Write(respBody)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(body, &loginRequest)
	if err != nil {
		h.lg.Warn("user handler: login error", zap.Error(err))
		respBody = CreateErrorResponse(r.Context(), UnmarshalHTTPBodyError, UnmarshalHTTPBodyMsg)
		w.Write(respBody)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	loginResponse, err = h.uc.Login(ctx, &loginRequest, h.lg)
	if err != nil {
		h.lg.Warn("user handler: login error", zap.Error(err))
		respBody = CreateErrorResponse(r.Context(), LoginUserError, LoginUserErrorMsg)
		w.Write(respBody)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	respBody, err = json.Marshal(loginResponse)
	if err != nil {
		h.lg.Warn("user handler: register error", zap.Error(err))
		respBody = CreateErrorResponse(r.Context(), MarshalHTTPBodyError, MarshalHTTPBodyErrorMsg)
		w.Write(respBody)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(respBody)
	w.WriteHeader(http.StatusOK)
}

func (h *UserHandler) DummyLogin(w http.ResponseWriter, r *http.Request) {
	var (
		respBody           []byte
		userType           string
		dummyLoginResponse domain.LoginUserResponse
	)
	defer r.Body.Close()

	userType = r.URL.Query().Get("user_type")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dummyLoginResponse, err := h.uc.DummyLogin(ctx, userType, h.lg)
	if err != nil {
		h.lg.Warn("user handler: dummy login error", zap.Error(err))
		respBody = CreateErrorResponse(r.Context(), DummyLoginError, DummyLoginErrorMsg)
		w.Write(respBody)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	respBody, err = json.Marshal(dummyLoginResponse)
	if err != nil {
		h.lg.Warn("user handler: register error", zap.Error(err))
		respBody = CreateErrorResponse(r.Context(), MarshalHTTPBodyError, MarshalHTTPBodyErrorMsg)
		w.Write(respBody)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(respBody)
	w.WriteHeader(http.StatusOK)
}