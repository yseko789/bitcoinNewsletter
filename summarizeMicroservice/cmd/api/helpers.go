package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

type envelope map[string]any

func (app *application) writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	js = append(js, '\n')

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}

func (app *application) yesterdaysDateString() string {
	loc, _ := time.LoadLocation("America/New_York")
	yesterday := time.Now().In(loc).AddDate(0, 0, -1).Format("2006-01-02")
	return yesterday
}

func (app *application) dayBeforeYesterdayDateString() string {
	loc, _ := time.LoadLocation("America/New_York")
	beforeYesterday := time.Now().In(loc).AddDate(0, 0, -2).Format("2006-01-02")
	return beforeYesterday
}

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
type chatgptRequest struct {
	Model    string    `json:"model"`
	Messages []message `json:"messages"`
}
type choice struct {
	Message message `json:"message"`
}
type chatgptResponse struct {
	Choices []choice `json:"choices"`
}

func (app *application) chatGPTSummarize(newsURLs []string) (string, error) {
	url := "https://api.openai.com/v1/chat/completions"
	apikey := os.Getenv("chatgptkey")

	contentString := strings.Join(newsURLs[:], " ")
	messageReq := message{
		Role:    "user",
		Content: fmt.Sprintf("For each article, summarize the main points: %s", contentString),
	}
	messagesReq := []message{messageReq}
	chatgptRequest := chatgptRequest{
		Model:    "gpt-3.5-turbo",
		Messages: messagesReq,
	}

	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(chatgptRequest)

	client := &http.Client{}
	req, err := http.NewRequest("POST", url, payloadBuf)
	if err != nil {
		app.logger.Error(fmt.Sprintf("%v", err))
		return "", err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", apikey))

	resp, err := client.Do(req)
	if err != nil {
		app.logger.Error(fmt.Sprintf("%v", err))
		return "", err
	}

	chatgptResponse := chatgptResponse{}
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&chatgptResponse)
	if err != nil {
		app.logger.Error(fmt.Sprintf("%v", err))
		return "", err
	}
	responseContent := chatgptResponse.Choices[0].Message.Content
	return responseContent, nil
}

type article struct {
	Author string `json:"author"`
	Title  string `json:"title"`
	URL    string `json:"url"`
}
type newsResponse struct {
	Status       string    `json:"status"`
	TotalResults int       `json:"totalResults"`
	Articles     []article `json:"articles"`
}

func (app *application) getNews(fromDate string, toDate string) (*newsResponse, error) {
	url := "https://newsapi.org/v2/everything?q=bitcoin&searchIn=title&sortBy=popularity&language=en&from=fromDate&to=toDate&pageSize=5"

	newsapi_key := os.Getenv("newskey")

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		app.logger.Error(fmt.Sprintf("%v", err))
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", newsapi_key)

	resp, err := client.Do(req)
	if err != nil {
		app.logger.Error(fmt.Sprintf("%v", err))
		return nil, err
	}
	newsResponse := newsResponse{}
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&newsResponse)
	if err != nil {
		app.logger.Error(fmt.Sprintf("%v", err))
		return nil, err
	}
	return &newsResponse, nil
}
