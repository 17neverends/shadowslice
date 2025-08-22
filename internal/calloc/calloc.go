package calloc

// #include <stdlib.h>
import "C"
import "unsafe"

// CreateSlice calculate need memory size and return her ptr
func CreateSlice[T any](n int) []T {
	if n == 0 {
		return nil
	}
	cPtr := C.calloc(C.size_t(n), C.size_t(unsafe.Sizeof(*(new(T)))))
	return unsafe.Slice((*T)(cPtr), n)
}

// FreeMemory reset useless data after slice swap
func FreeMemory(ptr unsafe.Pointer) {
	if ptr != nil {
		C.free(ptr)
	}
}
