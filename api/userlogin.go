package api

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"

	"github.com/CybexDex/cybex-go/crypto"
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
)

func Sha256f(input []byte, outputEncoding string) string {
	hash := sha256.New()
	_, err := hash.Write(input)
	if err != nil {
		fmt.Println(err)
	}
	sum := hash.Sum(nil)

	switch outputEncoding {
	case "base64":
		return base64.StdEncoding.EncodeToString(sum)

	case "base64url":
		return base64.URLEncoding.EncodeToString(sum)

	case "base64rawurl":
		return base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(sum)

	case "hex":
		return hex.EncodeToString(sum)

	default:
		return ""
	}

	return ""
}
func makeKey(username string, password string, flag string) string {
	tempSeed := fmt.Sprintf("%s%s%s", username, flag, password)
	sha256ed := Sha256f([]byte(tempSeed), "hex")

	//pri,_ := types.NewPrivateKeyFromHex(sha256ed)
	seedBytes, _ := hex.DecodeString(sha256ed)
	privKey, _ := btcec.PrivKeyFromBytes(btcec.S256(), seedBytes)
	tempP, _ := btcutil.NewWIF(privKey, &chaincfg.MainNetParams, false)
	//wif := tempP.ToWIF()
	wif := tempP.String()
	// fmt.Println(tempSeed,sha256ed,wif)
	return wif
}
func UserPasstoPri(username string, password string) []string {
	active := makeKey(username, password, "active")
	owner := makeKey(username, password, "owner")
	memo := makeKey(username, password, "memo")
	return []string{active, owner, memo}
}
func KeyBagByUserPass(username string, password string) *crypto.KeyBag {
	keys := UserPasstoPri(username, password)
	keyBag := crypto.NewKeyBag()
	for _, key := range keys {
		keyBag.Add(key)
	}
	return keyBag
}
