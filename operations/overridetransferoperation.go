package operations

//go:generate ffjson $GOFILE

import (
	"github.com/CybexDex/cybex-go/types"
	"github.com/CybexDex/cybex-go/util"
	"github.com/juju/errors"
)

func init() {
	types.OperationMap[types.OperationTypeOverrideTransfer] = func() types.Operation {
		op := &OverrideTransferOperation{}
		return op
	}
}

type OverrideTransferOperation struct {
	types.OperationFee
	Amount     types.AssetAmount `json:"amount"`
	Extensions types.Extensions  `json:"extensions"`
	From       types.GrapheneID  `json:"from"`
	Issuer     types.GrapheneID  `json:"issuer"`
	Memo       *types.Memo       `json:"memo,omitempty"`
	To         types.GrapheneID  `json:"to"`
}

func (p OverrideTransferOperation) Type() types.OperationType {
	return types.OperationTypeOverrideTransfer
}

func (p OverrideTransferOperation) Marshal(enc *util.TypeEncoder) error {
	if err := enc.Encode(int8(p.Type())); err != nil {
		return errors.Annotate(err, "encode OperationType")
	}

	if err := enc.Encode(p.Fee); err != nil {
		return errors.Annotate(err, "encode fee")
	}

	if err := enc.Encode(p.Issuer); err != nil {
		return errors.Annotate(err, "encode Issuer")
	}

	if err := enc.Encode(p.From); err != nil {
		return errors.Annotate(err, "encode From")
	}

	if err := enc.Encode(p.To); err != nil {
		return errors.Annotate(err, "encode To")
	}

	if err := enc.Encode(p.Amount); err != nil {
		return errors.Annotate(err, "encode Amount")
	}

	if err := enc.Encode(p.Memo != nil); err != nil {
		return errors.Annotate(err, "encode have Memo")
	}

	if err := enc.Encode(p.Memo); err != nil {
		return errors.Annotate(err, "encode Memo")
	}

	if err := enc.Encode(p.Extensions); err != nil {
		return errors.Annotate(err, "encode extensions")
	}

	return nil
}

//NewOverrideTransferOperation creates a new OverrideTransferOperation
func NewOverrideTransferOperation() *OverrideTransferOperation {
	tx := OverrideTransferOperation{
		Extensions: types.Extensions{},
	}
	return &tx
}
