package main

import (
	"net/http"
	"time"

	"github.com/yseko789/bitcoinSummarize/internal/data"
)

func (app *application) insertSummaryHandler(w http.ResponseWriter, r *http.Request) {

	randomSummary := "hi this is bitcoin news for today. Bitcoin finally hit $100,000! To the moon!"
	summary := &data.Summary{
		Content: randomSummary,
	}
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	err := app.models.Summary.Insert(summary, yesterday)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusAccepted, envelope{"date": yesterday}, nil)
	if err != nil {
		app.logger.Error(err.Error())
	}

}

func (app *application) getLatestSummaryHandler(w http.ResponseWriter, r *http.Request) {
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	summary, err := app.models.Summary.GetByDate(yesterday)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"summary": summary}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}