package redmine

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"mattermost/pkg/logger"

	querystring "github.com/google/go-querystring/query"
)

//pre-compiler time check.
// var _ RedmineInterface = (*Redmine)(nil)

type RedmineInterface interface {
	Issues() IssuesInterface
}

type Redmine struct {
	// HTTP Client to perform http operations.
	HTTPClient *http.Client

	// Redmine internal configurations.
	config *Config

	// Logger is alias for zap.Logger.
	// It can be used to add logging for all queries and methods.
	log *logger.Logger
}

// New Redmine API client.
func New(logger *logger.Logger, config *Config) *Redmine {
	return &Redmine{
		HTTPClient: http.DefaultClient,
		config:     config,
		log:        logger,
	}
}

func (c *Redmine) Request(ctx context.Context, method, path string,
	payload []byte,
	query, reply interface{}) (err error) {

	// TODO: init logger
	var uri string
	if uri, err = c.RequestURI(path, query); err != nil {
		return fmt.Errorf("unable to consttacty request uri")
	}

	// *** Perform HTTP Request ***
	var req *http.Request
	if req, err = http.NewRequestWithContext(ctx, method, uri, bytes.NewBuffer(payload)); err != nil {
		return fmt.Errorf("redmine_api: permorming http request error: %w", err)
	}
	// ***  Authentication section ***
	// Redmine auth by passing Redmine API Key as a Request HTTP Header as a "X-Redmine-API-Key" key.
	if c.config.APIKey != "" { //if config.APIKey ="" return
		req.Header.Add("X-Redmine-API-Key", c.config.APIKey)
	} else {
		return fmt.Errorf("redmine_api: missing Redmine API key")
	}

	if method == http.MethodPost || method == http.MethodPut || method == http.MethodDelete {
		req.Header.Add("Content-Type", "application/json")
	}

	// *** Sends an HTTP request and returns an HTTP response ***
	var resp *http.Response
	if resp, err = c.HTTPClient.Do(req); err != nil {
		return fmt.Errorf("redmine_api: repforming http request error: %w", err)
	}

	defer resp.Body.Close()

	// *** Read and check HTTP Status from Request Body ***
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusNoContent {
		var respErr string
		if respErr, err = c.decodeError(resp.Body); err != nil {
			return err
		}
		return fmt.Errorf("redmine_api: decoding http response error: %s", respErr)
	}

	// nothing to do, doesn't need to drain response body.
	if reply == nil {
		return
	}

	var response []byte
	if response, err = ioutil.ReadAll(resp.Body); err != nil {
		return fmt.Errorf("redmine_api: reading response error %w", err)
	}

	if err = json.Unmarshal(response, &reply); err != nil { //нужно передать указатель на реплай?
		return fmt.Errorf("redmine_api: decoding http response body error: %w", err)
	}

	return
}

func (c *Redmine) RequestURI(path string, query interface{}) (uri string, err error) {

	uri = join(c.config.Host, path)

	var values = make(url.Values)
	if query != nil {
		if values, err = querystring.Values(query); err != nil {
			return "", fmt.Errorf("redmine_api: encoding query: %w", err)
		}
	}

	if q := values.Encode(); q != "" {
		uri += "?" + q
	}

	return
}

// join method calculates uri used by host and path.
func join(host, path string) string {

	if strings.HasPrefix(host, "/") {
		if strings.HasPrefix(path, "/") {
			return host + path[1:]
		}
		return host + path
	}

	if strings.HasPrefix(path, "/") {
		return host + path
	}

	return host + path
}

// decodeError unmarshalling HTTP response error.
func (c *Redmine) decodeError(respBody io.Reader) (respErr string, err error) {

	// TODO: init logger

	type Error struct {
		Error string `json:"error"`
	}

	var errResponse []byte
	if errResponse, err = ioutil.ReadAll(respBody); err != nil {
		return "", fmt.Errorf("redmine_api: reading response body error: %w", err)
	}

	var exx Error
	if err = json.Unmarshal(errResponse, &exx); err != nil {
		return "", fmt.Errorf("redmine_api: decoding response body error: %w", err)
	}

	return exx.Error, nil
}
