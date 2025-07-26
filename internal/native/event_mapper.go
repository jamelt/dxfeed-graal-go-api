package native

/*
#include "graal/dxfg_api.h"
#include <stdlib.h>
*/
import "C"

import (
	"fmt"
	"unsafe"

	"github.com/dxfeed/dxfeed-graal-go-api/internal/native/mappers"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api/Osub"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events"
)

type eventMapperUtil int

const eventMapper = eventMapperUtil(0)

func (m eventMapperUtil) goEvents(eventsList *C.dxfg_event_type_list) []interface{} {
	if eventsList == nil || eventsList.elements == nil || int(eventsList.size) == 0 {
		return nil
	}

	size := int(eventsList.size)
	list := make([]interface{}, size)
	elementsSlice := unsafe.Slice(eventsList.elements, C.size_t(eventsList.size))

	for i, event := range elementsSlice {
		list[i] = m.goEvent(event)
	}

	return list
}

func (m eventMapperUtil) goEvent(event *C.dxfg_event_type_t) interface{} {
	mapper := mappers.SelectMapper(int32(event.clazz))
	if mapper != nil {
		return mapper.GoEvent(unsafe.Pointer(event))
	} else {
		panic(fmt.Sprintf("unknown event eventcodes %v", event.clazz))
	}
}

// TODO add recursive release for symbols
func (m eventMapperUtil) cSymbol(symbol any) unsafe.Pointer {
	switch value := symbol.(type) {
	case string:
		return unsafe.Pointer(m.cStringSymbol(value))
	case *Osub.WildcardSymbol:
		return unsafe.Pointer(m.cWildCardSymbol())
	case *Osub.IndexedEventSubscriptionSymbol:
		return unsafe.Pointer(m.cIndexedEventSubscriptionSymbol(value.Symbol(), value.Source()))
	case *Osub.TimeSeriesSubscriptionSymbol:
		return unsafe.Pointer(m.cTimeSeriesSymbol(value.Symbol(), value.FromTime()))
	default:
		return nil
	}
}

func (m eventMapperUtil) cStringSymbol(str string) *dxfg_symbol_t {
	ss := &dxfg_symbol_t{}
	ss.t = 0
	ss.symbol = C.CString(str)
	return ss
}

func (m eventMapperUtil) cWildCardSymbol() *dxfg_symbol_t {
	ss := &dxfg_symbol_t{}
	ss.t = 2
	return ss
}

func (m eventMapperUtil) cTimeSeriesSymbol(str any, fromTime int64) *dxfg_time_series_subscription_symbol_t {
	ss := &dxfg_time_series_subscription_symbol_t{}
	ss.t = 4
	ss.symbol = (*dxfg_symbol_t)(m.cSymbol(str))
	ss.from_time = C.int64_t(fromTime)
	return ss
}

func (m eventMapperUtil) cIndexedEventSubscriptionSymbol(str any, source events.IndexedEventSourceInterface) *dxfg_indexed_event_subscription_symbol_t {
	ss := &dxfg_indexed_event_subscription_symbol_t{}
	ss.t = 3
	ss.symbol = (*dxfg_symbol_t)(m.cSymbol(str))
	nativeSource := &dxfg_indexed_event_source_t{}
	nativeSource.id = C.int32_t(source.Id())
	nativeSource.name = C.CString(*source.Name())
	switch source.Type() {
	case events.IndexedEventSourceType:
		nativeSource.t = C.INDEXED_EVENT_SOURCE
	case events.OrderSourceType:
		nativeSource.t = C.ORDER_SOURCE
	default:
		panic(fmt.Sprintf("Undefined source %d", source.Type()))
	}
	ss.source = nativeSource
	return ss
}
