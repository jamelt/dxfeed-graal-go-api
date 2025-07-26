Below is a focused engineering review of the ‟internal/native” layer (the CGO boundary) from the perspective of extreme-throughput, low-latency systems that cross the Go ↔ Graal/Java frontier thousands or millions of times per second. I have divided the feedback into five areas: threading, allocation, memory management / leaks, CGO call overhead, and general Go-level hygiene. Each item includes concrete, actionable recommendations and, where appropriate, references to small code snippets you can drop in immediately.

────────────────────────────────────────────────────────

- [x] 1. Threading & isolate management
     ────────────────────────────────────────────────────────
     Current state
     • Every public operation wraps a closure in executeInIsolateThread().
     • executeInIsolateThread() does (1) LockOSThread(), (2) graal_attach_thread(), (3) call, (4) graal_detach_thread(), (5) UnlockOSThread().

Impact
Lock/attach/detach on every market-data event easily dwarfs the _real_ native work you do. On Linux this path is ±300 ns; on macOS it’s worse. Under heavy flow it becomes the bottleneck.

Recommendations

- [X]1-A Keep a dedicated Graal “pinned goroutine”
  • Start one goroutine at package-init that immediately calls runtime.LockOSThread() and _never_ unlocks.
  • Inside it, do the graal_attach_thread() **once** and then read requests from a channel.  
   • All other goroutines just marshal tiny structs onto that chan (or a fast ring-buffer) and wait on a sync.Pool’d object / waiter.  
   • Latency drops to ~40 ns (queueing cost) instead of hundreds.

~~1-B Re-use existing thread when executeInIsolateThread() nests*
You already re-use when the same goroutine calls twice (nice!), but avoid the attach/detach at the \_outermost* boundary when you are already on an attached thread. Track this with a context/flag on goroutine-local storage.~~

2. Heap allocation hot spots
   ────────────────────────────────────────────────────────
   2-A C.CString + C.free for every string
   • Symbol names, event fields, etc. create garbage on _both_ the Go and C heaps.
   • Allocate once per immutable symbol and cache in a sync.Map keyed by Go string.  
    The native side is read-only, so reuse is safe.

2-B eventClazzList / ListMapper
• createEventClazzList() mallocs N pointer slots + N tiny structs for _every_ call.
• Replace with:
– One malloc for N\*sizeof(dxfg_event_clazz_t)  
 – Pass the **slice’s data pointer** to C; the C API only needs the ints, not individual heap objects.
– Keep a sync.Pool of an 32/64/128-slot backing arrays; reuse.

2-C Save()/Unref() index pointer
• You malloc 1 byte for every callback registration. Pool these pointer stubs; they never hold data.

──────────────────────────────────────────────────────── 3. Memory-safety & leak audit
────────────────────────────────────────────────────────
3-A C strings created in all _Mapper.CEvent_ helpers are never freed.
That is a permanent leak whenever you publish events. Two fixes:

     Option 1:  After dxfg_DXPublisher_publishEvents() returns, iterate over the list and free any C strings you created.  That means ListMapper must capture clean-up closures along with the pointer.

     Option 2:  Provide a small C helper that copies the strings into Java heap immediately, then free in Go right away (preferred—minimises cross-heap lifetime).

3-B createEventClazzList() – missing destroy if graal call panics.
Use a defer inside executeInIsolateThread() wrapper so the allocated C memory is freed even on panic.

3-C InstrumentProfileReader.readFromFile\* – you copy profiles to Go but do not free individual C strings inside dxfg_instrument_profile_t. Verify the native release function actually frees nested fields, otherwise wrap your own.

──────────────────────────────────────────────────────── 4. CGO call overhead optimisation
────────────────────────────────────────────────────────
4-A Inline trivial wrappers
Tiny wrappers like GetSystemProperty() or AwaitProcessed() still cross CGO twice (Go→C, C→Graal). Batch related calls or expose coarse-grained APIs on the Graal side to minimise round-trips.

4-B Build with `-gcflags=all=-l -N` OFF in production
Ensure your Makefile releases are built with `-ldflags "-s -w"`, `-gcflags=all=-trimpath`, `-tags netgo`, and `-tags=noop`. Stripping helps i-cache pressure.

4-C Enable `//go:nosplit` where recursion is impossible
A few leaf wrappers (e.g. checkCall(), getJavaThreadErrorIfExist()) can be marked `//go:nosplit` to save stack-split checks. Benchmark first; don’t over-use.

──────────────────────────────────────────────────────── 5. Go-level clean-ups & best practice gaps
────────────────────────────────────────────────────────
5-A Error policy
User rule prefers github.com/pkg/errors but the layer uses the stdlib `errors`. Migrate:
err := errors.Wrap(call(), "dxfg attach failed")
everywhere; keeps caller stack.

5-B Panic vs error
newEventMapper().goEvent() panics on unknown clazz. In a live feed, a single unexpected event will tear the entire process down. Return an annotated error (or at least stats.Counter) instead.

5-C Avoid reflection / interface{} in hot path
DXFeedSubscription.AddSymbols([]any) forces interface{} assert for each element. Introduce typed ﹤ T ﹥ helpers (`AddStringSymbols`, `AddWildcardSymbols`, etc.) so the common case is zero-cost.

5-D Build tags
Add `//go:build cgo` to files that require CGO to prevent accidental `go test` on systems without CGO enabled.

────────────────────────────────────────────────────────
Fast tactical changes you can merge immediately
────────────────────────────────────────────────────────

1. Introduce a global symbol C-string cache:

```go
// internal/native/cstr.go
package native

import "C"
import (
	"sync"
	"unsafe"
)

var (
	cstrPool sync.Map // map[string]*C.char
)

// getCString returns a *C.char living on C heap that never changes.
func getCString(s string) *C.char {
	if v, ok := cstrPool.Load(s); ok {
		return v.(*C.char)
	}
	cs := C.CString(s)
	actual, _ := cstrPool.LoadOrStore(s, cs)
	// If another goroutine got here first, free ours.
	if actual != cs {
		C.free(unsafe.Pointer(cs))
	}
	return actual.(*C.char)
}
```

Then replace every `C.CString(symbol)` for immutable identifiers with `getCString(symbol)` (e.g. in all _Mapper.CEvent_ builders).

2. Replace createEventClazzList() with pooled contiguous allocation (≈ 2× throughput). Same for ListMapper.

3. Create a single pinned Graal goroutine (see §1-A). 10–20 lines of code but **orders-of-magnitude** impact under heavy load.

4. Clean up all outstanding C allocations after publish / listener detach.

────────────────────────────────────────────────────────
Next steps
────────────────────────────────────────────────────────
• Benchmark again (use go test -bench plus perf stat / flamegraph). You should see:
– ~70 % fewer CGO calls
– ~90 % drop in C.malloc traffic
– GC allocations in the hot path close to zero.
• After the low-hanging fruit, profile the Graal side—often string interning or reflection there dominates once the Go side is lean.

Feel free to ask for concrete code for any single change above; happy to provide targeted patches.
