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

func (m SummaryModel) GetByDate(date string) (*Summary, error) {

	summaries := m.client.Collection("summaries")
	docsnap, err := summaries.Doc(date).Get(context.Background())
	if err != nil {
		return nil, err
	}
	var summary Summary
	err = docsnap.DataTo(&summary)
	if err != nil {
		return nil, err
	}

	return &summary, nil
}
