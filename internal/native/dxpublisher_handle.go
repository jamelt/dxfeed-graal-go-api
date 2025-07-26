package native

/*
#include "graal/dxfg_api.h"
#include <stdlib.h>
*/
import "C"

import (
	"unsafe"
)

type DXPublisherHandle struct {
	handle Handler
}

func NewDXPublisherHandle(ptr *C.dxfg_publisher_t) *DXPublisherHandle {
	return &DXPublisherHandle{handle: NewJavaHandle(unsafe.Pointer(ptr))}
}

func (p *DXPublisherHandle) Free() error {
	if p != nil {
		return p.handle.Free()
	}
	return nil
}

// Publish publishes events to the DXFeed infrastructure.
//
// WARNING: This method has a known memory leak. Each string field in the published
// events will leak memory because the C API frees the event structures but not
// the string fields allocated by C.CString(). This is a limitation of the current
// C API design which provides no callback mechanism to free these strings.
//
// This leak primarily affects applications that publish events. Applications that
// only consume/receive events are not affected by this issue.
//
// If you must publish events, consider:
// - Using a limited set of symbols to bound the leak (with string interning)
// - Restarting the publisher periodically
// - Waiting for a future API version that addresses this issue
func (p *DXPublisherHandle) Publish(events []interface{}) error {
	err := dispatchOnIsolateThread(func(thread *isolateThread) error {
		l := NewListMapper[C.dxfg_event_type_list, interface{}](events)
		_ = C.dxfg_DXPublisher_publishEvents(thread.ptr, p.ptr(), (*C.dxfg_event_type_list)(unsafe.Pointer(l)))
		return nil
	})
	return err
}

func (p *DXPublisherHandle) ptr() *C.dxfg_publisher_t {
	return (*C.dxfg_publisher_t)(p.handle.Ptr())
}
