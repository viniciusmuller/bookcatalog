package server

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/arcticlimer/bookcatalog/business"
	"github.com/arcticlimer/bookcatalog/importers"
	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"
	renderPkg "github.com/unrolled/render"
)

var render = renderPkg.New()

type DocumentsResource struct {
	Repository business.DocumentsRepository // TODO: Use repositorier interface here
	Importer   importers.FsImporter
}

func (rs DocumentsResource) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", rs.List)    // GET /todos - read a list of todos
	r.Post("/", rs.Create) // POST /todos - create a new todo and persist it

	r.Route("/{id}", func(r chi.Router) {
		r.Get("/", rs.Get) // GET /todos/{id} - read a single todo by :id
		// r.Put("/", rs.Update)    // PUT /todos/{id} - update a single todo by :id
		r.Delete("/", rs.Delete) // DELETE /todos/{id} - delete a single todo by :id
	})

	return r
}

func (rs DocumentsResource) List(w http.ResponseWriter, r *http.Request) {
	docs, err := rs.Repository.ListDocuments()
	if err != nil {
		log.Println(err)
		// writeStatus(w, http.StatusInternalServerError)
		return
	}

	err = render.JSON(w, http.StatusOK, docs)
	if err != nil {
		fmt.Println(err)
	}
}

func (rs DocumentsResource) Create(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(0)
	if err != nil {
		log.Println(err)
		return
	}

	file, header, err := r.FormFile("document")
	if err != nil {
		// TODO: Return these as JSON errors to the client
		log.Println(err)
		return
	}
	defer file.Close()

	err = rs.Importer.ImportFile(header.Filename, file)
	if err != nil {
		log.Println(err)
		// writeStatus(w, http.StatusInternalServerError)
		return
	}

	doc, err := rs.Repository.CreateDocument(header.Filename)
	if err != nil {
		log.Println(err)
		// writeStatus(w, http.StatusInternalServerError)
		return
	}

	err = render.JSON(w, http.StatusCreated, doc)
	if err != nil {
		fmt.Println(err)
	}
}

func (rs DocumentsResource) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	doc, err := rs.Repository.ReadDocument(id)
	if err != nil {
		if errors.Is(err, business.ErrDocumentNotFound) {
			// writeStatus(w, http.StatusNotFound)
			return
		}

		log.Println(err)
		// writeStatus(w, http.StatusInternalServerError)
		return
	}
	err = render.JSON(w, http.StatusOK, doc)
	if err != nil {
		fmt.Println(err)
	}
}

func (rs DocumentsResource) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	err := rs.Repository.DeleteDocument(id)
	if err != nil {
		if errors.Is(err, business.ErrDocumentNotFound) {
			// writeStatus(w, http.StatusNotFound)
			return
		}

		log.Println(err)
		// writeStatus(w, http.StatusInternalServerError)
	}
}
