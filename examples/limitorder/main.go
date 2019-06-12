package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/CybexDex/cybex-go/api"
	"github.com/CybexDex/cybex-go/config"
	"github.com/CybexDex/cybex-go/crypto"
	"github.com/CybexDex/cybex-go/operations"
	"github.com/CybexDex/cybex-go/types"
	"github.com/juju/errors"
)

var (
	cyb    = types.NewGrapheneID("1.3.0")
	cny    = types.NewGrapheneID("1.3.113")
	keyBag *crypto.KeyBag
	seller *types.GrapheneID
)

const (
	wsFullApiUrl = "wss://cybex.openledger.info/ws"
)

func init() {
	// init is called before the API is initialized,
	// hence must define current chain config explicitly.
	config.SetCurrentConfig(config.ChainIDCYB)
	seller = types.NewGrapheneID(
		os.Getenv("CYB_TEST_ACCOUNT"),
	)
	keyBag = crypto.NewKeyBag()
	if err := keyBag.Add(os.Getenv("CYB_TEST_WIF")); err != nil {
		log.Fatal(errors.Annotate(err, "Add [wif]"))
	}
}

func main() {
	api := api.New(wsFullApiUrl, "")
	if err := api.Connect(); err != nil {
		log.Fatal(errors.Annotate(err, "OnConnect"))
	}

	api.OnError(func(err error) {
		log.Fatal(errors.Annotate(err, "OnError"))
	})

	op := operations.LimitOrderCreateOperation{
		FillOrKill: false,
		Seller:     *seller,
		Extensions: types.Extensions{},
		AmountToSell: types.AssetAmount{
			Amount: 100,
			Asset:  *cyb,
		},
		MinToReceive: types.AssetAmount{
			Amount: 1000000,
			Asset:  *cny,
		},
	}

	op.Expiration.Set(24 * time.Hour)

	tx, err := api.BuildSignedTransaction(keyBag, cyb, &op)
	if err != nil {
		log.Fatal(errors.Annotate(err, "BuildSignedTransaction"))
	}

	if err := api.BroadcastTransaction(tx); err != nil {
		log.Fatal(errors.Annotate(err, "BroadcastTransaction"))
	}

	fmt.Println("operation successfull broadcasted")
}
