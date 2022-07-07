package content

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"tec-doc/internal/config"
	"tec-doc/internal/model"
	"time"
)

type Client interface {
	CreateProductCard(ctx context.Context, contentConfig config.ContentClientConfig, cards []model.CreateProductCardRequest) error
}

type contentClient struct {
	http.Client
}

func NewClient(timeout time.Duration) *contentClient {
	return &contentClient{
		Client: http.Client{Timeout: timeout},
	}
}

func (c *contentClient) CreateProductCard(ctx context.Context, contentConfig config.ContentClientConfig, supplierID int, seviceUUID string, card model.CreateProductCardRequest) error {
	js, err := json.Marshal(card)
	if err != nil {
		return fmt.Errorf("can't marshal product card: %v", err)
	}
	reqBodyReader := bytes.NewReader(js)

	req, err := http.NewRequest(http.MethodPost, contentConfig.URL, reqBodyReader)
	if err != nil {
		return fmt.Errorf("can't create new request: %v", err)
	}

	req.Header = http.Header{"Content-Type": {"application/json"},
		"X-Int-Supplier-Id": {strconv.Itoa(supplierID)},
		"serviceUUID":       {seviceUUID}}

	resp, err := c.Do(req)
	if err != nil {
		return fmt.Errorf("can't get response: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("request failed with status code: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("can't read response")
	}

	type respStruct struct {
		Error     bool   `json:"error"`
		ErrorText string `json:"errorText"`
	}

	var r respStruct

	err = json.Unmarshal(body, &r)
	if r.Error == true {
		return fmt.Errorf("can't unmarshal body: %v", r.ErrorText)
	}

	return nil
}
