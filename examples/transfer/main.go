package main

import (
	"fmt"

	apim "github.com/CybexDex/cybex-go/api"
)

var api apim.CybexAPI

func init() {
	// 选择需要的接入点
	api = apim.New("", "")
	// api = apim.New("wss://shenzhen.51nebula.com/", "")
	api.SetCredentials("username", "password")
}

func main() {
	re, err := api.Send("", "sendtoUsername", "0.012", "CYB", "something fun", "")
	if err != nil {
		fmt.Println(1, err)
	}
	fmt.Println(re)
}
