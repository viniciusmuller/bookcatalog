package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/arcticlimer/bookcatalog/business"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	"github.com/rs/cors"
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

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:8080"},
		AllowCredentials: true,
		// Enable Debugging for testing, consider disabling in production
		Debug: true,
	})

	var corsMiddleware = func(handler http.Handler) http.Handler {
		// Insert the middleware
		return c.Handler(handler)
	}

	r.Use(corsMiddleware)

	documentsRepo := business.NewDocumentsRepository(db)
	r.Mount("/api/documents", DocumentsResource{Repository: documentsRepo}.Routes())

	// TODO: Use user provided image and library paths
	pdffs := http.FileServer(http.Dir("library"))
	r.Handle("/library/*", http.StripPrefix("/library/", pdffs))

	coverfs := http.FileServer(http.Dir("img"))
	r.Handle("/covers/*", http.StripPrefix("/covers/", coverfs))

	distfs := http.FileServer(http.Dir("dist"))
	r.Handle("/*", http.StripPrefix("/", distfs))

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
