package ga

import (
	"encoding/json"
	"fmt"
	"github.com/Mmx233/tool"
	"net/http"
)

func NewClient(domain string, httpClient *http.Client) *Client {
	return &Client{
		Http:   tool.NewHttpTool(httpClient),
		Domain: domain,
	}
}

type Client struct {
	Http   *tool.Http
	Domain string
}

func (c Client) Request(Type string, opt *tool.DoHttpReq) (*http.Response, error) {
	opt.Url = fmt.Sprintf("https://%s/api/%s", c.Domain, opt.Url)
	return c.Http.Request(Type, opt)
}

type RequestVerifyToken struct {
	Token  string   `json:"token"`
	Groups []string `json:"groups,omitempty"`
}

func (c Client) VerifyToken(req *RequestVerifyToken) (*VerifyTokenResponse, error) {
	res, e := c.Request("POST", &tool.DoHttpReq{
		Url:  "v1/public/login/verify",
		Body: req,
	})
	if e != nil {
		return nil, e
	}
	defer res.Body.Close()

	var resp VerifyTokenResponse
	return &resp, json.NewDecoder(res.Body).Decode(&resp)
}
