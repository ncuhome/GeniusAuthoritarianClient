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
	Res  any
}

func (c Client) Request(Type string, opt *DoReq) error {
	res, err := c.Http.Request(Type, &tool.DoHttpReq{
		Url: fmt.Sprintf("https://%s/api/v1/%s", c.Domain, opt.Url),
		Header: map[string]interface{}{
			"Content-Type": "application/json",
		},
		Body: c.signHeader.SignMap(opt.Body),
	})
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return json.NewDecoder(res.Body).Decode(opt.Res)
}

type RequestVerifyToken struct {
	Token     string `json:"token"`
	ClientIp  string `json:"clientIp,omitempty"`
	GrantType string `json:"grantType,omitempty"`

	Payload string `json:"payload,omitempty"`
	Valid   int64  `json:"valid,omitempty"`
}

func (c Client) VerifyToken(req *RequestVerifyToken) (*Response[VerifyToken], error) {
	var resp Response[VerifyToken]
	return &resp, c.Request("POST", &DoReq{
		Url:  "public/login/verify",
		Body: req,
		Res:  &resp,
	})
}

func (c Client) LoginUrl() string {
	return fmt.Sprintf("https://%s/?appCode=%s", c.Domain, url.QueryEscape(c.signHeader.AppCode))
}

type RequestRefreshToken struct {
	Token string `json:"token"`
}

func (c Client) RefreshToken(req *RequestRefreshToken) (*Response[RefreshToken], error) {
	var resp Response[RefreshToken]
	return &resp, c.Request("POST", &DoReq{
		Url:  "public/token/refresh",
		Body: req,
		Res:  &resp,
	})
}

type RequestVerifyAccessToken struct {
	Token string `json:"token"`
}

func (c Client) VerifyAccessToken(req *RequestVerifyAccessToken) (*Response[VerifyAccessToken], error) {
	var resp Response[VerifyAccessToken]
	return &resp, c.Request("POST", &DoReq{
		Url:  "public/token/access/verify",
		Body: req,
		Res:  &resp,
	})
}

func (c Client) GetUserInfo(req *RequestVerifyToken) (*Response[UserInfo], error) {
	var resp Response[UserInfo]
	return &resp, c.Request("POST", &DoReq{
		Url:  "public/token/access/user/info",
		Body: req,
		Res:  &resp,
	})
}
