package api

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/CybexDex/cybex-go/types"
	"github.com/CybexDex/cybex-go/util"
	"github.com/btcsuite/btcd/btcec"
)

// SignStr ...
func (p *cybexAPI) SignStr(toSign string, prikey string) (sign string, err error) {
	pri, err := types.NewPrivateKeyFromWif(prikey)
	if err != nil {
		return "", fmt.Errorf("SignStr NewPrivateKeyFromWif: %v", err)
	}
	writer := sha256.New()
	writer.Write([]byte(toSign))
	// fmt.Println([]byte(toSign))
	digest := writer.Sum(nil)
	dig := digest[:]
	// fmt.Println("dig", dig, pri.PublicKey().String())
	sig, err := pri.SignCompact(dig)
	// fmt.Println("sig", sig)
	sign = hex.EncodeToString(sig)
	// ps, _, err := btcec.RecoverCompact(btcec.S256(), sig, dig)
	// if err != nil {
	// 	return "false", fmt.Errorf("gen raw pubkey error: %v", err)
	// }
	// pub, err := types.NewPublicKey(ps)
	// if err != nil {
	// 	return "false", fmt.Errorf("gen pubkey error: %v", err)
	// }
	// fmt.Println(pub)
	return sign, nil
}

// VerifySign ...
func (p *cybexAPI) VerifySign(accountName string, toSign string, sign string) (re bool, err error) {
	account, err := p.GetAccountByName(accountName)
	if err != nil {
		return false, fmt.Errorf("get account error: %v", err)
	}
	writer := sha256.New()
	writer.Write([]byte(toSign))
	// fmt.Println([]byte(toSign))
	digest := writer.Sum(nil)
	dig := digest[:]
	sig, _ := hex.DecodeString(sign)
	ps, _, err := btcec.RecoverCompact(btcec.S256(), sig, dig)
	if err != nil {
		return false, fmt.Errorf("gen raw pubkey error: %v", err)
	}
	pub, err := types.NewPublicKey(ps)
	if err != nil {
		return false, fmt.Errorf("gen pubkey error: %v", err)
	}
	// fmt.Println(pub)
	for k, _ := range account.Active.KeyAuths {
		if k.Equal(pub) {
			return true, nil
		}
	}

	return false, nil
}
func (p *cybexAPI) LoginVerify(fund types.Fund, sign string) (re bool, err error) {
	account, err := p.GetAccountByName(string(fund.AccountName))
	if err != nil {
		return false, fmt.Errorf("get account error: %v", err)
	}
	var b bytes.Buffer
	//fmt.Println("name",fund.AccountName.Bytes())
	enc := util.NewTypeEncoder(&b)
	if err := enc.Encode(fund); err != nil {
		//fmt.Println("fucked")
		return false, fmt.Errorf("fund error: %v", err)
	}
	writer := sha256.New()
	writer.Write(b.Bytes())
	digest := writer.Sum(nil)
	dig := digest[:]
	sig, _ := hex.DecodeString(sign)
	ps, _, err := btcec.RecoverCompact(btcec.S256(), sig, dig)
	if err != nil {
		return false, fmt.Errorf("gen raw pubkey error: %v", err)
	}

	pub, err := types.NewPublicKey(ps)
	if err != nil {
		return false, fmt.Errorf("gen pubkey error: %v", err)
	}
	for k, _ := range account.Active.KeyAuths {
		if k.Equal(pub) {
			return true, nil
		}
	}

	return false, nil
}
