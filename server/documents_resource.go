package server

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/arcticlimer/bookcatalog/business"
	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"
	renderPkg "github.com/unrolled/render"
)

var render = renderPkg.New()

type CreateDocument struct {
	name string
}

type DocumentsResource struct {
	Repository business.DocumentsRepository // TODO: Use repositorier interface here
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
	render.JSON(w, http.StatusOK, docs)
}

func (rs DocumentsResource) Create(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(0)
	fmt.Println(r.FormValue("file"))

	doc, err := rs.Repository.CreateDocument("foobar")
	if err != nil {
		log.Println(err)
		// writeStatus(w, http.StatusInternalServerError)
		return
	}

	render.JSON(w, http.StatusCreated, doc)
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
	render.JSON(w, http.StatusOK, doc)
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
