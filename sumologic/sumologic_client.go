package sumologic

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/hashicorp/go-retryablehttp"
)

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Client struct {
	AccessID      string
	AccessKey     string
	AuthJwt       string
	Environment   string
	BaseURL       *url.URL
	IsInAdminMode bool
	httpClient    HttpClient
}

var ProviderVersion string

var endpoints = map[string]string{
	"us1": "https://api.sumologic.com/api/",
	"us2": "https://api.us2.sumologic.com/api/",
	"fed": "https://api.fed.sumologic.com/api/",
	"eu":  "https://api.eu.sumologic.com/api/",
	"au":  "https://api.au.sumologic.com/api/",
	"de":  "https://api.de.sumologic.com/api/",
	"jp":  "https://api.jp.sumologic.com/api/",
	"ca":  "https://api.ca.sumologic.com/api/",
	"in":  "https://api.in.sumologic.com/api/",
}

var rateLimiter = time.NewTicker(time.Minute / 240)

func createNewRequest(method, url string, body io.Reader, accessID string, accessKey string, authJwt string) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", "SumoLogicTerraformProvider/"+ProviderVersion)
	if authJwt == "" {
		req.SetBasicAuth(accessID, accessKey)
	} else {
		req.Header.Add("Authorization", "Bearer "+authJwt)
	}
	return req, nil
}

func logRequestAndResponse(req *http.Request, resp *http.Response) {
	var maskedHeader = req.Header.Clone()
	maskedHeader.Set("Authorization", "xxxxxxxxxxx")
	log.Printf("[DEBUG] Request: [Method=%s] [URL=%s] [Headers=%s]. Response: [Status=%s]\n", req.Method, req.URL, maskedHeader, resp.Status)
}

func (s *Client) PostWithCookies(urlPath string, payload interface{}) ([]byte, []*http.Cookie, error) {
	relativeURL, err := url.Parse(urlPath)
	if err != nil {
		return nil, nil, err
	}

	sumoURL := s.BaseURL.ResolveReference(relativeURL)

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, nil, err
	}

	req, err := createNewRequest(http.MethodPost, sumoURL.String(), bytes.NewBuffer(body), s.AccessID, s.AccessKey, s.AuthJwt)
	if err != nil {
		return nil, nil, err
	}

	<-rateLimiter.C
	resp, err := s.httpClient.Do(req)

	if err != nil {
		return nil, nil, err
	}
	logRequestAndResponse(req, resp)
	defer resp.Body.Close()

	respCookie := resp.Cookies()

	d, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, nil, errors.New(string(d))
	}

	return d, respCookie, nil
}

func (s *Client) GetWithCookies(urlPath string, cookies []*http.Cookie) ([]byte, string, error) {
	relativeURL, err := url.Parse(urlPath)
	if err != nil {
		return nil, "", err
	}

	sumoURL := s.BaseURL.ResolveReference(relativeURL)

	req, err := createNewRequest(http.MethodGet, sumoURL.String(), nil, s.AccessID, s.AccessKey, s.AuthJwt)
	if err != nil {
		return nil, "", err
	}

	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	<-rateLimiter.C
	resp, err := s.httpClient.Do(req)

	if err != nil {
		return nil, "", err
	}
	logRequestAndResponse(req, resp)
	defer resp.Body.Close()

	d, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}

	if resp.StatusCode == 404 {
		return nil, "", nil
	} else if resp.StatusCode >= 400 {
		return nil, "", errors.New(string(d))
	}

	return d, resp.Header.Get("ETag"), nil
}

