package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type APIServer struct {
	listenAddr string
	db         *gorm.DB
}

func NewApiServer(listenAddr string, db *gorm.DB) *APIServer {
	return &APIServer{listenAddr: listenAddr, db: db}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	//user service
	userService := NewUserHandler(s.db)
	userService.RegisterRoutes(subrouter)

	//url service
	urlService := NewUrlHandler(s.db)
	urlService.RegisterRoutes(subrouter)
	log.Println("Application is listening on port", s.listenAddr)
	http.ListenAndServe(s.listenAddr, router)
}
