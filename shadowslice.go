package main

import (
	"errors"
	"fmt"
	"sync"
	"unsafe"

	"github.com/17neverends/shadowslice/internal/calloc"
)

const growFactor = 2

type ShadowSlice[T any] struct {
	slice         []T
	growSlice     []T
	mu            sync.RWMutex
	currentIdx    int
	offset        int
	cAllocEnabled bool
}

func NewShadowSlice[T any](initSize int, cAllocEnabled bool) (*ShadowSlice[T], error) {
	if initSize%2 != 0 {
		return nil, errors.New("incorrect init size")
	}

	var slice []T
	var growSlice []T

	if cAllocEnabled {
		slice = calloc.CreateSlice[T](initSize)
		growSlice = calloc.CreateSlice[T](initSize * growFactor)
	} else {
		slice = make([]T, initSize)
		growSlice = make([]T, initSize*growFactor)
	}

	return &ShadowSlice[T]{
		slice:         slice,
		growSlice:     growSlice,
		cAllocEnabled: cAllocEnabled,
	}, nil
}

// Return zero value for T type
func (ss *ShadowSlice[T]) getZeroValue() T {
	var zeroValue T
	return zeroValue
}

// Append add data on current index pos in main slice and check to need relocate 2 values in grow slice
func (ss *ShadowSlice[T]) Append(val T) {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	if ss.currentIdx == len(ss.slice) {
		if ss.cAllocEnabled {
			var oldMemory unsafe.Pointer
			if len(ss.slice) > 0 {
				oldMemory = unsafe.Pointer(&ss.slice[0])
			}

			ss.slice = ss.growSlice
			ss.growSlice = calloc.CreateSlice[T](cap(ss.slice) * growFactor)

			go calloc.FreeMemory(oldMemory)
		} else {
			ss.slice = ss.growSlice
			ss.growSlice = make([]T, cap(ss.slice)*growFactor)
		}
		ss.offset = 0
	}

	ss.slice[ss.currentIdx] = val

	if (cap(ss.slice) / growFactor) <= ss.currentIdx {
		for i := 0; i < growFactor; i++ {
			ss.growSlice[ss.offset+i] = ss.slice[ss.offset+i]
		}
		ss.offset += growFactor
	}
	ss.currentIdx++
}

// Get return value, exist by index in main slice
func (ss *ShadowSlice[T]) Get(idx int) (T, bool) {
	ss.mu.RLock()
	defer ss.mu.RUnlock()
	if idx >= ss.currentIdx {
		return ss.getZeroValue(), false
	}

	return ss.slice[idx], true
}

// Modify replace value by index in main slice and addittionaly can replace data in grow slice
func (ss *ShadowSlice[T]) Modify(idx int, newValue T) {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	if idx >= ss.currentIdx {
		return
	}

	ss.slice[idx] = newValue

	if (cap(ss.slice) / growFactor) <= idx {
		ss.growSlice[ss.currentIdx-idx] = newValue
	}
}

// Cleanup do memory clear if memory was allocated with cAllocEnabled otherwise return error
func (ss *ShadowSlice[T]) Cleanup() error {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	if !ss.cAllocEnabled {
		return errors.New("unable to execute without cAllocEnabled setting")
	}

	if len(ss.slice) > 0 {
		calloc.FreeMemory(unsafe.Pointer(&ss.slice[0]))
	}
	if len(ss.growSlice) > 0 {
		calloc.FreeMemory(unsafe.Pointer(&ss.growSlice[0]))
	}

	return nil
}

// Len return 2 int with value length of 2 slices
func (ss *ShadowSlice[T]) Len() (int, int) {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	return len(ss.slice), len(ss.growSlice)
}

// Cap return 2 int with value capacity of 2 slices
func (ss *ShadowSlice[T]) Cap() (int, int) {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	return cap(ss.slice), cap(ss.growSlice)
}

// String return debug info
func (ss *ShadowSlice[T]) String() string {
	ss.mu.RLock()
	defer ss.mu.RUnlock()

	return fmt.Sprintf("slice: %v (len=%d, cap=%d)\ngrowSlice: %v (len=%d, cap=%d)\ncurrentIdx: %d",
		ss.slice, len(ss.slice), cap(ss.slice),
		ss.growSlice, len(ss.growSlice), cap(ss.growSlice),
		ss.currentIdx)
}
