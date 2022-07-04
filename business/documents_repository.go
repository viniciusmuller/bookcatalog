package business

import (
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Document struct {
	Id         string    `db:"id"`
	Name       string    `db:"name"`
	UpdatedAt  time.Time `db:"updated_at"`
	InsertedAt time.Time `db:"inserted_at"`
}

type DocumentsRepository struct {
	DB *sqlx.DB
}

type DocumentsRepositories interface {
	CreateDocument(name string) (Document, error)
	UpdateDocument(id string) (Document, error)
	ReadDocument(id string) (Document, error)
	DeleteDocument(id string) error
}

func (b *DocumentsRepository) CreateDocument(name string) (Document, error) {
	uuid := uuid.New()
	doc := Document{
		Id:         uuid.String(),
		Name:       name,
		InsertedAt: time.Now(),
		UpdatedAt:  time.Now(),
	}
	_, err := b.DB.NamedExec(`
    insert into Documents (id,name,inserted_at,updated_at)
    values (:id,:name,:inserted_at,:updated_at)`, doc)
	if err != nil {
		return Document{}, err
	}

	return doc, nil
}

func (b *DocumentsRepository) ReadDocument(id string) (Document, error) {
	return Document{}, nil
}

func (b *DocumentsRepository) UpdateDocument(id string) (Document, error) {
	return Document{}, nil
}

func (b *DocumentsRepository) DeleteDocument(id string) (Document, error) {
	return Document{}, nil
}

func NewDocumentsRepository(db *sqlx.DB) DocumentsRepository {
	return DocumentsRepository{DB: db}
}
