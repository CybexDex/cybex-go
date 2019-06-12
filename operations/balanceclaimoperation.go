package operations

//go:generate ffjson $GOFILE

import (
	"github.com/CybexDex/cybex-go/types"
	"github.com/CybexDex/cybex-go/util"
	"github.com/juju/errors"
)

func init() {
	types.OperationMap[types.OperationTypeBalanceClaim] = func() types.Operation {
		op := &BalanceClaimOperation{}
		return op
	}
}

type BalanceClaimOperation struct {
	types.OperationFee
	BalanceToClaim   types.GrapheneID  `json:"balance_to_claim"`
	BalanceOwnerKey  types.PublicKey   `json:"balance_owner_key"`
	DepositToAccount types.GrapheneID  `json:"deposit_to_account"`
	TotalClaimed     types.AssetAmount `json:"total_claimed"`
}

func (p BalanceClaimOperation) Type() types.OperationType {
	return types.OperationTypeBalanceClaim
}

func (p BalanceClaimOperation) Marshal(enc *util.TypeEncoder) error {
	if err := enc.Encode(int8(p.Type())); err != nil {
		return errors.Annotate(err, "encode OperationType")
	}

	if err := enc.Encode(p.Fee); err != nil {
		return errors.Annotate(err, "encode fee")
	}

	if err := enc.Encode(p.DepositToAccount); err != nil {
		return errors.Annotate(err, "encode DepositToAccount")
	}

	if err := enc.Encode(p.BalanceToClaim); err != nil {
		return errors.Annotate(err, "encode BalanceToClaim")
	}

	if err := enc.Encode(p.BalanceOwnerKey); err != nil {
		return errors.Annotate(err, "encode BalanceOwnerKey")
	}

	if err := enc.Encode(p.TotalClaimed); err != nil {
		return errors.Annotate(err, "encode TotalClaimed")
	}

	return nil
}

//NewBalanceClaimOperation creates a new BalanceClaimOperation
func NewBalanceClaimOperation() *BalanceClaimOperation {
	tx := BalanceClaimOperation{}
	return &tx
}
