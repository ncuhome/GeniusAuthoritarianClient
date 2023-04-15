# GeniusAuthoritarianClient

使用：

```shell
~$ export GOPRIVATE=github.com/ncuhome
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

	info, e := GaClient.VerifyToken(&ga.RequestVerifyToken{
		Token: f.Token,
    })
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

在接口中限制可登录组：

```go
GaClient.VerifyToken(&ga.RequestVerifyToken{
    Token: f.Token,
    Groups: {"中心", "研发"}
})
```