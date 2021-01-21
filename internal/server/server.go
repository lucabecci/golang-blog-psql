package server

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	v1 "github.com/lucabecci/golang-blog-psql.git/internal/server/v1"
)

// Server is a base server configuration.
type Server struct {
	server *http.Server
}

func New(port string) (*Server, error) {
	r := chi.NewRouter()

	r.Mount("/api/v1", v1.New())

	serv := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	server := Server{server: serv}

	return &server, nil
}

func (serv *Server) Close() error {
	return nil
}

//This func start the srv
func (serv *Server) Start() {
	log.Printf("Server running on http://localhost%s", serv.server.Addr)
	log.Fatal(serv.server.ListenAndServe())
}
