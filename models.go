package ga

type Response struct {
	Code uint        `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

type VerifyTokenSuccess struct {
	UserID    uint     `json:"userID"`
	Name      string   `json:"name"`
	Groups    []string `json:"groups"`
	AvatarUrl string   `json:"avatarUrl"`
}

type VerifyTokenResponse struct {
	Response
	Data *VerifyTokenSuccess `json:"data"`
}
