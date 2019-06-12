package operations

//go:generate ffjson $GOFILE

import (
	"github.com/CybexDex/cybex-go/types"
	"github.com/CybexDex/cybex-go/util"
	"github.com/juju/errors"
)

func init() {
	types.OperationMap[types.OperationTypeAssetUpdate] = func() types.Operation {
		op := &AssetUpdateOperation{}
		return op
	}
}

type AssetUpdateOperation struct {
	types.OperationFee
	AssetToUpdate types.GrapheneID   `json:"asset_to_update"`
	Issuer        types.GrapheneID   `json:"issuer"`
	Extensions    types.Extensions   `json:"extensions"`
	NewIssuer     *types.GrapheneID  `json:"new_issuer"`
	NewOptions    types.AssetOptions `json:"new_options"`
}

func (p AssetUpdateOperation) Type() types.OperationType {
	return types.OperationTypeAssetUpdate
}

func (p AssetUpdateOperation) Marshal(enc *util.TypeEncoder) error {
	if err := enc.Encode(int8(p.Type())); err != nil {
		return errors.Annotate(err, "encode OperationType")
	}

	if err := enc.Encode(p.Fee); err != nil {
		return errors.Annotate(err, "encode Fee")
	}

	if err := enc.Encode(p.Issuer); err != nil {
		return errors.Annotate(err, "encode Issuer")
	}

	if err := enc.Encode(p.AssetToUpdate); err != nil {
		return errors.Annotate(err, "encode AssetToUpdate")
	}

	if err := enc.Encode(p.NewIssuer != nil); err != nil {
		return errors.Annotate(err, "encode have NewIssuer")
	}

	if err := enc.Encode(p.NewIssuer); err != nil {
		return errors.Annotate(err, "NewIssuer")
	}

	if err := enc.Encode(p.NewOptions); err != nil {
		return errors.Annotate(err, "encode new NewOptions")
	}

	if err := enc.Encode(p.Extensions); err != nil {
		return errors.Annotate(err, "encode Extensions")
	}

	return nil
}

//NewAssetUpdateOperation creates a new AssetUpdateOperation
func NewAssetUpdateOperation() *AssetUpdateOperation {
	tx := AssetUpdateOperation{
		Extensions: types.Extensions{},
	}
	return &tx
}
