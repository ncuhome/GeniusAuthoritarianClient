package ga

import (
	"encoding/json"
	"fmt"
	"github.com/Mmx233/tool"
	"github.com/ncuhome/GeniusAuthoritarianClient/pkg/signature"
	"net/http"
	"net/url"
	"time"
)

func NewClient(domain, appCode, appSecret string, httpClient *http.Client) *Client {
	return &Client{
		Http:      tool.NewHttpTool(httpClient),
		Domain:    domain,
		appCode:   appCode,
		appSecret: appSecret,
	}
}

type Client struct {
	Http   *tool.Http
	Domain string

	appCode   string
	appSecret string
}

func (c Client) Request(Type string, opt *tool.DoHttpReq) (*http.Response, error) {
	opt.Url = fmt.Sprintf("https://%s/api/%s", c.Domain, opt.Url)
	return c.Http.Request(Type, opt)
}

type RequestVerifyToken struct {
	Token string `json:"token"`
}
type reqVerifyToken struct {
	RequestVerifyToken
	AppCode   string `json:"appCode"`
	TimeStamp int64  `json:"timeStamp"`
	Signature string `json:"signature"`
}

func (c Client) VerifyToken(req RequestVerifyToken) (*VerifyTokenResponse, error) {
	var body = reqVerifyToken{
		RequestVerifyToken: req,
		AppCode:            c.appCode,
		TimeStamp:          time.Now().Unix(),
	}
	body.Signature = signature.Gen(&signature.VerifyClaims{
		Token:     body.Token,
		TimeStamp: body.TimeStamp,
		AppCode:   c.appCode,
		AppSecret: c.appSecret,
	})
	res, e := c.Request("POST", &tool.DoHttpReq{
		Url:  "v1/public/login/verify",
		Body: &body,
	})
	if e != nil {
		return nil, e
	}
	defer res.Body.Close()

	var resp VerifyTokenResponse
	return &resp, json.NewDecoder(res.Body).Decode(&resp)
}

func (c Client) LoginUrl(appCode string) string {
	return fmt.Sprintf("https://%s/?appCode=%s", c.Domain, url.QueryEscape(appCode))
}