func (s *Client) Post(urlPath string, payload interface{}) ([]byte, error) {
	relativeURL, _ := url.Parse(urlPath)
	sumoURL := s.BaseURL.ResolveReference(relativeURL)

	body, _ := json.Marshal(payload)
	req, err := createNewRequest(http.MethodPost, sumoURL.String(), bytes.NewBuffer(body), s.AccessID, s.AccessKey, s.AuthJwt)
	if err != nil {
		return nil, err
	}

	if s.IsInAdminMode {
		req.Header.Add("isAdminMode", "true")
	}

	<-rateLimiter.C
	resp, err := s.httpClient.Do(req)

	if err != nil {
		return nil, err
	}
	logRequestAndResponse(req, resp)
	defer resp.Body.Close()

	d, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, errors.New(string(d))
	}

	return d, nil
}

func (s *Client) PostRawPayload(urlPath string, payload string) ([]byte, error) {
	relativeURL, _ := url.Parse(urlPath)
	sumoURL := s.BaseURL.ResolveReference(relativeURL)
	req, err := createNewRequest(http.MethodPost, sumoURL.String(), bytes.NewBuffer([]byte(payload)), s.AccessID, s.AccessKey, s.AuthJwt)
	if err != nil {
		return nil, err
	}

	if s.IsInAdminMode {
		req.Header.Add("isAdminMode", "true")
	}

	<-rateLimiter.C
	resp, err := s.httpClient.Do(req)

	if err != nil {
		return nil, err
	}
	logRequestAndResponse(req, resp)

	d, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode >= 400 {
		return nil, errors.New(string(d))
	}

	return d, nil
}

func (s *Client) Put(urlPath string, payload interface{}) ([]byte, error) {
	relativeURL, _ := url.Parse(urlPath)
	sumoURL := s.BaseURL.ResolveReference(relativeURL)

	_, etag, _ := s.Get(sumoURL.String())

	body, _ := json.Marshal(payload)

	req, err := createNewRequest(http.MethodPut, sumoURL.String(), bytes.NewBuffer(body), s.AccessID, s.AccessKey, s.AuthJwt)
	if err != nil {
		return nil, err
	}
	req.Header.Add("If-Match", etag)

	if s.IsInAdminMode {
		req.Header.Add("isAdminMode", "true")
	}

	<-rateLimiter.C
	resp, err := s.httpClient.Do(req)

	if err != nil {
		return nil, err
	}
	logRequestAndResponse(req, resp)
	defer resp.Body.Close()

	d, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, errors.New(string(d))
	}

	return d, nil
}

func (s *Client) Get(urlPath string) ([]byte, string, error) {
	return s.GetWithErrOpt(urlPath, false)
}

func (s *Client) GetWithErrOpt(urlPath string, return404Err bool) ([]byte, string, error) {
	relativeURL, _ := url.Parse(urlPath)
	sumoURL := s.BaseURL.ResolveReference(relativeURL)

	req, err := createNewRequest(http.MethodGet, sumoURL.String(), nil, s.AccessID, s.AccessKey, s.AuthJwt)
	if err != nil {
		return nil, "", err
	}

	if s.IsInAdminMode {
		req.Header.Add("isAdminMode", "true")
	}

	<-rateLimiter.C
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, "", err
	}
	logRequestAndResponse(req, resp)

	defer resp.Body.Close()

	d, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}

	if resp.StatusCode == 404 {
		if return404Err {
			return nil, "", errors.New(string(d))
		} else {
			return nil, "", nil
		}
	} else if resp.StatusCode >= 400 {
		return nil, "", errors.New(string(d))
	}

	return d, resp.Header.Get("ETag"), nil
}

func (s *Client) Delete(urlPath string) ([]byte, error) {
	relativeURL, _ := url.Parse(urlPath)
	sumoURL := s.BaseURL.ResolveReference(relativeURL)

	req, err := createNewRequest(http.MethodDelete, sumoURL.String(), nil, s.AccessID, s.AccessKey, s.AuthJwt)
	if err != nil {
		return nil, err
	}

	if s.IsInAdminMode {
		req.Header.Add("isAdminMode", "true")
	}

	<-rateLimiter.C
	resp, err := s.httpClient.Do(req)

	if err != nil {
		return nil, err
	}
	logRequestAndResponse(req, resp)
	defer resp.Body.Close()

	d, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, errors.New(string(d))
	}

	return d, nil
}

