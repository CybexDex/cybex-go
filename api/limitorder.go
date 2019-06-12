package api

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/CybexDex/cybex-go/crypto"
	"github.com/CybexDex/cybex-go/operations"
	"github.com/CybexDex/cybex-go/types"
	"github.com/juju/errors"
	"github.com/shopspring/decimal"
)

func (p *cybexAPI) LimitOrderGet(user string) (types.LimitOrders, error) {
	if user == "" {
		user = p.username
	}
	fromUser, err := p.GetAccountByName(user)
	fulls, err := p.GetFullAccounts(fromUser.ID)
	return fulls[0].AccountInfo.LimitOrders, err
}
func (p *cybexAPI) LimitOrderCancel(user string, orderid string, password string) (*types.SignedTransaction, error) {
	if user == "" {
		user = p.username
	}
	if password == "" {
		password = p.password
	}
	fromUser, err := p.GetAccountByName(user)
	keyBag := KeyBagByUserPass(user, password)
	if err != nil {
		return nil, err
	}
	orderID := types.NewGrapheneID(orderid)
	op := operations.LimitOrderCancelOperation{
		FeePayingAccount: fromUser.ID,
		Extensions:       types.Extensions{},
		Order:            *orderID,
	}

	// op.Expiration.Set(24 * time.Hour)
	tx, err := p.BuildSignedTransaction(keyBag, cybAsset, &op)
	if err != nil {
		log.Fatal(errors.Annotate(err, "BuildSignedTransaction"))
	}
	crypto.NewTransactionSigner(tx)
	if err := p.BroadcastTransaction(tx); err != nil {
		//log.Fatal(errors.Annotate(err, "BroadcastTransaction"))
		return nil, err
	}
	return tx, nil
}
func (p *cybexAPI) LimitOrder(user string, base string, quote string, action string, price string, amount string, password string) (*types.SignedTransaction, error) {
	if user == "" {
		user = p.username
	}
	if password == "" {
		password = p.password
	}
	if action != "sell" && action != "buy" {
		return nil, fmt.Errorf("action only support sell/buy")
	}
	fromUser, err := p.GetAccountByName(user)
	keyBag := KeyBagByUserPass(user, password)
	if err != nil {
		return nil, err
	}
	baseObj, err := p.GetAsset(base)
	if err != nil {
		return nil, err
	}
	quoteObj, err := p.GetAsset(quote)
	if err != nil {
		return nil, err
	}
	amountD, err := decimal.NewFromString(amount)
	tenD, err := decimal.NewFromString("10")
	pstr := strconv.Itoa(quoteObj.Precision)
	precisionD, err := decimal.NewFromString(pstr)
	quoteAmountRaw := tenD.Pow(precisionD).Mul(amountD).IntPart()
	amountInt := types.Int64(quoteAmountRaw)
	priceD, err := decimal.NewFromString(price)
	amountB := amountD.Mul(priceD)

	bstr := strconv.Itoa(baseObj.Precision)
	precisionB, err := decimal.NewFromString(bstr)
	baseAmountRaw := tenD.Pow(precisionB).Mul(amountB).IntPart()
	bamountInt := types.Int64(baseAmountRaw)
	var sell types.AssetAmount
	var buy types.AssetAmount
	if action == "sell" {
		sell = types.AssetAmount{
			Amount: amountInt,
			Asset:  quoteObj.ID,
		}
		buy = types.AssetAmount{
			Amount: bamountInt,
			Asset:  baseObj.ID,
		}
	} else {
		sell = types.AssetAmount{
			Amount: bamountInt,
			Asset:  baseObj.ID,
		}
		buy = types.AssetAmount{
			Amount: amountInt,
			Asset:  quoteObj.ID,
		}
	}

	op := operations.LimitOrderCreateOperation{
		FillOrKill:   false,
		Seller:       fromUser.ID,
		Extensions:   types.Extensions{},
		AmountToSell: sell,
		MinToReceive: buy,
	}

	op.Expiration.Set(24 * time.Hour)

	tx, err := p.BuildSignedTransaction(keyBag, cybAsset, &op)
	if err != nil {
		log.Fatal(errors.Annotate(err, "BuildSignedTransaction"))
	}
	crypto.NewTransactionSigner(tx)
	if err := p.BroadcastTransaction(tx); err != nil {
		//log.Fatal(errors.Annotate(err, "BroadcastTransaction"))
		return nil, err
	}
	return tx, nil
}
