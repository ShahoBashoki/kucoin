package util

import (
	"github.com/ShahoBashoki/kucoin/object"
)

// KlineTypeToSecond is a function.
func KlineTypeToSecond(
	klineType object.KlineTypeType,
) int64 {
	switch klineType {
	case object.KlineTypeType1min:
		return object.NUM1MinToSecond

	case object.KlineTypeType3min:
		return object.NUM3MinToSecond

	case object.KlineTypeType5min:
		return object.NUM5MinToSecond

	case object.KlineTypeType15min:
		return object.NUM15MinToSecond

	case object.KlineTypeType30min:
		return object.NUM30MinToSecond

	case object.KlineTypeType1hour:
		return object.NUM1HourToSecond

	case object.KlineTypeType2hour:
		return object.NUM2HourToSecond

	case object.KlineTypeType4hour:
		return object.NUM4HourToSecond

	case object.KlineTypeType6hour:
		return object.NUM6HourToSecond

	case object.KlineTypeType8hour:
		return object.NUM8HourToSecond

	case object.KlineTypeType12hour:
		return object.NUM12HourToSecond

	case object.KlineTypeType1day:
		return object.NUM1DayToSecond

	case object.KlineTypeType1week:
		return object.NUM1WeekToSecond

	default:
		return object.NUM5MinToSecond
	}
}
