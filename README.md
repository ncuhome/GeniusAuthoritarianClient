# GeniusAuthoritarianClient

使用：

```shell
~$ go get github.com/ncuhome/GeniusAuthoritarianClient
```

示例：

```go
package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	ga "github.com/ncuhome/GeniusAuthoritarianClient"
	"net/http"
)

var GaClient = ga.NewClient("v.ncuos.com" ,"your app code", "your app secret", http.DefaultClient)

// Login 一次性校验模式验证登录身份
func Login(c *gin.Context) {
	var f struct {
		// 一次性令牌
		Token string `json:"token" form:"token" binding:"required"`
	}
	if err := c.ShouldBind(&f); err != nil {
		panic(err)
		return
	}

	info, err := GaClient.VerifyToken(&ga.RequestVerifyToken{
		Token: f.Token,
		ClientIp: c.ClientIp(),
	})
	if err != nil {
		panic(err)
		return
	} else if info.Code != 0 {
		panic(info.Msg)
		return
	}

	// 登录成功
	fmt.Println(info.Data)
}

// GoLogin 跳转到 GeniusAuth 登录
func GoLogin(c *gin.Context)  {
    c.Redirect(302, GaClient.LoginUrl())
}
```
