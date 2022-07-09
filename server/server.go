package server

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"time"

	"github.com/arcticlimer/bookcatalog/business"
	"github.com/arcticlimer/bookcatalog/importers"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
)

//go:embed dist/assets/*
//go:embed dist/index.html
var dist embed.FS

type Server struct {
	Router *chi.Mux
	DB     *sqlx.DB
	err    error
}

const (
	port = 8080
)

func (a *Server) Initialize(db *sqlx.DB, libraryPath, coversPath string) {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	c := cors.New(cors.Options{
		AllowOriginFunc:  func(origin string) bool { return true },
		AllowCredentials: true,
		AllowedMethods:   []string{"HEAD", "GET", "POST", "DELETE"},
		// Enable Debugging for testing, consider disabling in production
		Debug: true,
	})

	var corsMiddleware = func(handler http.Handler) http.Handler {
		// Insert the middleware
		return c.Handler(handler)
	}

	r.Use(corsMiddleware)

	impcfg := importers.FsImporterConfig{
		LibraryPath: libraryPath,
		ImagesPath:  coversPath,
	}
	documentsRepo := business.NewDocumentsRepository(db)
	importer := importers.NewFsImporter(impcfg, documentsRepo)
	r.Mount("/api/documents", DocumentsResource{Repository: documentsRepo, Importer: importer}.Routes())

	// TODO: Use user provided image and library paths
	pdffs := http.FileServer(http.Dir(libraryPath))
	r.Handle("/library/*", http.StripPrefix("/library/", pdffs))

	coverfs := http.FileServer(http.Dir(coversPath))
	r.Handle("/covers/*", http.StripPrefix("/covers/", coverfs))

	subFS, _ := fs.Sub(dist, "dist")
	distfs := http.FileServer(http.FS(subFS))
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

func Start(db *sqlx.DB, addr, libraryPath, coversPath string) error {
	log.Println(fmt.Sprintf("starting server at %s", addr))
	a := Server{}
	a.Initialize(db, libraryPath, coversPath)
	return a.Run(addr)
}
