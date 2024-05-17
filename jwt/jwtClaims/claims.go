package jwtClaims

type RefreshToken struct {
	UserClaims
	ID      uint64 `json:"id"`
	AppCode string `json:"appCode"`
	Payload string `json:"payload,omitempty"`
}

type AccessToken struct {
	RefreshToken
}
