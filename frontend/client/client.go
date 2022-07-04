package client

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
	"net/url"
	"sync"
	"tec-doc/frontend/config"
	"time"
)

var once sync.Once
var client *Client

type Client struct {
	client       *http.Client
	frontPortURL *url.URL
	engine       *gin.Engine
	backendURL   *url.URL
}

func New(config *config.Config) *Client {
	var (
		back, front *url.URL
		err         error
	)

	if back, err = url.Parse("http://" + config.ServerAddress); err != nil {
		log.Error().Err(err).Send()
		return nil
	}

	if front, err = url.Parse("http://" + config.FrontendAddress); err != nil {
		log.Error().Err(err).Send()
		return nil
	}

	once.Do(func() {
		httpClient := &http.Client{Timeout: 10 * time.Second}
		client = &Client{
			client:       httpClient,
			backendURL:   back,
			engine:       gin.Default(),
			frontPortURL: front,
		}
		configureRouter(client)
	})
	return client
}

func (cl *Client) Run(ctx context.Context) error {
	errCh := make(chan error, 1)
	go func() {
		errCh <- cl.engine.Run(cl.frontPortURL.Host)
	}()
	select {
	case <-ctx.Done():
		return nil
	case e := <-errCh:
		return e
	}
}

func configureRouter(c *Client) {
	c.engine.Static("/css", "./frontend/templates/css")
	c.engine.Static("/js", "./frontend/templates/js")
	c.engine.LoadHTMLGlob("./frontend/templates/*.gohtml")

	c.engine.GET(frontMainPage, c.indexGet)
	c.engine.POST(frontMainPage, c.indexPost)
	c.engine.GET(frontExcelTemplate, c.downloadExcel)
}
