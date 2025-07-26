package native

/*
#cgo CFLAGS: -I${SRCDIR}/graal
#include "graal/dxfg_api.h"
*/
import "C"

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"sync"
	"sync/atomic"
)

const (
	// DefaultWorkerCount is the default number of worker goroutines
	DefaultWorkerCount = 2

	// WorkerEnvVar is the environment variable name for configuring worker count
	WorkerEnvVar = "DXFEED_ISOLATE_WORKERS"

	// ChannelBufferSize is the buffer size for each worker's request channel
	ChannelBufferSize = 1024
)

// isolate represents a GraalVM isolate instance
type isolate struct {
	ptr *C.graal_isolate_t
}

// isolateThread represents a thread attached to the GraalVM isolate
type isolateThread struct {
	ptr          *C.graal_isolatethread_t
	shouldDetach bool // whether this thread should be detached when done
}

// request encapsulates a function to be executed within the isolate
// along with a channel to communicate the result
type request struct {
	fn         func(*isolateThread) error
	resultChan chan error
}

// workerPool manages the pool of worker goroutines that execute isolate operations
type workerPool struct {
	workers      int
	requestChans []chan request
	requestCount uint64 // for round-robin load balancing
}

var (
	// Global isolate instance (singleton)
	isolateOnce     sync.Once
	isolateInstance *isolate

	// Global worker pool
	pool *workerPool
)

func init() {
	workerCount := getConfiguredWorkerCount()
	pool = newWorkerPool(workerCount)
	pool.start()
}

// getConfiguredWorkerCount determines the number of workers based on environment
// variable or defaults. The count is capped at runtime.NumCPU().
func getConfiguredWorkerCount() int {
	// Try to read from environment variable
	if envValue := os.Getenv(WorkerEnvVar); envValue != "" {
		if count, err := strconv.Atoi(envValue); err == nil && count > 0 {
			return min(count, runtime.NumCPU())
		}
	}

	// Use default, but ensure it doesn't exceed available CPUs
	return min(DefaultWorkerCount, runtime.NumCPU())
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// newWorkerPool creates a new worker pool with the specified number of workers
func newWorkerPool(workerCount int) *workerPool {
	p := &workerPool{
		workers:      workerCount,
		requestChans: make([]chan request, workerCount),
	}

	// Initialize request channels
	for i := 0; i < workerCount; i++ {
		p.requestChans[i] = make(chan request, ChannelBufferSize)
	}

	return p
}

// start launches all worker goroutines
func (p *workerPool) start() {
	for i := 0; i < p.workers; i++ {
		go p.runWorker(i, p.requestChans[i])
	}
}

// submit sends a request to the next worker in round-robin fashion
func (p *workerPool) submit(fn func(*isolateThread) error) error {
	// Select worker using round-robin
	workerIndex := int(atomic.AddUint64(&p.requestCount, 1) % uint64(p.workers))

	// Create request with result channel
	req := request{
		fn:         fn,
		resultChan: make(chan error, 1),
	}

	// Submit request and wait for result
	p.requestChans[workerIndex] <- req
	return <-req.resultChan
}

// runWorker is the main loop for a worker goroutine
func (p *workerPool) runWorker(id int, requests chan request) {
	// Pin this goroutine to its OS thread for the entire lifetime
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	// Get or create the isolate
	iso := getOrCreateIsolate()

	// Attach this OS thread to the isolate once
	thread, err := attachThreadToIsolate(iso)
	if err != nil {
		panic(fmt.Errorf("worker %d: failed to attach thread to isolate: %w", id, err))
	}

	// Ensure graceful cleanup on exit
	defer func() {
		if err := detachThreadFromIsolate(thread); err != nil {
			// Log error but don't panic during cleanup
			fmt.Fprintf(os.Stderr, "worker %d: failed to detach thread: %v\n", id, err)
		}
	}()

	// Process requests until channel is closed
	for req := range requests {
		req.resultChan <- req.fn(thread)
	}
}

// getOrCreateIsolate returns the singleton isolate instance, creating it if necessary
func getOrCreateIsolate() *isolate {
	isolateOnce.Do(func() {
		isolateInstance = &isolate{}
		err := checkIsolateCall(func() C.int {
			return C.graal_create_isolate(nil, &isolateInstance.ptr, nil)
		})
		if err != nil {
			panic(fmt.Errorf("failed to create GraalVM isolate: %w", err))
		}
	})
	return isolateInstance
}

// attachThreadToIsolate attaches the current OS thread to the given isolate
func attachThreadToIsolate(iso *isolate) (*isolateThread, error) {
	var threadPtr *C.graal_isolatethread_t
	err := checkIsolateCall(func() C.int {
		return C.graal_attach_thread(iso.ptr, &threadPtr)
	})
	if err != nil {
		return nil, fmt.Errorf("failed to attach thread to isolate: %w", err)
	}

	return &isolateThread{
		ptr:          threadPtr,
		shouldDetach: false, // Worker threads don't detach
	}, nil
}

// detachThreadFromIsolate detaches a thread from the isolate
func detachThreadFromIsolate(thread *isolateThread) error {
	if thread == nil || thread.ptr == nil {
		return nil
	}

	return checkIsolateCall(func() C.int {
		return C.graal_detach_thread(thread.ptr)
	})
}

// dispatchOnIsolateThread executes the given function within the GraalVM isolate context.
// The function is executed by one of the worker threads in the pool.
// This function is kept for backward compatibility.
func dispatchOnIsolateThread(fn func(*isolateThread) error) error {
	return pool.submit(fn)
}

// attachCurrentThread attaches the current OS thread to the isolate.
// This is primarily used for testing and special cases where direct
// thread attachment is needed. The caller must ensure proper cleanup
// by calling the thread's detach method.
// This function is kept for backward compatibility.
func attachCurrentThread() *isolateThread {
	runtime.LockOSThread()

	iso := getOrCreateIsolate()
	thread := &isolateThread{
		ptr:          C.graal_get_current_thread(iso.ptr),
		shouldDetach: false,
	}

	// If not already attached, attach now
	if thread.ptr == nil {
		err := checkIsolateCall(func() C.int {
			return C.graal_attach_thread(iso.ptr, &thread.ptr)
		})
		if err != nil {
			runtime.UnlockOSThread()
			panic(fmt.Errorf("failed to attach current thread: %w", err))
		}
		thread.shouldDetach = true
	}

	return thread
}

// detach detaches this thread from the isolate if it was attached
// by attachCurrentThread. This also unlocks the OS thread.
// This method is kept for backward compatibility.
func (t *isolateThread) detach() {
	defer runtime.UnlockOSThread()

	if t.ptr != nil && t.shouldDetach {
		err := checkIsolateCall(func() C.int {
			return C.graal_detach_thread(t.ptr)
		})
		if err != nil {
			// Don't panic during cleanup, just log
			fmt.Fprintf(os.Stderr, "failed to detach thread: %v\n", err)
		}
	}
	t.ptr = nil
}
