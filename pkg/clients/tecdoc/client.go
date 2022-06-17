package tecdoc

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"tec-doc/internal/model"
	"time"
)

//Client интерфейс с методами для получения запчастей с TecDoc
type TecDocClient interface {
	GetAllParts(ctx *context.Context, ID []string) ([]model.Autopart, error)
}

type tecDocClient struct {
	http.Client
	baseURL string
}

func NewClient(baseURL string, timeout time.Duration) (*tecDocClient, error) {
	return &tecDocClient{
		Client:  http.Client{Timeout: timeout},
		baseURL: baseURL,
	}, nil
}

func (c *tecDocClient) GetAllParts(ctx *context.Context, ID []string) ([]model.Autopart, error) {
	req, err := http.NewRequest(http.MethodGet, c.baseURL, nil)
	if err != nil {
		return nil, fmt.Errorf("can't create new request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	//Здесь будет что-то, добавляющее в запрос ID запчастей

	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("can't get response: %v", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("can't read response")
	}

	defer resp.Body.Close()

	var parts []model.Autopart

	err = json.Unmarshal(body, &parts)
	if err != nil {
		return nil, fmt.Errorf("can't unmarshal body: %v", err)
	}

	return parts, nil
}