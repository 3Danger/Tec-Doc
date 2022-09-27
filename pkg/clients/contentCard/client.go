package contentCard

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"tec-doc/internal/tec-doc/config"
)

type ContentCardClient interface {
	Upload(supplierIdString string) (upload func(io.Reader) error, err error)
}

type contentCardClient struct {
	host   *url.URL
	client http.Client
}

type errorResponse struct {
	Data             string `json:"data"`
	Error            bool   `json:"error"`
	ErrorText        string `json:"errorText"`
	AdditionalErrors string `json:"additionalErrors"`
}

func New(cnf *config.ContentClientConfig) (ContentCardClient, error) {
	host, err := url.Parse(cnf.URL)
	if err != nil {
		return nil, err
	}
	return &contentCardClient{
		host:   host,
		client: http.Client{Timeout: cnf.Timeout},
	}, nil
}

func (c *contentCardClient) Upload(supplierIdString string) (upload func(io.Reader) error, err error) {
	request, err := http.NewRequest(http.MethodPost, c.host.JoinPath("/source/upload").String(), nil)
	if err != nil {
		return nil, err
	}
	request.AddCookie(&http.Cookie{Name: "x-supplier-id", Value: supplierIdString, HttpOnly: true})
	request.Header.Add("Content-Type", "application/json")
	return func(body io.Reader) error {
		var errResp errorResponse
		request.Body = io.NopCloser(body)
		if err := c.doRequest(request, &errResp); err != nil {
			return err
		}
		if errResp.Error {
			return fmt.Errorf("from server: %s", errResp.ErrorText)
		}
		return nil
	}, nil
}

// doRequest делает запрос и заполняет данными JSON структуру outStructPtr. аналог BindJSON() из gin
func (c *contentCardClient) doRequest(request *http.Request, outStructPtr interface{}) (err error) {
	var (
		response *http.Response
	)
	if response, err = c.client.Do(request); err != nil {
		return err
	}
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status code %d", response.StatusCode)
	}
	if err = json.NewDecoder(response.Body).Decode(&outStructPtr); err != nil {
		return err
	}
	return nil
}
