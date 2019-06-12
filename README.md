# cybexgolib

fork from github.com/denkhaus/bitshares

## limitOrder example

```
package main

import (
	"fmt"

	apim "github.com/CybexDex/cybex-go/api"
	"github.com/CybexDex/cybex-go/types"
)

var api apim.CybexAPI

func init() {
	// empty use default
	api = apim.New("", "")
	// api = apim.New("wss://shenzhen.51nebula.com/", "")
	// set username and password, you can ignore this and pass them with real call
	api.SetCredentials("username", "password")
}
func creatOrder() {
	// market CYB/TEST.ETH sell by price 0.012 amount 100
	re, err := api.LimitOrder("", "TEST.ETH", "CYB", "sell", "0.012", "100", "")
	if err != nil {
		fmt.Println(1, err)
	}
	fmt.Println(re)
}
func getOrder() types.LimitOrders {
	re2, err := api.LimitOrderGet("")
	if err != nil {
		fmt.Println(1, err)
	}
	return re2
}
func cancalOrder(id string) {
	re2, err := api.LimitOrderCancel("", id, "")
	if err != nil {
		fmt.Println(1, err)
	}
	fmt.Println(re2)
}
func main() {
	creatOrder()
	s := getOrder()
	for _, order := range s {
		fmt.Println(order.ID)
		cancalOrder(order.ID.String())
	}
}
```

## transfer example

```
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

```