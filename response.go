package geniusAuth

import "encoding/json"

type Response struct {
	Code uint        `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

type VerifyTokenData struct {
	UserID       uint     `json:"userID"`
	Name         string   `json:"name"`
	Groups       []string `json:"groups"`
	AvatarUrl    string   `json:"avatarUrl"`
	RefreshToken string   `json:"refreshToken"`
	AccessToken  string   `json:"accessToken"`
}

type VerifyTokenResponse struct {
	Response
	Data *VerifyTokenData `json:"data"`
}

type RefreshTokenData struct {
	AccessToken string          `json:"access_token"`
	Payload     json.RawMessage `json:"payload,omitempty"`
}

type RefreshTokenResponse struct {
	Response
	Data *RefreshTokenData `json:"data"`
}

type VerifyAccessTokenData struct {
	UID     uint            `json:"uid"`
	Payload json.RawMessage `json:"payload,omitempty"`
}

type VerifyAccessTokenResponse struct {
	Response
	Data *VerifyAccessTokenData `json:"data"`
}