func ErrorHandler(resp *http.Response, err error, numTries int) (*http.Response, error) {
	log.Printf("[ERROR] Request %s failed after %d attempts with response: [%s]", resp.Request.URL, numTries, resp.Status)
	return resp, err
}

func NewClient(accessID, accessKey, authJwt, environment, base_url string, admin bool) (*Client, error) {
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 10
	// Disable DEBUG logs (https://github.com/hashicorp/go-retryablehttp/issues/31)
	retryClient.Logger = nil
	retryClient.ErrorHandler = ErrorHandler
	client := Client{
		AccessID:      accessID,
		AccessKey:     accessKey,
		AuthJwt:       authJwt,
		httpClient:    retryClient.StandardClient(),
		Environment:   environment,
		IsInAdminMode: admin,
	}

	if base_url == "" {

		endpoint, found := endpoints[client.Environment]
		if !found {
			return nil, fmt.Errorf("could not find endpoint for environment %s", environment)
		}
		base_url = endpoint
	}
	parsed, err := url.Parse(base_url)
	if err != nil {
		return nil, fmt.Errorf("failed to parse base_url %s", base_url)
	}
	client.BaseURL = parsed

	return &client, nil
}

func HasErrorCode(errorJsonStr string, errorCodeChoices []string) string {
	var apiError ApiError
	jsonErr := json.Unmarshal([]byte(errorJsonStr), &apiError)
	if jsonErr != nil {
		// when fail to unmarshal JSON, we should consider the errorCode is not found
		return ""
	}
	for i := range apiError.Errors {
		for j := range errorCodeChoices {
			if apiError.Errors[i].Code == errorCodeChoices[j] {
				return errorCodeChoices[j]
			}
		}
	}
	return ""
}

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Detail  string `json:"detail"`
}

// e.g.:
// {"id":"8UQOI-82VTR-YBQ8G","errors":[{"code":"not_implemented_yet","message":"Not implemented yet"}]}
// {"id":"RO4X1-BZW7P-Q8KJF","errors":[{"code":"api_not_enabled","message":"This API is not enabled for your organization."}]}
type ApiError struct {
	Id     string  `json:"id"`
	Errors []Error `json:"errors"`
}

type Status struct {
	Status        string `json:"status"`
	StatusMessage string `json:"statusMessage"`
	Error         Error  `json:"error"`
}

type FolderUpdate struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type JobId struct {
	ID string `json:"id,omitempty"`
}

type Folder struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ParentId    string `json:"parentId"`
	CreatedBy   string `json:"createdBy,omitempty"`
}

type Content struct {
	ID             string      `json:"id,omitempty"`
	Type           string      `json:"type"`
	Name           string      `json:"name"`
	Description    string      `json:"description"`
	Children       []Content   `json:"children,omitempty"`
	Search         interface{} `json:"search"`
	SearchSchedule interface{} `json:"searchSchedule"`
	Config         string      `json:"-"`
	ParentId       string      `json:"-"`
}

// Connection is used to describe a connection.
type Connection struct {
	ID                string    `json:"id,omitempty"`
	Type              string    `json:"type"`
	Name              string    `json:"name"`
	Description       string    `json:"description,omitempty"`
	URL               string    `json:"url"`
	Headers           []Headers `json:"headers,omitempty"`
	CustomHeaders     []Headers `json:"customHeaders,omitempty"`
	DefaultPayload    string    `json:"defaultPayload"`
	WebhookType       string    `json:"webhookType"`
	ConnectionSubtype string    `json:"connectionSubtype,omitempty"`
	ResolutionPayload string    `json:"resolutionPayload,omitempty"`
}

// Headers is used to describe headers for http requests.
type Headers struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
