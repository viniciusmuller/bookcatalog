package business

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

var ErrDocumentNotFound = errors.New("document not found")

type Document struct {
  Id         string    `db:"id" json:"id"`
	Name       string    `db:"name" json:"name"`
	CoverUrl   string    `db:"cover_url" json:"coverUrl"`
	LibraryUrl string    `db:"library_url" json:"libraryUrl"`
	UpdatedAt  time.Time `db:"updated_at" json:"updatedAt"`
	InsertedAt time.Time `db:"inserted_at" json:"insertedAt"`
}

type DocumentsRepository struct {
	DB *sqlx.DB
}

type DocumentsRepositorier interface {
	CreateDocument(name string) (Document, error)
	// UpdateDocument(id string) (Document, error)
	ReadDocument(id string) (Document, error)
	DeleteDocument(id string) error
	ListDocuments() ([]Document, error)
}

func (d *DocumentsRepository) CreateDocument(name string) (Document, error) {
	uuid := uuid.New()
	doc := Document{
		Id:         uuid.String(),
		Name:       name,
		LibraryUrl: fmt.Sprintf("/library/%s", name),
		CoverUrl:   fmt.Sprintf("/covers/%s.jpg", removeExt(name)),
		InsertedAt: time.Now(),
		UpdatedAt:  time.Now(),
	}
	_, err := d.DB.NamedExec(`
    insert into Documents (id,name,library_url,cover_url,inserted_at,updated_at)
    values (:id,:name,:library_url,:cover_url,:inserted_at,:updated_at)`, doc)
	if err != nil {
		return Document{}, err
	}

	return doc, nil
}

// TODO: handle not found cases
func (d *DocumentsRepository) ReadDocument(id string) (Document, error) {
	doc := Document{}
	err := d.DB.Get(&doc,
		`select id,name,inserted_at,updated_at from Documents where id=$1`, id)
	if err != nil {
		return Document{}, fmt.Errorf("couldn't get document: %w", err)
	}

	return doc, nil
}

// TODO: update
// func (b *DocumentsRepository) UpdateDocument(id string) (Document, error) {
// 	return Document{}, nil
// }

func (d *DocumentsRepository) DeleteDocument(id string) error {
	_, err := d.DB.Exec("delete from Documents where id=$1", id)
	if err != nil {
		return fmt.Errorf("couldn't delete document: %w", err)
	}

	return nil
}

func (d *DocumentsRepository) ListDocuments() ([]Document, error) {
	docs := []Document{}
	err := d.DB.Select(&docs, "select id,name,library_url,cover_url,inserted_at,updated_at from Documents")
	if err != nil {
		return []Document{}, fmt.Errorf("couldn't select documents: %w", err)
	}

	return docs, nil
}

func NewDocumentsRepository(db *sqlx.DB) DocumentsRepository {
	return DocumentsRepository{DB: db}
}

func removeExt(fileName string) string {
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}
