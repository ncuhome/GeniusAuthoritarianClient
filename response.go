package geniusAuth

type Response[T any] struct {
	Code uint   `json:"code"`
	Data T      `json:"data"`
	Msg  string `json:"msg"`
}

type VerifyToken struct {
	UserID       uint     `json:"userID"`
	Name         string   `json:"name"`
	Groups       []string `json:"groups"`
	AvatarUrl    string   `json:"avatarUrl"`
	RefreshToken string   `json:"refreshToken"`
	AccessToken  string   `json:"accessToken"`
}

type RefreshToken struct {
	AccessToken string `json:"access_token"`
	Payload     string `json:"payload,omitempty"`
}

type VerifyAccessToken struct {
	UID     uint   `json:"uid"`
	Payload string `json:"payload,omitempty"`
}

type UserInfo struct {
	UserID    uint     `json:"userID"`
	Name      string   `json:"name"`
	Groups    []string `json:"groups"`
	AvatarUrl string   `json:"avatarUrl"`
}

type Tokens struct {
	RefreshToken string `json:"refreshToken"`
	AccessToken  string `json:"accessToken,omitempty"`
}
