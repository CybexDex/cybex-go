package types

import (
	"github.com/CybexDex/cybex-go/util"
	"github.com/juju/errors"
)

type Fund struct {
	AccountName Buffer `json:"accountName"`
	//Asset Buffer `json:"asset"`
	//FundType Buffer `json:"fundType"`
	Size UInt32 `json:"size"`
	//Offset UInt32 `json:"offset"`
	Expiration int `json:"expiration"`
}

func (p Fund) Marshal(enc *util.TypeEncoder) error {
	if err := enc.Encode(p.AccountName); err != nil {
		return errors.Annotate(err, "encode RefBlockNum")
	}

	//if err := enc.Encode(p.Asset); err != nil {
	//	return errors.Annotate(err, "encode RefBlockPrefix")
	//}
	//
	//if err := enc.Encode(p.FundType); err != nil {
	//	return errors.Annotate(err, "encode Expiration")
	//}
	//
	if err := enc.Encode(p.Size); err != nil {
		return errors.Annotate(err, "encode Operations")
	}

	//if err := enc.Encode(p.Offset); err != nil {
	//	return errors.Annotate(err, "encode Extension")
	//}
	if err := enc.Encode(p.Expiration); err != nil {
		return errors.Annotate(err, "encode Expiration")
	}
	return nil
}
