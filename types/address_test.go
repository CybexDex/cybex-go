package types

import (
	"testing"

	"github.com/CybexDex/cybex-go/config"
	"github.com/stretchr/testify/assert"
)

var addresses = []string{
	"CYBFN9r6VYzBK8EKtMewfNbfiGCr56pHDBFi",
}

func TestAddress(t *testing.T) {
	config.SetCurrentConfig(config.ChainIDCYB)

	for _, add := range addresses {
		address, err := NewAddressFromString(add)
		if err != nil {
			assert.FailNow(t, err.Error(), "NewAddressFromString")
		}

		assert.Equal(t, add, address.String())
	}
}
