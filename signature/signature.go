package signature

import (
	"crypto/sha256"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"time"
)

type SignHeader struct {
	AppCode   string
	AppSecret string
}

func (s SignHeader) SignMap(body any) map[string]interface{} {
	v := reflect.ValueOf(body).Elem()
	t := v.Type()
	num := t.NumField()
	bodyMap := make(map[string]interface{}, num+3)
	for i := 0; i < num; i++ {
		tagStr := t.Field(i).Tag.Get("json")
		tags := strings.Split(tagStr, ",")
		if len(tags) >= 2 && tags[1] == "omitempty" && v.Field(i).IsZero() {
			continue
		}
		bodyMap[tags[0]] = v.Field(i).Interface()
	}

	unix := time.Now().Unix()
	bodyMap["appCode"] = s.AppCode
	bodyMap["timeStamp"] = unix

	signMap := make(map[string]string, len(bodyMap)+1)
	signMap["appSecret"] = s.AppSecret
	for key, value := range bodyMap {
		signMap[key] = fmt.Sprint(value)
	}

	var keySlice = make([]string, len(signMap))
	var signStrLen = len(signMap)*2 - 1
	i := 0
	for key, value := range signMap {
		keySlice[i] = key
		signStrLen += len(key) + len(value)
		i++
	}
	sort.Strings(keySlice)

	var signStr strings.Builder
	signStr.Grow(signStrLen)
	for i, key := range keySlice {
		if i != 0 {
			signStr.Write([]byte("&"))
		}
		signStr.Write([]byte(key + "=" + signMap[key]))
	}

	h := sha256.New()
	h.Write([]byte(signStr.String()))
	bodyMap["signature"] = fmt.Sprintf("%x", h.Sum(nil))
	return bodyMap
}
