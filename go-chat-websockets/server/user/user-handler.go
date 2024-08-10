package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type Handler struct {
	Service
}

func NewHandler(s Service) *Handler {
	return &Handler{
		s,
	}
}

type httpError struct {
	Err string `json:"error"`
}

func writeError(w http.ResponseWriter, status int, err error) {
	he := httpError{
		Err: err.Error(),
	}
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(&he)
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	logrus.WithFields(logrus.Fields{
		"Remote Host": r.RemoteAddr,
	}).Info("signup request")
	var cr *CreateUserReq = &CreateUserReq{}
	if err := json.NewDecoder(r.Body).Decode(cr); err != nil {
		if errors.Is(err, io.EOF) {
			writeError(w, http.StatusBadRequest, fmt.Errorf("empty request body %w", err))
		} else {
			writeError(w, http.StatusBadRequest, err)
		}

		return
	}
	resp, err := h.Service.CreateUser(r.Context(), cr)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&resp)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	loginUserReq := LoginUserReq{}
	if err := json.NewDecoder(r.Body).Decode(&loginUserReq); err != nil {
		if errors.Is(err, io.EOF) {
			writeError(w, http.StatusBadRequest, fmt.Errorf("empty request body %w", err))
		} else {
			writeError(w, http.StatusBadRequest, err)
		}
		return
	}
	logrus.WithFields(logrus.Fields{
		"Remote Host": r.RemoteAddr,
		"email":       loginUserReq.Email,
	}).Info("login request")
	resp, err := h.Service.Login(r.Context(), &loginUserReq)
	if err != nil {
		writeError(w, http.StatusUnauthorized, err)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    resp.accessToken,
		Path:     "/",
		Domain:   "localhost",
		Expires:  time.Now().Add(time.Minute * 15),
		Secure:   false,
		HttpOnly: true,
	})
	json.NewEncoder(w).Encode(resp)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {

	http.SetCookie(w, &http.Cookie{
		Name:   "jwt",
		Value:  "",
		MaxAge: -1,
	})
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "logged out"})
}
