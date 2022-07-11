package client

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"tec-doc/internal/config"
	"time"
)

type Client struct {
	client        *http.Client
	frontPort     string
	engine        *gin.Engine
	backendAddres string
}

func New(config *config.Config) *Client {
	httpClient := &http.Client{Timeout: 10 * time.Second}
	if config.ServerPort != "" {
		config.ServerPort = ":" + config.ServerPort
	}
	client := &Client{
		client:        httpClient,
		backendAddres: "http://" + config.ServerHost + config.ServerPort,
		engine:        gin.Default(),
		frontPort:     config.FrontendPort,
	}
	configureRouter(client)
	return client
}

func (cl *Client) Run(ctx context.Context) error {
	errCh := make(chan error, 1)
	go func() {
		errCh <- cl.engine.Run(":" + cl.frontPort)
	}()
	select {
	case <-ctx.Done():
		return nil
	case e := <-errCh:
		return e
	}
}

func configureRouter(c *Client) {
	c.engine.Static("/css", "./internal/templates/css")
	c.engine.Static("/js", "./internal/templates/js")
	c.engine.LoadHTMLGlob("./internal/templates/*.gohtml")

	c.engine.GET("/", c.indexGet)
	c.engine.POST("/", c.indexPost)
	c.engine.GET("/excel_template", c.downloadExcel)
}
