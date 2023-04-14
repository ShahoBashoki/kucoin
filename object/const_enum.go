package object

//go:generate stringer -output=./const_enum_string.go -type=OrderStateType ./

type (
	// KlineTypeType is an enumeration.
	KlineTypeType string

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
	// KlineTypeType1min is KlineTypeType.
	KlineTypeType1min KlineTypeType = "1min"
	// KlineTypeType3min is a KlineTypeType.
	KlineTypeType3min KlineTypeType = "3min"
	// KlineTypeType5min is a KlineTypeType.
	KlineTypeType5min KlineTypeType = "5min"
	// KlineTypeType15min is a KlineTypeType.
	KlineTypeType15min KlineTypeType = "15min"
	// KlineTypeType30min is a KlineTypeType.
	KlineTypeType30min KlineTypeType = "30min"
	// KlineTypeType1hour is a KlineTypeType.
	KlineTypeType1hour KlineTypeType = "1hour"
	// KlineTypeType2hour is a KlineTypeType.
	KlineTypeType2hour KlineTypeType = "2hour"
	// KlineTypeType4hour is a KlineTypeType.
	KlineTypeType4hour KlineTypeType = "4hour"
	// KlineTypeType6hour is a KlineTypeType.
	KlineTypeType6hour KlineTypeType = "6hour"
	// KlineTypeType8hour is a KlineTypeType.
	KlineTypeType8hour KlineTypeType = "8hour"
	// KlineTypeType12hour is a KlineTypeType.
	KlineTypeType12hour KlineTypeType = "12hour"
	// KlineTypeType1day is a KlineTypeType.
	KlineTypeType1day KlineTypeType = "1day"
	// KlineTypeType1week is a KlineTypeType.
	KlineTypeType1week KlineTypeType = "1week"

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
