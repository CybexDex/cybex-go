package types

//go:generate ffjson $GOFILE

type SettleOrders []SettleOrder

type SettleOrder struct {
	ID             GrapheneID  `json:"id"`
	Owner          GrapheneID  `json:"owner"`
	SettlementDate Time        `json:"settlement_date"`
	Balance        AssetAmount `json:"balance"`
}
