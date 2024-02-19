package glpi

import "github.com/go-gorp/gorp"

type Repository struct {
	db *gorp.DbMap
}

func NewRepository(cnn *gorp.DbMap) *Repository {
	return &Repository{
		db: cnn,
	}
}
