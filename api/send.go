package api

import (
	"log"
	"math/rand"
	"strings"

	"strconv"

	"github.com/CybexDex/cybex-go/crypto"
	"github.com/CybexDex/cybex-go/operations"
	"github.com/CybexDex/cybex-go/types"
	"github.com/juju/errors"
	"github.com/shopspring/decimal"
)

// var userkeybags map[[]string]*crypto.KeyBag
var (
	cybAsset = types.NewGrapheneID("1.3.0")
)

func (p *cybexAPI) makeOp(from string, to string, amount string, asset string, memo string, password string) (*operations.TransferOperation, error) {
	nonce := types.UInt64(rand.Int63())
	keyBag := crypto.NewKeyBag()
	passArr := strings.Split(password, ",")
	lenpass := len(passArr)
	if lenpass == 1 {
		keyBag1 := KeyBagByUserPass(from, password)
		keyBag.Merge(keyBag1)
	} else if lenpass >= 2 {
		for _, newkey := range passArr {
			keyBag.Add(newkey)
		}
	}
	fromUser, err := p.GetAccountByName(from)
	if err != nil {
		return nil, err
	}
	toUser, err := p.GetAccountByName(to)
	if err != nil {
		return nil, err
	}
	assetObj, err := p.GetAsset(asset)
	if err != nil {
		return nil, err
	}
	amountD, err := decimal.NewFromString(amount)
	tenD, err := decimal.NewFromString("10")
	pstr := strconv.Itoa(assetObj.Precision)
	precisionD, err := decimal.NewFromString(pstr)
	amountRaw := tenD.Pow(precisionD).Mul(amountD).IntPart()
	amountInt := types.Int64(amountRaw)
	op := operations.TransferOperation{
		Extensions: types.Extensions{},
		// Memo:       &memo1,
		From: fromUser.ID,
		To:   toUser.ID,
		Amount: types.AssetAmount{
			Amount: amountInt,
			Asset:  assetObj.ID,
		},
	}
	if memo != "" {
		memo1 := types.Memo{
			From:  fromUser.Options.MemoKey,
			To:    toUser.Options.MemoKey,
			Nonce: nonce,
		}
		pubkeys := types.PublicKeys{fromUser.Options.MemoKey}
		priKeyA := keyBag.PrivatesByPublics(pubkeys)
		for _, prv := range priKeyA {
			memo1.Encrypt(&prv, memo)
		}
		op.Memo = &memo1
	}
	return &op, nil
}
func (p *cybexAPI) Sends(tosends []types.SimpleSend) (tx *types.SignedTransaction, err error) {
	tx, err = p.PreSends(tosends)
	if err := p.BroadcastTransaction(tx); err != nil {
		//log.Fatal(errors.Annotate(err, "BroadcastTransaction"))
		return nil, err
	}
	return tx, nil
}

func (p *cybexAPI) PreSends(tosends []types.SimpleSend) (tx *types.SignedTransaction, err error) {
	keyBag := crypto.NewKeyBag()
	ops := make([]*operations.TransferOperation, 0)
	for _, tosend := range tosends {
		op, err := p.makeOp(tosend.From, tosend.To, tosend.Amount, tosend.Asset, tosend.Memo, tosend.Password)
		if err != nil {
			return nil, err
		}
		ops = append(ops, op)
		passArr := strings.Split(tosend.Password, ",")
		lenpass := len(passArr)
		if lenpass == 1 {
			keyBag1 := KeyBagByUserPass(tosend.From, tosend.Password)
			keyBag.Merge(keyBag1)
		} else if lenpass >= 2 {
			for _, newkey := range passArr {
				keyBag.Add(newkey)
			}
		}
	}
	switch len(ops) {
	case 1:
		tx, err = p.BuildSignedTransaction(keyBag, cybAsset, ops[0])
		if err != nil {
			log.Fatal(errors.Annotate(err, "BuildSignedTransaction"))
		}
	case 2:
		tx, err = p.BuildSignedTransaction(keyBag, cybAsset, ops[0], ops[1])
		if err != nil {
			log.Fatal(errors.Annotate(err, "BuildSignedTransaction"))
		}
	case 3:
		tx, err = p.BuildSignedTransaction(keyBag, cybAsset, ops[0], ops[1], ops[2])
		if err != nil {
			log.Fatal(errors.Annotate(err, "BuildSignedTransaction"))
		}
	}
	crypto.NewTransactionSigner(tx)
	return tx, err
}
func (p *cybexAPI) Send(from string, to string, amount string, asset string, memo string, password string) (*types.SignedTransaction, error) {
	if from == "" {
		from = p.username
	}
	if password == "" {
		password = p.password
	}
	nonce := types.UInt64(rand.Int63())
	fromUser, err := p.GetAccountByName(from)
	keyBag := KeyBagByUserPass(from, password)
	if err != nil {
		return nil, err
	}
	toUser, err := p.GetAccountByName(to)
	if err != nil {
		return nil, err
	}
	assetObj, err := p.GetAsset(asset)
	if err != nil {
		return nil, err
	}
	amountD, err := decimal.NewFromString(amount)
	tenD, err := decimal.NewFromString("10")
	pstr := strconv.Itoa(assetObj.Precision)
	precisionD, err := decimal.NewFromString(pstr)
	amountRaw := tenD.Pow(precisionD).Mul(amountD).IntPart()
	amountInt := types.Int64(amountRaw)
	op := operations.TransferOperation{
		Extensions: types.Extensions{},
		// Memo:       &memo1,
		From: fromUser.ID,
		To:   toUser.ID,
		Amount: types.AssetAmount{
			Amount: amountInt,
			Asset:  assetObj.ID,
		},
	}
	if memo != "" {
		memo1 := types.Memo{
			From:  fromUser.Options.MemoKey,
			To:    toUser.Options.MemoKey,
			Nonce: nonce,
		}
		pubkeys := types.PublicKeys{fromUser.Options.MemoKey}
		priKeyA := keyBag.PrivatesByPublics(pubkeys)
		for _, prv := range priKeyA {
			memo1.Encrypt(&prv, memo)
		}
		op.Memo = &memo1
	}
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
