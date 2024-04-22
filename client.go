package geniusAuth

import (
	"encoding/json"
	"fmt"
	"github.com/Mmx233/tool"
	"github.com/ncuhome/GeniusAuthoritarianClient/signature"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func NewClient(domain, appCode, appSecret string, httpClient *http.Client) *Client {
	return &Client{
		Http:    tool.NewHttpTool(httpClient),
		Domain:  domain,
		AppCode: appCode,

		signHeader: &signature.SignHeader{
			AppCode:   appCode,
			AppSecret: appSecret,
		},
	}
}

type Client struct {
	Http *tool.Http

	Domain  string
	AppCode string

	signHeader *signature.SignHeader
}

type DoReq struct {
	Url  string
	Form interface{}
}

func Request[T any](c Client, Type string, opt *DoReq) (*T, error) {
	req := &tool.DoHttpReq{
		Url: fmt.Sprintf("https://%s/api/v1/%s", c.Domain, opt.Url),
		Header: map[string]interface{}{
			"Content-Type": "application/json",
		},
	}
	signedMap := c.signHeader.SignMap(opt.Form)
	if Type == "GET" {
		req.Query = signedMap
	} else {
		req.Body = signedMap
	}

	var resp Response[T]
	res, err := c.Http.Request(Type, req)
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
		Form: req,
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
		Url:  "app/token/refresh",
		Form: req,
	})
}

type RequestModifyPayload struct {
	Token       string `json:"token"`
	Payload     string `json:"payload"`
	AccessToken bool   `json:"accessToken"`
}

func (c Client) ModifyPayload(req *RequestModifyPayload) (*Tokens, error) {
	return Request[Tokens](c, "PATCH", &DoReq{
		Url:  "app/token/refresh",
		Form: req,
	})
}

type RequestVerifyAccessToken struct {
	Token string `json:"token"`
}

func (c Client) VerifyAccessToken(req *RequestVerifyAccessToken) (*VerifyAccessToken, error) {
	return Request[VerifyAccessToken](c, "POST", &DoReq{
		Url:  "app/token/access/verify",
		Form: req,
	})
}

func (c Client) GetUserInfo(req *RequestVerifyToken) (*UserInfo, error) {
	return Request[UserInfo](c, "POST", &DoReq{
		Url:  "app/token/access/user/info",
		Form: req,
	})
}

type RequestGetUserPublicInfo struct {
	ID string `json:"id"`
}

func (c Client) GetUserPublicInfo(uid ...uint) ([]UserPublicInfo, error) {
	idStrArr := make([]string, len(uid))
	for i, id := range uid {
		idStrArr[i] = strconv.FormatUint(uint64(id), 10)
	}
	resp, err := Request[[]UserPublicInfo](c, "GET", &DoReq{
		Url: "app/user/info",
		Form: &RequestGetUserPublicInfo{
			ID: strings.Join(idStrArr, ","),
		},
	})
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

func (c Client) GetServerPublicKeys() (*ServerPublicKeys, error) {
	return Request[ServerPublicKeys](c, "GET", &DoReq{
		Url: "app/keypair/server",
	})
}

func (c Client) CreateRpcClientCredential() (*RpcClientCredential, error) {
	return Request[RpcClientCredential](c, "POST", &DoReq{
		Url: "app/keypair/rpc",
	})
}
