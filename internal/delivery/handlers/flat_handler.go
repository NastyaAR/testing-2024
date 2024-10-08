package handlers

import (
	"avito-test-task/internal/domain"
	"avito-test-task/pkg"
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"io"
	"net/http"
	"time"
)

type FlatHandler struct {
	uc        domain.FlatUsecase
	lg        *zap.Logger
	dbTimeout time.Duration
}

func NewFlatHandler(uc domain.FlatUsecase, timeout time.Duration, lg *zap.Logger) *FlatHandler {
	return &FlatHandler{
		uc:        uc,
		lg:        lg,
		dbTimeout: timeout,
	}
}

func (h *FlatHandler) Create(w http.ResponseWriter, r *http.Request) {
	var (
		respBody     []byte
		flatRequest  domain.CreateFlatRequest
		flatResponse domain.CreateFlatResponse
	)

	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.lg.Warn("flat handler: create error", zap.Error(err))
		respBody = CreateErrorResponse(r.Context(), ReadHTTPBodyError, ReadHTTPBodyMsg)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(respBody)
		return
	}
	err = json.Unmarshal(body, &flatRequest)
	if err != nil {
		h.lg.Warn("flat handler: create error", zap.Error(err))
		respBody = CreateErrorResponse(r.Context(), UnmarshalHTTPBodyError, UnmarshalHTTPBodyMsg)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(respBody)
		return
	}

	userID, err := pkg.ExtractPayloadFromToken(r.Header.Get("authorization"), "userID")
	if err != nil {
		h.lg.Warn("flat handler: create error: extract id", zap.Error(err))
		respBody = CreateErrorResponse(r.Context(), CreateFlatError, CreateFlatErrorMsg)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(respBody)
		return
	}
	userUuid, err := uuid.Parse(userID)
	if err != nil {
		h.lg.Warn("flat handler: create error: extract id", zap.Error(err))
		respBody = CreateErrorResponse(r.Context(), CreateFlatError, CreateFlatErrorMsg)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(respBody)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), h.dbTimeout*time.Second)
	defer cancel()

	flatResponse, err = h.uc.Create(ctx, userUuid, &flatRequest, h.lg)
	if err != nil {
		h.lg.Warn("flat handler: create error", zap.Error(err))
		respBody = CreateErrorResponse(r.Context(), CreateFlatError, CreateFlatErrorMsg)
		w.WriteHeader(GetReturnHTTPCode(w, err))
		w.Write(respBody)
		return
	}

	respBody, err = json.Marshal(flatResponse)
	if err != nil {
		h.lg.Warn("flat handler: create error", zap.Error(err))
		respBody = CreateErrorResponse(r.Context(), MarshalHTTPBodyError, MarshalHTTPBodyErrorMsg)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(respBody)
		return
	}

	w.Write(respBody)
}

func (h *FlatHandler) Update(w http.ResponseWriter, r *http.Request) {
	var (
		respBody     []byte
		flatRequest  domain.UpdateFlatRequest
		flatResponse domain.CreateFlatResponse
	)

	defer r.Body.Close()

	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.lg.Warn("flat handler: update error", zap.Error(err))
		respBody = CreateErrorResponse(r.Context(), ReadHTTPBodyError, ReadHTTPBodyMsg)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(respBody)
		return
	}
	err = json.Unmarshal(body, &flatRequest)
	if err != nil {
		h.lg.Warn("flat handler: update error", zap.Error(err))
		respBody = CreateErrorResponse(r.Context(), UnmarshalHTTPBodyError, UnmarshalHTTPBodyMsg)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(respBody)
		return
	}

	userID, err := pkg.ExtractPayloadFromToken(r.Header.Get("authorization"), "userID")
	if err != nil {
		h.lg.Warn("flat handler: create error", zap.Error(err))
		respBody = CreateErrorResponse(r.Context(), CreateFlatError, CreateFlatErrorMsg)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(respBody)
		return
	}
	userUuid, err := uuid.Parse(userID)
	if err != nil {
		h.lg.Warn("flat handler: create error: extract id", zap.Error(err))
		respBody = CreateErrorResponse(r.Context(), CreateFlatError, CreateFlatErrorMsg)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(respBody)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), h.dbTimeout*time.Second)
	defer cancel()

	flatResponse, err = h.uc.Update(ctx, userUuid, &flatRequest, h.lg)
	if err != nil {
		h.lg.Warn("flat handler: update error", zap.Error(err))
		respBody = CreateErrorResponse(r.Context(), UpdateFlatError, UpdateFlatErrorMsg)
		w.WriteHeader(GetReturnHTTPCode(w, err))
		w.Write(respBody)
		return
	}

	respBody, err = json.Marshal(flatResponse)
	if err != nil {
		h.lg.Warn("flat handler: update error", zap.Error(err))
		respBody = CreateErrorResponse(r.Context(), MarshalHTTPBodyError, MarshalHTTPBodyErrorMsg)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(respBody)
		return
	}

	w.Write(respBody)
}
