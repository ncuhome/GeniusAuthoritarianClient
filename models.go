package ga

type Response struct {
	Code uint        `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

type VerifyTokenSuccess struct {
	Name   string   `json:"name"`
	Groups []string `json:"groups"`
}

type VerifyTokenResponse struct {
	Response
	Data *VerifyTokenSuccess `json:"data"`
}
