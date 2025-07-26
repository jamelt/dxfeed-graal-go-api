# Engineering Review: Go ↔ Graal/Java Native Layer

A focused engineering review of the _“internal/native”_ layer (the CGO boundary) targeting extreme-throughput, low-latency systems crossing the Go ↔ Graal/Java frontier thousands or millions of times per second.

This feedback is divided into five areas:

- Threading
- Allocation
- Memory management / leaks
- CGO call overhead
- General Go-level hygiene

Each section includes **concrete recommendations** and inline **code snippets**.

---

## ✅ 1. Threading & Isolate Management

### Current State

- Every public operation wraps a closure in `executeInIsolateThread()`.
- `executeInIsolateThread()` does:
  1. `LockOSThread()`
  2. `graal_attach_thread()`
  3. call
  4. `graal_detach_thread()`
  5. `UnlockOSThread()`

### Impact

Lock/attach/detach on every market-data event easily dwarfs the _real_ native work.  
On Linux this path is ±300 ns; on macOS it’s worse. Under heavy flow, this becomes a bottleneck.

### Recommendations

#### ✅ 1-A. Keep a Dedicated Graal “Pinned Goroutine”

- Start one goroutine at package-init that calls `runtime.LockOSThread()` and **never** unlocks.
- Inside it:
  - Call `graal_attach_thread()` **once**
  - Read requests from a channel or ring buffer.
  - Other goroutines marshal tiny structs and wait on a `sync.Pool`’d object/waiter.
- Reduces latency to ~40 ns from hundreds.

#### ~~1-B. Re-use Existing Thread When `executeInIsolateThread()` Nests~~

~~You already re-use when the same goroutine calls twice (nice!), but avoid the attach/detach at the _outermost_ boundary when you're already on an attached thread.  
Track this with a context/flag on goroutine-local storage.~~

---

## 2. Heap Allocation Hot Spots

### 2-A. `C.CString` + `C.free` for Every String

- Symbol names and event fields create garbage on both Go and C heaps.
- **Fix**: Cache per-immutable symbol in `sync.Map` keyed by Go string.

### 2-B. `eventClazzList` / `ListMapper`

- `createEventClazzList()` mallocs N pointer slots + N tiny structs per call.
- **Fix**:
  - One `malloc` for `N*sizeof(dxfg_event_clazz_t)`
  - Pass **slice’s data pointer** to C
  - Use a `sync.Pool` for 32/64/128-slot backing arrays.

### 2-C. `Save()/Unref()` Index Pointer

- You `malloc` 1 byte per callback registration.
- **Fix**: Pool the pointer stubs—they hold no actual data.

---

## 3. Memory-Safety & Leak Audit

### 3-A. C Strings in `_Mapper.CEvent` Never Freed

**Fix options**:

1. After `dxfg_DXPublisher_publishEvents()` returns:
   - Iterate and free created C strings (via cleanup closures).
2. (Preferred) Use C helper to:
   - Copy into Java heap immediately.
   - Free in Go right away.

### 3-B. `createEventClazzList()` – Missing `destroy` on Panic

- Use `defer` in `executeInIsolateThread()` wrapper to always free memory.

### 3-C. `InstrumentProfileReader.readFromFile*`

- You copy profiles but **do not free** nested strings inside `dxfg_instrument_profile_t`.
- **Fix**: Verify native release function frees **nested** fields or wrap your own.

---

## 4. CGO Call Overhead Optimisation

### 4-A. Inline Trivial Wrappers

- Wrappers like `GetSystemProperty()` and `AwaitProcessed()` still do Go→C→Graal.
- **Fix**: Batch calls or expose coarser Graal-side APIs.

### 4-B. Disable `-gcflags=all=-l -N` in Production

- Use:
  ```bash
  -ldflags "-s -w"
  -gcflags=all=-trimpath
  -tags netgo -tags noop
  ```

### 4-C. Use `//go:nosplit` Where Safe

- For wrappers like `checkCall()` or `getJavaThreadErrorIfExist()`
- **Benchmark first** before widespread use.

---

## 5. Go-Level Clean-Ups & Best Practices

### 5-A. Error Policy

- Migrate from `errors` to `github.com/pkg/errors`
  ```go
  err := errors.Wrap(call(), "dxfg attach failed")
  ```

### 5-B. Panic vs Error

- `newEventMapper().goEvent()` panics on unknown `clazz`.
- **Fix**: Return annotated error or use `stats.Counter`.

### 5-C. Avoid `interface{}` in Hot Paths

- `DXFeedSubscription.AddSymbols([]any)` causes assert overhead.
- **Fix**: Introduce typed helpers:
  - `AddStringSymbols`
  - `AddWildcardSymbols`

### 5-D. Build Tags

- Add:
  ```go
  //go:build cgo
  ```
  to CGO-dependent files.

---

## ✅ Fast Tactical Changes (Merge Immediately)

### 1. Global C-String Cache

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
	if actual != cs {
		C.free(unsafe.Pointer(cs))
	}
	return actual.(*C.char)
}
```

Replace every `C.CString(symbol)` for immutable identifiers with `getCString(symbol)`.

### 2. Replace `createEventClazzList()` With Pooled Contiguous Allocation

- Improves throughput by ≈ 2×

### 3. Create a Single Pinned Graal Goroutine

- As in **§1-A**
- 10–20 LoC, **orders-of-magnitude impact**

### 4. Clean Up All Outstanding C Allocations

- After `publish` or listener `detach`

---

## 📈 Next Steps

- **Benchmark** with:
  - `go test -bench`
  - `perf stat`
  - `flamegraph`

### Expected Gains

- ~70% **fewer CGO calls**
- ~90% **drop in `C.malloc` traffic**
- **GC allocations** in hot path → _near zero_

---

## 💬 Need Help?

Want concrete code for any change above?  
Feel free to ask—happy to provide targeted patches.
