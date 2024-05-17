package jwtClaims

type RefreshToken struct {
	UserClaims
	ID      uint64 `json:"id"`
	AppCode string `json:"appCode"`
	Payload string `json:"payload,omitempty"`
}

func (t RefreshToken) GetID() uint64 {
	return t.ID
}
func (t RefreshToken) GetAppCode() string {
	return t.AppCode
}

type AccessToken struct {
	RefreshToken
}
