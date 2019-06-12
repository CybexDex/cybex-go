package operations

//go:generate ffjson $GOFILE

import (
	"github.com/CybexDex/cybex-go/types"
	"github.com/CybexDex/cybex-go/util"
	"github.com/juju/errors"
)

func init() {
	types.OperationMap[types.OperationTypeVestingBalanceCreate] = func() types.Operation {
		op := &VestingBalanceCreateOperation{}
		return op
	}
}

type VestingBalanceCreateOperation struct {
	types.OperationFee
	Amount  types.AssetAmount   `json:"amount"`
	Creator types.GrapheneID    `json:"creator"`
	Owner   types.GrapheneID    `json:"owner"`
	Policy  types.VestingPolicy `json:"policy"`
}

func (p VestingBalanceCreateOperation) Type() types.OperationType {
	return types.OperationTypeVestingBalanceCreate
}

func (p VestingBalanceCreateOperation) Marshal(enc *util.TypeEncoder) error {
	if err := enc.Encode(int8(p.Type())); err != nil {
		return errors.Annotate(err, "encode OperationType")
	}

	if err := enc.Encode(p.Fee); err != nil {
		return errors.Annotate(err, "encode fee")
	}

	if err := enc.Encode(p.Creator); err != nil {
		return errors.Annotate(err, "encode Creator")
	}

	if err := enc.Encode(p.Owner); err != nil {
		return errors.Annotate(err, "encode Owner")
	}

	if err := enc.Encode(p.Amount); err != nil {
		return errors.Annotate(err, "encode Amount")
	}

	if err := enc.Encode(p.Policy); err != nil {
		return errors.Annotate(err, "encode Policy")
	}

	return nil
}

//NewVestingBalanceCreateOperation creates a new VestingBalanceCreateOperation
func NewVestingBalanceCreateOperation() *VestingBalanceCreateOperation {
	tx := VestingBalanceCreateOperation{}
	return &tx
}
