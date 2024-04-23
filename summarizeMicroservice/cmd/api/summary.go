package main

import (
	"net/http"

	"github.com/yseko789/bitcoinSummarize/internal/data"
)

func (app *application) insertSummaryHandler(w http.ResponseWriter, r *http.Request) {

	news, err := app.getNews(app.yesterdaysDateString(), app.dayBeforeYesterdayDateString())
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	var newsURLs []string
	for _, n := range news.Articles {
		newsURLs = append(newsURLs, n.URL)
	}

	summaryResponse, err := app.chatGPTSummarize(newsURLs)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	summary := &data.Summary{
		Content: summaryResponse,
	}

	err = app.models.Summary.Insert(summary, app.yesterdaysDateString())
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusAccepted, envelope{"date": app.yesterdaysDateString()}, nil)
	if err != nil {
		app.logger.Error(err.Error())
	}

}
