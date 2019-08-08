package automation

import (
	"encoding/json"

	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

var log = logrus.StandardLogger()

type Client struct {
	httpCli *fasthttp.Client
	url     string
}

func NewClient(url string, httpCli *fasthttp.Client) *Client {
	return &Client{
		url:     url,
		httpCli: httpCli,
	}
}

func (c *Client) Trigger(data interface{}) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	req.Header.SetMethod("POST")
	req.Header.SetContentType("application/json")
	req.SetRequestURI(c.url + "/trigger")
	body, err := json.Marshal(data)
	if err != nil {
		log.WithError(err).Error("Failed to trigger: Failed to decode JSON")
		return
	}
	req.SwapBody(body)
	if err := c.httpCli.Do(req, resp); err != nil {
		log.WithError(err).Error("Failed to trigger: Failed to send request")
	}
}
