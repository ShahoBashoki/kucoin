package object

//go:generate stringer -output=./const_enum_string.go -type=OrderStateType ./

type (
	// OrderSideType is an enumeration.
	OrderSideType string

	// OrderStateType is an enumeration.
	OrderStateType string

	// OrderTradeTypeType is an enumeration.
	OrderTradeTypeType string

	// OrderTypeType is an enumeration.
	OrderTypeType string
)

const (
	// OrderSideTypeBuy is OrderSideType.
	OrderSideTypeBuy OrderSideType = "buy"
	// OrderSideTypeSell is a OrderSideType.
	OrderSideTypeSell OrderSideType = "sell"

	// OrderStateTypeActive is OrderStateType.
	OrderStateTypeActive OrderStateType = "active"
	// OrderStateTypeDone is a OrderStateType.
	OrderStateTypeDone OrderStateType = "done"

	// OrderTypeTypeLimit is OrderTypeType.
	OrderTypeTypeLimit OrderTypeType = "limit"
	// OrderTypeTypeLimitStop is a OrderTypeType.
	OrderTypeTypeLimitStop OrderTypeType = "limit_stop"
	// OrderTypeTypeMarket is a OrderTypeType.
	OrderTypeTypeMarket OrderTypeType = "market"
	// OrderTypeTypeMarketStop is a OrderTypeType.
	OrderTypeTypeMarketStop OrderTypeType = "market_stop"

	// OrderTypeTypeTrade is OrderTypeType.
	OrderTypeTypeTrade OrderTypeType = "TRADE"
	// OrderTypeTypeMarginIsolatedTrade is a OrderTypeType.
	OrderTypeTypeMarginIsolatedTrade OrderTypeType = "MARGIN_ISOLATED_TRADE"
	// OrderTypeTypeMarginTrade is OrderTypeType.
	OrderTypeTypeMarginTrade OrderTypeType = "MARGIN_TRADE"
)
