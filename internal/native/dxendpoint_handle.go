package native

/*
#include "graal/dxfg_api.h"
#include <stdlib.h>
extern void OnStateChanged(graal_isolatethread_t *thread, dxfg_endpoint_state_t old_state, dxfg_endpoint_state_t new_state, void *user_data);
*/
import "C"

import (
	"errors"
	"sync"
	"unsafe"

	"github.com/dxfeed/dxfeed-graal-go-api/pkg/common"
)

type DXEndpointHandle struct {
	self     Handler
	feedOnce sync.Once
	feed     *DXFeedHandle

	publisherOnce sync.Once
	publisher     *DXPublisherHandle
}

func NewDXEndpointHandle(role common.Role) (*DXEndpointHandle, error) {
	var ptr *C.dxfg_endpoint_t
	err := dispatchOnIsolateThread(func(thread *isolateThread) error {
		return checkCall(func() {
			ptr = C.dxfg_DXEndpoint_create2(thread.ptr, (C.dxfg_endpoint_role_t)(role))
		})
	})
	if err != nil {
		return nil, err
	}

	return &DXEndpointHandle{self: NewJavaHandle(unsafe.Pointer(ptr))}, nil
}

func NewDXEndpointHandleWithProperties(role common.Role, properties map[string]string) (*DXEndpointHandle, error) {
	var ptr *C.dxfg_endpoint_t
	err := dispatchOnIsolateThread(func(thread *isolateThread) error {
		return checkCall(func() {
			builder := C.dxfg_DXEndpoint_newBuilder(thread.ptr)
			C.dxfg_DXEndpoint_Builder_withRole(thread.ptr, builder, (C.dxfg_endpoint_role_t)(role))
			for key, value := range properties {
				C.dxfg_DXEndpoint_Builder_withProperty(thread.ptr, builder, C.CString(key), C.CString(value))
			}

			ptr = C.dxfg_DXEndpoint_Builder_build(thread.ptr, builder)
		})
	})
	if err != nil {
		return nil, err
	}

	return &DXEndpointHandle{self: NewJavaHandle(unsafe.Pointer(ptr))}, nil
}

func (e *DXEndpointHandle) Close() error {
	err := dispatchOnIsolateThread(func(thread *isolateThread) error {
		return checkCall(func() {
			C.dxfg_DXEndpoint_close(thread.ptr, e.ptr())
		})
	})
	return errors.Join(err, e.Free())
}

func (e *DXEndpointHandle) CloseAndAwaitTermination() error {
	err := dispatchOnIsolateThread(func(thread *isolateThread) error {
		return checkCall(func() {
			C.dxfg_DXEndpoint_closeAndAwaitTermination(thread.ptr, e.ptr())
		})
	})
	return errors.Join(err, e.Free())
}

func (e *DXEndpointHandle) Connect(address string) error {
	return dispatchOnIsolateThread(func(thread *isolateThread) error {
		addressPtr := C.CString(address)
		defer C.free(unsafe.Pointer(addressPtr))

		return checkCall(func() {
			C.dxfg_DXEndpoint_connect(thread.ptr, e.ptr(), addressPtr)
		})
	})
}

func (e *DXEndpointHandle) Reconnect() error {
	return dispatchOnIsolateThread(func(thread *isolateThread) error {
		return checkCall(func() {
			C.dxfg_DXEndpoint_reconnect(thread.ptr, e.ptr())
		})
	})
}

func (e *DXEndpointHandle) Disconnect() error {
	return dispatchOnIsolateThread(func(thread *isolateThread) error {
		return checkCall(func() {
			C.dxfg_DXEndpoint_disconnect(thread.ptr, e.ptr())
		})
	})
}

func (e *DXEndpointHandle) DisconnectAndClear() error {
	return dispatchOnIsolateThread(func(thread *isolateThread) error {
		return checkCall(func() {
			C.dxfg_DXEndpoint_disconnectAndClear(thread.ptr, e.ptr())
		})
	})
}

func (e *DXEndpointHandle) AwaitProcessed() error {
	return dispatchOnIsolateThread(func(thread *isolateThread) error {
		return checkCall(func() {
			C.dxfg_DXEndpoint_awaitProcessed(thread.ptr, e.ptr())
		})
	})
}

func (e *DXEndpointHandle) AwaitNotConnected() error {
	return dispatchOnIsolateThread(func(thread *isolateThread) error {
		return checkCall(func() {
			C.dxfg_DXEndpoint_awaitNotConnected(thread.ptr, e.ptr())
		})
	})
}

func (e *DXEndpointHandle) GetFeed() (*DXFeedHandle, error) {
	var err error
	e.feedOnce.Do(func() {
		var ptr *C.dxfg_feed_t
		err = dispatchOnIsolateThread(func(thread *isolateThread) error {
			return checkCall(func() {
				ptr = C.dxfg_DXEndpoint_getFeed(thread.ptr, e.ptr())
			})
		})
		e.feed = NewDXFeedHandle(ptr)
	})

	return e.feed, err
}

func (e *DXEndpointHandle) GetPublisher() (*DXPublisherHandle, error) {
	var err error
	e.publisherOnce.Do(func() {
		var ptr *C.dxfg_publisher_t
		err = dispatchOnIsolateThread(func(thread *isolateThread) error {
			return checkCall(func() {
				ptr = C.dxfg_DXEndpoint_getPublisher(thread.ptr, e.ptr())
			})
		})
		e.publisher = NewDXPublisherHandle(ptr)
	})

	return e.publisher, err
}

//export OnStateChanged
func OnStateChanged(thread *C.graal_isolatethread_t, old C.dxfg_endpoint_state_t, new C.dxfg_endpoint_state_t, userData unsafe.Pointer) {
	Restore(userData).(common.ConnectionStateListener).UpdateState(common.ConnectionState(old), common.ConnectionState(new))
}

func (e *DXEndpointHandle) AttachListener(listener common.ConnectionStateListener) error {
	err := dispatchOnIsolateThread(func(thread *isolateThread) error {
		l := C.dxfg_PropertyChangeListener_new(thread.ptr, (*[0]byte)(C.OnStateChanged), Save(listener))
		C.dxfg_DXEndpoint_addStateChangeListener(thread.ptr, e.ptr(), l)
		return nil
	})
	return err
}

func (e *DXEndpointHandle) Free() error {
	return errors.Join(e.feed.Free(), e.publisher.Free(), e.self.Free())
}

func (e *DXEndpointHandle) ptr() *C.dxfg_endpoint_t {
	return (*C.dxfg_endpoint_t)(e.self.Ptr())
}
