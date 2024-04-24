package api

import (
	"net/http"

	"github.com/astaxie/beego/httplib"
)

type Client struct {
	Request *Request
}

func (c *Client) SendRequest() (*http.Response, error) {
	beegoRequest := httplib.NewBeegoRequest(
		c.Request.BaseURL+c.Request.Path,
		c.Request.Method.String())
	for key, value := range c.Request.HeaderField {
		beegoRequest.Header(key, value)
	}
	for key, value := range c.Request.Params {
		beegoRequest.Param(key, value)
	}
	res, err := beegoRequest.Response()
	if err != nil {
		return nil, err
	}
	return res, nil
}
