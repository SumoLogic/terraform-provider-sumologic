package sumologic

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type SearchJobLink struct {
	Rel  string `json:"rel,omitempty"`
	Href string `json:"href,omitempty"`
}

type SearchJobResponse struct {
	SearchJobLink `json:"link"`
	ID            string `json:"id,omitempty"`
}

type SearchJobResults struct {
	Fields   map[string]*json.RawMessage `json:"omitempty"`
	Messages []SearchJobMessage          `json:"messages"`
}

type SearchJobMessage struct {
	Map SearchJobMessageMap `json:"map"`
}

type SearchJobMessageMap struct {
	BlockID        string `json:"_blockid"`
	MessageTime    string `json:"_messagetime"`
	RawData        string `json:"_raw"`
	CollectorID    string `json:"_collectorid"`
	SourceID       string `json:"_sourceid"`
	CollectorName  string `json:"_collector"`
	MessageCount   string `json:"_messagecount"`
	SourceHost     string `json:"_sourcehost"`
	MessageID      string `json:"_messageid"`
	SourceType     string `json:"_sourcename"`
	MessageSize    string `json:"_size"`
	ReceiptTime    string `json:"_receipttime"`
	SourceCategory string `json:"_sourcecategory"`
	Format         string `json:"_format"`
	SourceName     string `json:"_source"`
}

type SearchJobRequest struct {
	Query    string `json:"query"`
	From     string `json:"from"`
	To       string `json:"to"`
	TimeZone string `json:"timeZone"`
}

func (s *Client) createSearchJob(category string, searchTime time.Time) (string, []*http.Cookie, error) {

	request := SearchJobRequest{
		Query:    fmt.Sprintf("_sourceCategory=\"%s\" SubscriptionConfirmation", category),
		From:     searchTime.UTC().Add(time.Second * -80).Format(time.RFC3339),
		To:       searchTime.UTC().Add(time.Second * 80).Format(time.RFC3339),
		TimeZone: "Etc/UTC",
	}

	body, newCookie, err := s.PostWithCookies("search/jobs/", request)

	if err != nil {
		return "", newCookie, err
	}

	var response SearchJobResponse

	if err := json.Unmarshal(body, &response); err != nil {
		return "", newCookie, err
	}

	return response.ID, newCookie, nil
}

func (s *Client) searchForConfirmationMessage(jobID string, cookies []*http.Cookie) (string, error) {
	const waitForSearchJobRetries = 5
	// Wait for Search job
	if err := s.waitForSearchJob(jobID, cookies, waitForSearchJobRetries); err != nil {
		return "", err
	}

	urlPath := fmt.Sprintf("search/jobs/%v/messages?offset=0&limit=1", jobID)

	//ETag returned value ignored
	body, _, err := s.GetWithCookies(urlPath, cookies)

	if err != nil {
		return "", err
	}

	var response SearchJobResults

	if err := json.Unmarshal(body, &response); err != nil {
		return "", err
	}

	if len(response.Messages) == 0 {
		err = fmt.Errorf("Message not found")
		return "", err
	}

	return response.Messages[0].Map.RawData, nil
}

func (s *Client) createHTTPSourceSNSAutoConfirm(category string, now time.Time) error {

	// For creating SearchJob and use must use cookies
	jobID, cookies, err := s.createSearchJob(category, now)
	if err != nil {
		return err
	}

	message, err := s.searchForConfirmationMessage(jobID, cookies)

	if err != nil {
		return err
	}

	var schemaJSON struct {
		SubscribeURL string `json:"SubscribeURL"`
	}

	if err := json.Unmarshal([]byte(message), &schemaJSON); err != nil {
		return err
	}
	url := schemaJSON.SubscribeURL

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("failed to auto-confirm: %s", resp.Status) // resp.Status will be something like `403 Unauthorized`
	}

	if err != nil {
		return err
	}

	return nil
}

func (s *Client) waitForSearchJob(jobID string, cookies []*http.Cookie, retries int) error {
	var schemaJSON struct {
		State string `json:"state"`
	}

	urlPath := fmt.Sprintf("search/jobs/%v/", jobID)

	for state, cycle := schemaJSON.State, 0; state != "DONE GATHERING RESULTS"; state, cycle = schemaJSON.State, cycle+1 {
		if cycle > retries {
			return fmt.Errorf("Reached maximum retries of %v", retries)
		}

		body, _, err := s.GetWithCookies(urlPath, cookies)
		if err != nil {
			//Skip json parse
			continue
		}

		if err := json.Unmarshal(body, &schemaJSON); err != nil {
			return err
		}

		time.Sleep(time.Second)
	}

	return nil
}
