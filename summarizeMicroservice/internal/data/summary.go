package data

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
)

type Summary struct {
	Content string `firestore:"content"`
}

type SummaryModel struct {
	client *firestore.Client
}

func (m SummaryModel) Insert(summary *Summary, date string) error {
	documentPath := fmt.Sprintf("summaries/%s", date)
	_, err := m.client.Doc(documentPath).Create(context.Background(), summary)
	if err != nil {
		return err
	}
	return nil
}
