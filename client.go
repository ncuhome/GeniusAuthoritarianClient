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

func Request[T any](c Client, Type string, opt *DoReq) (*T, error) {
	var resp Response[T]
	res, err := c.Http.Request(Type, &tool.DoHttpReq{
		Url: fmt.Sprintf("https://%s/api/v1/%s", c.Domain, opt.Url),
		Header: map[string]interface{}{
			"Content-Type": "application/json",
		},
		Body: c.signHeader.SignMap(opt.Body),
	})
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&resp)
	if err != nil {
		return nil, err
	}

	if resp.Code != 0 {
		return nil, &ApiErr{
			Code: resp.Code,
			Msg:  resp.Msg,
		}
	}
	return &resp.Data, nil
}

type RequestVerifyToken struct {
	Token     string `json:"token"`
	ClientIp  string `json:"clientIp,omitempty"`
	GrantType string `json:"grantType,omitempty"`

	Payload string `json:"payload,omitempty"`
	Valid   int64  `json:"valid,omitempty"`
}

func (c Client) VerifyToken(req *RequestVerifyToken) (*VerifyToken, error) {
	return Request[VerifyToken](c, "POST", &DoReq{
		Url:  "public/login/verify",
		Body: req,
	})
}

func (c Client) LoginUrl() string {
	return fmt.Sprintf("https://%s/?appCode=%s", c.Domain, url.QueryEscape(c.signHeader.AppCode))
}

type RequestRefreshToken struct {
	Token string `json:"token"`
}

func (c Client) RefreshToken(req *RequestRefreshToken) (*RefreshToken, error) {
	return Request[RefreshToken](c, "POST", &DoReq{
		Url:  "public/token/refresh",
		Body: req,
	})
}

type RequestModifyPayload struct {
	Token       string `json:"token"`
	Payload     string `json:"payload"`
	AccessToken bool   `json:"accessToken"`
}

func (c Client) ModifyPayload(req *RequestModifyPayload) (*Tokens, error) {
	return Request[Tokens](c, "PATCH", &DoReq{
		Url:  "public/token/refresh",
		Body: req,
	})
}

type RequestVerifyAccessToken struct {
	Token string `json:"token"`
}

func (c Client) VerifyAccessToken(req *RequestVerifyAccessToken) (*VerifyAccessToken, error) {
	return Request[VerifyAccessToken](c, "POST", &DoReq{
		Url:  "public/token/access/verify",
		Body: req,
	})
}

func (c Client) GetUserInfo(req *RequestVerifyToken) (*UserInfo, error) {
	return Request[UserInfo](c, "POST", &DoReq{
		Url:  "public/token/access/user/info",
		Body: req,
	})
}
