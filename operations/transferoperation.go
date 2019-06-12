package operations

//go:generate ffjson $GOFILE

import (
	"github.com/CybexDex/cybex-go/types"
	"github.com/CybexDex/cybex-go/util"
	"github.com/juju/errors"
)

func init() {
	types.OperationMap[types.OperationTypeTransfer] = func() types.Operation {
		op := &TransferOperation{}
		return op
	}
}

type TransferOperation struct {
	types.OperationFee
	From       types.GrapheneID  `json:"from"`
	To         types.GrapheneID  `json:"to"`
	Amount     types.AssetAmount `json:"amount"`
	Memo       *types.Memo       `json:"memo,omitempty"`
	Extensions types.Extensions  `json:"extensions"`
}

func (p TransferOperation) Type() types.OperationType {
	return types.OperationTypeTransfer
}

func (p TransferOperation) Marshal(enc *util.TypeEncoder) error {
	if err := enc.Encode(int8(p.Type())); err != nil {
		return errors.Annotate(err, "encode OperationType")
	}

	if err := enc.Encode(p.Fee); err != nil {
		return errors.Annotate(err, "encode fee")
	}

	if err := enc.Encode(p.From); err != nil {
		return errors.Annotate(err, "encode from")
	}

	if err := enc.Encode(p.To); err != nil {
		return errors.Annotate(err, "encode to")
	}

	if err := enc.Encode(p.Amount); err != nil {
		return errors.Annotate(err, "encode amount")
	}

	if err := enc.Encode(p.Memo != nil); err != nil {
		return errors.Annotate(err, "encode have Memo")
	}

	if err := enc.Encode(p.Memo); err != nil {
		return errors.Annotate(err, "encode memo")
	}

	if err := enc.Encode(p.Extensions); err != nil {
		return errors.Annotate(err, "encode extensions")
	}

	return nil
}

//NewTransferOperation creates a new TransferOperation
func NewTransferOperation() *TransferOperation {
	tx := TransferOperation{
		Extensions: types.Extensions{},
	}
	return &tx
}
