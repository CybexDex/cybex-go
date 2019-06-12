package operations

//go:generate ffjson $GOFILE

import (
	"github.com/CybexDex/cybex-go/types"
	"github.com/CybexDex/cybex-go/util"
	"github.com/juju/errors"
)

func init() {
	types.OperationMap[types.OperationTypeProposalCreate] = func() types.Operation {
		op := &ProposalCreateOperation{}
		return op
	}
}

type ProposalCreateOperation struct {
	types.OperationFee
	ExpirationTime      types.Time                     `json:"expiration_time"`
	Extensions          types.Extensions               `json:"extensions"`
	FeePayingAccount    types.GrapheneID               `json:"fee_paying_account"`
	ReviewPeriodSeconds *types.UInt32                  `json:"review_period_seconds,omitempty"`
	ProposedOps         types.OperationEnvelopeHolders `json:"proposed_ops"`
}

func (p ProposalCreateOperation) Type() types.OperationType {
	return types.OperationTypeProposalCreate
}

func (p ProposalCreateOperation) Marshal(enc *util.TypeEncoder) error {
	if err := enc.Encode(int8(p.Type())); err != nil {
		return errors.Annotate(err, "encode OperationType")
	}

	if err := enc.Encode(p.Fee); err != nil {
		return errors.Annotate(err, "encode Fee")
	}

	if err := enc.Encode(p.FeePayingAccount); err != nil {
		return errors.Annotate(err, "encode FeePayingAccount")
	}

	if err := enc.Encode(p.ExpirationTime); err != nil {
		return errors.Annotate(err, "encode ExpirationTime")
	}

	if err := enc.Encode(p.ProposedOps); err != nil {
		return errors.Annotate(err, "encode ProposedOps")
	}

	if err := enc.Encode(p.ReviewPeriodSeconds != nil); err != nil {
		return errors.Annotate(err, "encode have ReviewPeriodSeconds")
	}

	if err := enc.Encode(p.ReviewPeriodSeconds); err != nil {
		return errors.Annotate(err, "encode ReviewPeriodSeconds")
	}

	if err := enc.Encode(p.Extensions); err != nil {
		return errors.Annotate(err, "encode Extensions")
	}

	return nil
}

//NewProposalCreateOperation creates a new ProposalCreateOperation
func NewProposalCreateOperation() *ProposalCreateOperation {
	tx := ProposalCreateOperation{
		Extensions: types.Extensions{},
	}
	return &tx
}
