package geniusAuth

import (
	"encoding/json"
	"fmt"
	"github.com/Mmx233/tool"
	"github.com/ncuhome/GeniusAuthoritarianClient/signature"
	"net/http"
	"net/url"
)

func NewClient(domain, appCode, appSecret string, httpClient *http.Client) *Client {
	return &Client{
		Http:   tool.NewHttpTool(httpClient),
		Domain: domain,

		signHeader: &signature.SignHeader{
			AppCode:   appCode,
			AppSecret: appSecret,
		},
	}
}

type Client struct {
	Http   *tool.Http
	Domain string

	signHeader *signature.SignHeader
}

type DoReq struct {
	Url  string
	Body interface{}
}

func (c Client) Request(Type string, opt *DoReq) (*http.Response, error) {
	return c.Http.Request(Type, &tool.DoHttpReq{
		Url: fmt.Sprintf("https://%s/api/v1/%s", c.Domain, opt.Url),
		Header: map[string]interface{}{
			"Content-Type": "application/json",
		},
		Body: c.signHeader.SignMap(opt.Body),
	})
}

type RequestVerifyToken struct {
	Token     string `json:"token"`
	ClientIp  string `json:"clientIp,omitempty"`
	GrantType string `json:"grantType,omitempty"`

	Payload string `json:"payload,omitempty"`
	Valid   int64  `json:"valid,omitempty"`
}

func (c Client) VerifyToken(req *RequestVerifyToken) (*VerifyTokenResponse, error) {
	res, err := c.Request("POST", &DoReq{
		Url:  "public/login/verify",
		Body: req,
	})
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var resp VerifyTokenResponse
	return &resp, json.NewDecoder(res.Body).Decode(&resp)
}

func (c Client) LoginUrl() string {
	return fmt.Sprintf("https://%s/?appCode=%s", c.Domain, url.QueryEscape(c.signHeader.AppCode))
}

type RequestRefreshToken struct {
	Token string `json:"token"`
}

func (c Client) RefreshToken(req *RequestRefreshToken) (*RefreshTokenResponse, error) {
	res, err := c.Request("POST", &DoReq{
		Url:  "public/token/refresh",
		Body: req,
	})
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var resp RefreshTokenResponse
	return &resp, json.NewDecoder(res.Body).Decode(&resp)
}

type RequestVerifyAccessToken struct {
	Token string `json:"token"`
}

func (c Client) VerifyAccessToken(req *RequestVerifyAccessToken) (*VerifyAccessTokenResponse, error) {
	res, err := c.Request("POST", &DoReq{
		Url:  "public/token/access/verify",
		Body: req,
	})
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var resp VerifyAccessTokenResponse
	return &resp, json.NewDecoder(res.Body).Decode(&resp)
}
