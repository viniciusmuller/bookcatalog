package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
)

type Server struct {
	Router *chi.Mux
	DB     *sqlx.DB
	err    error
}

const (
	port = 8080
)

func (a *Server) Initialize(db *sqlx.DB) {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// booksRepo := NewBooksRepository(a.DB)
	// r.Mount("/books", Books{Repository: usersRepo}.Routes())

	a.Router = r
}

func (a *Server) Run(addr string) error {
	if a.err != nil {
		return a.err
	}
	log.Println(fmt.Sprintf("Listening at %s", addr))
	twoMinutes := time.Minute * 2
	srv := &http.Server{ReadTimeout: twoMinutes, WriteTimeout: twoMinutes,
		IdleTimeout: twoMinutes, Addr: addr, Handler: a.Router}
	return srv.ListenAndServe()
}

func Start(db *sqlx.DB, addr string) error {
	log.Println(fmt.Sprintf("starting server at %s", addr))
	a := Server{}
	a.Initialize(db)
	return a.Run(addr)
}
