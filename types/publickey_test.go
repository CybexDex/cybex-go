package types

import (
	"testing"

	"github.com/CybexDex/cybex-go/config"
	"github.com/juju/errors"
	"github.com/stretchr/testify/assert"
)

var keys = []string{
	"CYB6K35Bajw29N4fjP4XADHtJ7bEj2xHJ8CoY2P2s1igXTB5oMBhR",
	"CYB4txNeAoSWcDX7oWceKppMb956z5oRx6mQyCJXCUB7aUh1EJp5y",
	"CYB6iUXJDmAPNbHWHtDDcmPTQ6F3nMBqi6pUHdhSkzWNd6grob2JP",
	"CYB5KCRzL27VLBvhPJ1DaXViuUPxyEXjDvVtWaifUkouNr2MkMGSH",
	"CYB6ThjMq97v6dLQUAmdsZfWG9ENq8nghVUhmLMQi52MDqXvtRGNc",
	"CYB5zzvbDtkbUVU1gFFsKqCE55U7JbjTp6mTh1usFv7KGgXL7HDQk",
	"CYB5yXqEBShUgcVm7Mve8Fg4RzQ2ftPpmo77aMbz884eX9aeGVvwD",
	"CYB5Z3vsgH6xj6HLXcsU38yo4TyoZs9AUzpfbaXbuxsAYPbutWvEP",
	"CYB5shffTjVoT4J8Zrj3f2mQJw4UVKrnbx5FWYhVgov45EpBf2NYi",
	"CYB56fy8qpkLzNoguGMPgPNkkznxnx2woEg1qPq7E6gF2SeGSRyK5",
}

func TestNewPublicKey(t *testing.T) {
	config.SetCurrentConfig(config.ChainIDCYB)
	for _, k := range keys {
		key, err := NewPublicKeyFromString(k)
		if err != nil {
			assert.FailNow(t, errors.Annotate(err, "NewPublicKeyFromString").Error())
		}

		assert.Equal(t, k, key.String())
	}
}
