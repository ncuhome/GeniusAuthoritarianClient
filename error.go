package geniusAuth

import "fmt"

type ApiErr struct {
	Code uint
	Msg  string
}

func (e ApiErr) Error() string {
	return fmt.Sprintf("genius-auth api server return code: %d, error msg: %s", e.Code, e.Msg)
}
