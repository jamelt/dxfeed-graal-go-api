package native

/*
#cgo CFLAGS: -I${SRCDIR}/graal
#include "dxfg_api.h"
*/
import "C"

import (
	"runtime"
	"sync"
)

type isolate struct {
	ptr *C.graal_isolate_t
}

var (
	// once guard and global isolate pointer
	isolateOnce     sync.Once
	isolateInstance *isolate

	// channel used to marshal function calls onto the dedicated pinned goroutine
	isolateReqCh chan isolateRequest

	// pointer to the isolate thread owned by the worker goroutine
	workerThreadPtr *C.graal_isolatethread_t
)

// isolateRequest wraps a closure together with a channel to transmit the result
type isolateRequest struct {
	fn   func(thread *isolateThread) error
	done chan error
}

func init() {
	// Create request channel with some buffering to reduce contention
	isolateReqCh = make(chan isolateRequest, 1024)

	go isolateWorker()
}

// isolateWorker runs forever on a dedicated OS thread attached to the Graal isolate.
// It executes all incoming requests sequentially, thus avoiding costly attach/detach
// operations for every native call.
func isolateWorker() {
	// Pin this goroutine to its current OS thread for the entire lifetime.
	runtime.LockOSThread()

	iso := getOrCreateIsolate()

	// Attach this thread **once** to the isolate.
	var threadPtr *C.graal_isolatethread_t
	err := checkIsolateCall(func() C.int {
		return C.graal_attach_thread(iso.ptr, &threadPtr)
	})
	if err != nil {
		panic(err)
	}

	workerThread := &isolateThread{ptr: threadPtr, shouldDetach: false}

	// Save for re-entrancy detection.
	workerThreadPtr = threadPtr

	for req := range isolateReqCh {
		req.done <- req.fn(workerThread)
	}

	// Not expected to reach here under normal conditions, but detach gracefully.
	_ = checkIsolateCall(func() C.int {
		return C.graal_detach_thread(workerThread.ptr)
	})
}

func getOrCreateIsolate() *isolate {
	isolateOnce.Do(func() {
		isolateInstance = &isolate{}
		err := checkIsolateCall(func() C.int {
			return C.graal_create_isolate(nil, &isolateInstance.ptr, nil)
		})
		if err != nil {
			panic(err)
		}
	})
	return isolateInstance
}

type isolateThread struct {
	ptr          *C.graal_isolatethread_t
	shouldDetach bool
}

func executeInIsolateThread(call func(thread *isolateThread) error) error {
	// Fast-path for re-entrant calls originating from the worker goroutine itself.
	if C.graal_get_current_thread(getOrCreateIsolate().ptr) == workerThreadPtr {
		// Re-entrant invocation from within the worker goroutine. We can cheaply obtain the
		// current thread handle (already attached) without the cost of attach/detach but we
		// MUST balance the LockOSThread() done inside attachCurrentThread with an Unlock in
		// thread.detach() to keep the lock count correct.
		thread := attachCurrentThread()
		err := call(thread)
		thread.detach()
		return err
	}

	done := make(chan error, 1)
	isolateReqCh <- isolateRequest{fn: call, done: done}
	return <-done
}

func attachCurrentThread() *isolateThread {
	runtime.LockOSThread()
	isolate := getOrCreateIsolate()
	thread := &isolateThread{ptr: C.graal_get_current_thread(isolate.ptr), shouldDetach: false}
	if thread.ptr == nil {
		err := checkIsolateCall(func() C.int {
			return C.graal_attach_thread(isolate.ptr, &thread.ptr)
		})
		if err != nil {
			panic(err)
		}
		thread.shouldDetach = true
	}
	return thread
}

func (t *isolateThread) detach() {
	defer runtime.UnlockOSThread()
	if t.ptr != nil && t.shouldDetach {
		err := checkIsolateCall(func() C.int {
			return C.graal_detach_thread(t.ptr)
		})
		if err != nil {
			panic(err)
		}
	}
	t.ptr = nil
}
