package types

import (
	"bytes"
	"testing"

	"github.com/CybexDex/cybex-go/config"
	"github.com/CybexDex/cybex-go/util"
	"github.com/stretchr/testify/assert"
)

var privKeys = [][]string{
	{"5Hx8KiHLnc3pDLkwe2jujkTTJev72n3Qx7xtyaRNBsJDuejzh9u", "CYB5zzvbDtkbUVU1gFFsKqCE55U7JbjTp6mTh1usFv7KGgXL7HDQk"},
	{"5JyuWmopuyxFyvM9xm8fxXyujzfVnsg9cvE6z3pcib5NW1Av4rP", "CYB5yXqEBShUgcVm7Mve8Fg4RzQ2ftPpmo77aMbz884eX9aeGVvwD"},
	{"5KRZv3ZmkcE71K9KwEKG6pV6pyufkMQgCJrCu8xKLf2y7R7J8YK", "CYB5Z3vsgH6xj6HLXcsU38yo4TyoZs9AUzpfbaXbuxsAYPbutWvEP"},
	{"5JTge2oTwFqfNPhUrrm6upheByG2VXvaXBAqWdDUvK2CsygMG3Z", "CYB5shffTjVoT4J8Zrj3f2mQJw4UVKrnbx5FWYhVgov45EpBf2NYi"},
	{"5JqmjeakPoTz3ComQ7Jgg11jHxywfkJHZPhMJoBomZLrZSfRAvr", "CYB56fy8qpkLzNoguGMPgPNkkznxnx2woEg1qPq7E6gF2SeGSRyK5"},
}

var data = [][]string{
	{"5JWHY5DxTF6qN5grTtChDCYBmWHfY9zaSsw4CxEKN5eZpH9iBma", "5ad2b8df2c255d4a2996ee7d065e013e1bbb35c075ee6e5208aca44adc9a9d4c"},
	{"5KPipdRzoxrp6dDqsBfMD6oFZG356trVHV5QBGx3rABs1zzWWs8", "cf9d6121ed458f24ea456ad7ff700da39e86688988cfe5c6ed6558642cf1e32f"},
}

func Test_PrivatePublic(t *testing.T) {
	config.SetCurrentConfig(config.ChainIDCYB)

	for _, k := range privKeys {
		wif := k[0]
		pub := k[1]

		key, err := NewPrivateKeyFromWif(wif)
		if err != nil {
			assert.FailNow(t, err.Error(), "NewPrivateKeyFromWif")
		}

		assert.Equal(t, pub, key.PublicKey().String())
		assert.Equal(t, wif, key.ToWIF())
	}
}

func TestDecode(t *testing.T) {
	config.SetCurrentConfig(config.ChainIDCYB)

	for _, k := range data {
		wif := k[0]
		hx := k[1]

		key, err := NewPrivateKeyFromWif(wif)
		if err != nil {
			assert.FailNow(t, err.Error(), "NewPrivateKeyFromWif")
		}

		assert.Equal(t, hx, key.ToHex())
	}
}

func TestMarshalUnmarshal(t *testing.T) {
	config.SetCurrentConfig(config.ChainIDCYB)

	for _, k := range data {
		wif := k[0]
		key1, err := NewPrivateKeyFromWif(wif)
		if err != nil {
			assert.FailNow(t, err.Error(), "NewPrivateKeyFromWif")
		}

		var buf bytes.Buffer
		enc := util.NewTypeEncoder(&buf)

		if err := key1.Marshal(enc); err != nil {
			assert.FailNow(t, err.Error(), "Marshal")
		}

		dec := util.NewTypeDecoder(&buf)
		var key2 PrivateKey
		if err := key2.Unmarshal(dec); err != nil {
			assert.FailNow(t, err.Error(), "Unmarshal")
		}

		assert.Equal(t, key2.Bytes(), key1.Bytes())
	}
}
