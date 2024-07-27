package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/deepanshuemblinux/go-rate-limiter/service"
	"github.com/deepanshuemblinux/go-rate-limiter/tokenbucket"
)

type apiServer struct {
	listenAddr       string
	srvc             service.MessageService
	rate_limiter_map map[string]tokenbucket.TokenBucket
}

func NewAPIServer(listenAddr string, srvc service.MessageService) *apiServer {
	return &apiServer{
		listenAddr:       listenAddr,
		srvc:             srvc,
		rate_limiter_map: make(map[string]tokenbucket.TokenBucket, 0),
	}
}

func (s *apiServer) Run() {
	http.HandleFunc("/limited", s.handleLimited)
	http.HandleFunc("/unlimited", s.handleUnlimited)
	fmt.Printf("API Server listening on %s\n", s.listenAddr)
	err := http.ListenAndServe(s.listenAddr, nil)
	if err != nil {
		fmt.Println(err)
	}
}

func (s *apiServer) handleLimited(w http.ResponseWriter, r *http.Request) {
	ip := strings.Split(r.RemoteAddr, ":")[0]
	log.Printf("Request came from %s\n", ip)
	_, ok := s.rate_limiter_map[ip]
	if !ok {
		token_bucket := tokenbucket.NewTokenBucket(10)
		s.rate_limiter_map[ip] = token_bucket
		go s.rate_limiter_map[ip].StartPushing()
	}
	if !s.rate_limiter_map[ip].GetToken() {
		fmt.Println("Rejecting request")
		w.WriteHeader(http.StatusTooManyRequests)
		return
	}
	resp := s.srvc.GetMessage("Limited, don't over use me!")
	err := writeJSON(w, http.StatusOK, resp)
	if err != nil {
		log.Println(err)
	}
}

func (s *apiServer) handleUnlimited(w http.ResponseWriter, r *http.Request) {
	resp := s.srvc.GetMessage("Unlimited, Let's go!")
	err := writeJSON(w, http.StatusOK, resp)
	if err != nil {
		log.Println(err)
	}
}

func writeJSON(w http.ResponseWriter, status int, value any) error {
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(value)
	if err != nil {
		return err
	}
	return nil
}
