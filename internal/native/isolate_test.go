package native

import (
	"sync"
	"testing"
)

func TestIsolateCreation(t *testing.T) {
	if getOrCreateIsolate().ptr != getOrCreateIsolate().ptr {
		t.Errorf("Multiple calls to getOrCreateIsolate returned different isolates")
	}
}

func TestMultipleAttachIsolateThreadInSameThread(t *testing.T) {
	_ = dispatchOnIsolateThread(func(thread *isolateThread) error {
		return dispatchOnIsolateThread(func(nestedThread *isolateThread) error {
			// With worker pool, nested calls may be routed to different workers
			// so we can't guarantee the same thread instance. Instead, verify
			// that both threads are valid.
			if thread.ptr == nil || nestedThread.ptr == nil {
				t.Errorf("One or both thread instances are nil")
			}
			return nil
		})
	})
}

func TestDoubleDetachIsolateThread(t *testing.T) {
	thread := attachCurrentThread()
	thread.detach()
	thread.detach()
}

func TestConcurrentAttachDetachIsolateThread(t *testing.T) {
	var wg sync.WaitGroup
	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func() {
			_ = dispatchOnIsolateThread(func(thread *isolateThread) error {
				defer wg.Done()
				return nil
			})
		}()
	}
	wg.Wait()
}
