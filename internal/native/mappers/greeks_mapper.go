package mappers

/*
#include "../graal/dxfg_api.h"
#include <stdlib.h>
*/
import "C"

import (
	"unsafe"

	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/greeks"
)

type GreeksMapper struct{}

func (GreeksMapper) CEvent(event interface{}) unsafe.Pointer {
	greeksEvent := event.(*greeks.Greeks)

	g := (*C.dxfg_greeks_t)(C.malloc(C.size_t(unsafe.Sizeof(C.dxfg_greeks_t{}))))
	g.market_event.event_type.clazz = C.DXFG_EVENT_GREEKS
	g.market_event.event_symbol = C.CString(*greeksEvent.EventSymbol())
	g.market_event.event_time = C.int64_t(greeksEvent.EventTime())
	g.event_flags = C.int32_t(greeksEvent.EventFlags())
	g.index = C.int64_t(greeksEvent.Index())
	g.price = C.double(greeksEvent.Price())
	g.volatility = C.double(greeksEvent.Volatility())
	g.delta = C.double(greeksEvent.Delta())
	g.gamma = C.double(greeksEvent.Gamma())
	g.theta = C.double(greeksEvent.Theta())
	g.rho = C.double(greeksEvent.Rho())
	g.vega = C.double(greeksEvent.Vega())
	return unsafe.Pointer(g)
}

func (GreeksMapper) GoEvent(native unsafe.Pointer) interface{} {
	greeksNative := (*C.dxfg_greeks_t)(native)
	g := greeks.NewGreeks(C.GoString(greeksNative.market_event.event_symbol))
	g.SetEventSymbol(C.GoString(greeksNative.market_event.event_symbol))
	g.SetEventTime(int64(greeksNative.market_event.event_time))
	g.SetEventFlags(int32(greeksNative.event_flags))
	g.SetIndex(int64(greeksNative.index))
	g.SetPrice(float64(greeksNative.price))
	g.SetVolatility(float64(greeksNative.volatility))
	g.SetDelta(float64(greeksNative.delta))
	g.SetGamma(float64(greeksNative.gamma))
	g.SetTheta(float64(greeksNative.theta))
	g.SetRho(float64(greeksNative.rho))
	g.SetVega(float64(greeksNative.vega))
	return g
}
