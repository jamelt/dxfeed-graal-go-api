package native

/*
#include "dxfg_api.h"
#include <stdlib.h>
*/
import "C"

import (
	"fmt"
	"unsafe"

	"github.com/dxfeed/dxfeed-graal-go-api/internal/native/mappers"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events"
)

type CMapper interface {
	// Define methods that your C types must satisfy
}

type ListMapper[T CMapper] struct {
	size     C.int32_t
	elements **T
}

func NewListMapper[T CMapper, U comparable](elements []U) *ListMapper[T] {
	size := len(elements)
	e := (**T)(C.malloc(C.size_t(size) * C.size_t(unsafe.Sizeof((*int)(nil)))))
	slice := unsafe.Slice(e, C.size_t(size))
	for i, element := range elements {
		slice[i] = allocElement[T, U](element)
	}

	return &ListMapper[T]{
		elements: e,
		size:     C.int32_t(size),
	}
}

func allocElement[T CMapper, U comparable](element U) *T {
	switch t := any(element).(type) {
	case int32:
		return (*T)(C.malloc(C.size_t(unsafe.Sizeof(element))))
	case events.EventType:
		// all market events have to implement this interface
		mapper := mappers.SelectMapper(int32(t.Type()))
		return (*T)(mapper.CEvent(t))
	default:
		symbol := eventMapper.cSymbol(t)
		if symbol != nil {
			return (*T)(symbol)
		} else {
			fmt.Printf("Couldn't alloc element for %T\n", element)
			return nil
		}

	}
}
