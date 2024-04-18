package data

import (
	"errors"

	"cloud.google.com/go/firestore"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Models struct {
	Summary SummaryModel
}

func NewModels(client *firestore.Client) Models {
	return Models{
		Summary: SummaryModel{client: client},
	}
}
