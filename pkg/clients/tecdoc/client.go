package tecdoc

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"tec-doc/internal/config"
	"tec-doc/internal/model"

	"github.com/rs/zerolog"
)

//Client интерфейс с методами для получения запчастей с TecDoc
type Client interface {
	GetAllParts(ID []string) ([]model.Autopart, error)
}

type client struct {
	http.Client
	cfg *config.Config
	log *zerolog.Logger
}

func NewClient(cfg *config.Config, log *zerolog.Logger) (*client, error) {
	return &client{
		Client: http.Client{Timeout: cfg.TecDocConfig.Timeout},
		cfg:    cfg,
		log:    log,
	}, nil
}

func (c *client) GetAllParts(ID []string) ([]model.Autopart, error) {
	req, err := http.NewRequest(http.MethodGet, c.cfg.TecDocConfig.Url, nil)
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
