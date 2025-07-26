package mappers

/*
#include "../graal/dxfg_api.h"
#include <stdlib.h>
*/
import "C"

// Singleton instances to avoid map allocations
var (
	quoteMapper         = QuoteMapper{}
	timeAndSaleMapper   = TimeAndSaleMapper{}
	profileMapper       = ProfileMapper{}
	orderMapper         = OrderMapper{}
	spreadOrderMapper   = SpreadOrderMapper{}
	candleMapper        = CandleMapper{}
	tradeMapper         = TradeMapper{}
	tradeETHMapper      = TradeETHMapper{}
	analyticOrderMapper = AnalyticOrderMapper{}
)

// GetMapper returns the appropriate mapper singleton for a given event type
// No map allocation - just direct singleton access
func SelectMapper(eventType int32) MapperInterface {
	switch eventType {
	case C.DXFG_EVENT_QUOTE:
		return quoteMapper
	case C.DXFG_EVENT_TIME_AND_SALE:
		return timeAndSaleMapper
	case C.DXFG_EVENT_PROFILE:
		return profileMapper
	case C.DXFG_EVENT_ORDER:
		return orderMapper
	case C.DXFG_EVENT_SPREAD_ORDER:
		return spreadOrderMapper
	case C.DXFG_EVENT_CANDLE:
		return candleMapper
	case C.DXFG_EVENT_TRADE:
		return tradeMapper
	case C.DXFG_EVENT_TRADE_ETH:
		return tradeETHMapper
	case C.DXFG_EVENT_ANALYTIC_ORDER:
		return analyticOrderMapper
	default:
		return nil
	}
}
