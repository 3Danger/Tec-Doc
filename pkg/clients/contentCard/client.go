package contentCard

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"tec-doc/internal/tec-doc/config"
)

type ContentCardClient interface {
	Upload(body io.Reader) (err error)
}

type contentCardClient struct {
	host string
	http.Client
	header http.Header
}

type errorResponse struct {
	Data             string `json:"data"`
	Error            bool   `json:"error"`
	ErrorText        string `json:"errorText"`
	AdditionalErrors string `json:"additionalErrors"`
}

func New(cnf *config.ContentClientConfig) (ContentCardClient, error) {
	var (
		err       error
		sourceUrl *url.URL
		cookie    *cookiejar.Jar
	)
	if sourceUrl, err = url.Parse(cnf.URL); err != nil {
		return nil, err
	}
	if cookie, err = cookiejar.New(nil); err != nil {
		return nil, err
	}
	cookie.SetCookies(sourceUrl, []*http.Cookie{
		{
			Name:     "x-supplier-id",
			Value:    cnf.SupplierId,
			HttpOnly: true,
		},
	})
	return &contentCardClient{
		host:   cnf.URL,
		header: map[string][]string{"Content-Type": {"application/json"}},
		Client: http.Client{Jar: cookie, Timeout: cnf.Timeout},
	}, nil
}

func (c *contentCardClient) Upload(body io.Reader) (err error) {
	var errResp errorResponse
	if err = c.doRequest(http.MethodPost, c.host+"/source/upload", body, &errResp); err != nil {
		return err
	}
	if errResp.Error {
		return fmt.Errorf("from server: %s", errResp.ErrorText)
	}
	return nil
}

// doRequest делает запрос и заполняет данными JSON структуру outStructPtr. аналог BindJSON() из gin
func (c *contentCardClient) doRequest(method, path string, body io.Reader, outStructPtr interface{}) (err error) {
	var (
		response *http.Response
		request  *http.Request
	)
	if request, err = http.NewRequest(method, path, body); err != nil {
		return fmt.Errorf("can't create new request: %w", err)
	}

	request.Header = c.header
	if response, err = c.Do(request); err != nil {
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
