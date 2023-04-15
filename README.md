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

var GaClient = ga.NewClient("v.ncuos.com", http.DefaultClient)

func Login(c *gin.Context) {
	var f struct {
		Token string `json:"token" form:"token" binding:"required"`
	}
	if e := c.ShouldBind(&f); e != nil {
		panic(e)
		return
	}

	info, e := GaClient.VerifyToken(f.Token)
	if e != nil {
		panic(e)
		return
	} else if info.Code != 0 {
		panic(info.Msg)
		return
	}

	// 登录成功
	fmt.Println(info.Data)
}
```