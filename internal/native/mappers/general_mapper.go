package mappers

/*
#include "../graal/dxfg_api.h"
#include <stdlib.h>
*/
import "C"

func convertString(value *C.char) *string {
	if value == nil {
		return nil
	} else {
		result := C.GoString(value)
		return &result
	}
}

func CString(str *string) *C.char {
	if str == nil {
		return nil
	}
	return C.CString(*str)
}
