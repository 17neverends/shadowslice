package main

import (
	"testing"
)

func TestNewShadowSlice(t *testing.T) {
	t.Run("with C allocation", func(t *testing.T) {
		ss, err := NewShadowSlice[int](10, true)
		if err != nil {
			t.Errorf("err should be nil")
		}

		defer func() {
			err := ss.Cleanup()

			if err != nil {
				t.Errorf("Expected nil cleanup error")

			}
		}()

		sliceLen, growSliceLen := ss.Len()
		if sliceLen != 10 || growSliceLen != 20 {
			t.Errorf("Expected lengths (10, 20), got (%d, %d)", sliceLen, growSliceLen)
		}

		sliceCap, growSliceCap := ss.Cap()
		if sliceCap != 10 || growSliceCap != 20 {
			t.Errorf("Expected capacities (10, 20), got (%d, %d)", sliceCap, growSliceCap)
		}
	})

	t.Run("without C allocation", func(t *testing.T) {
		ss, err := NewShadowSlice[string](10, false)
		if err != nil {
			t.Errorf("err should be nil")
		}

		sliceLen, growSliceLen := ss.Len()
		if sliceLen != 10 || growSliceLen != 20 {
			t.Errorf("Expected lengths (10, 20), got (%d, %d)", sliceLen, growSliceLen)
		}
	})
}

func TestAppendAndGet(t *testing.T) {
	ss, err := NewShadowSlice[int](2, false)
	if err != nil {
		t.Errorf("err should be nil")
	}

	ss.Append(10)
	ss.Append(20)
	ss.Append(30)

	if val, ok := ss.Get(0); !ok || val != 10 {
		t.Errorf("Expected 10, got %v (ok: %v)", val, ok)
	}
	if val, ok := ss.Get(1); !ok || val != 20 {
		t.Errorf("Expected 20, got %v (ok: %v)", val, ok)
	}
	if val, ok := ss.Get(2); !ok || val != 30 {
		t.Errorf("Expected 30, got %v (ok: %v)", val, ok)
	}
}

func TestModify(t *testing.T) {
	ss, err := NewShadowSlice[string](2, false)
	if err != nil {
		t.Errorf("err should be nil")
	}

	ss.Append("hello")
	ss.Append("world")
	ss.Modify(0, "modified")

	if val, ok := ss.Get(0); !ok || val != "modified" {
		t.Errorf("Expected 'modified', got %v", val)
	}
}

func TestCleanup(t *testing.T) {
	t.Run("with C allocation", func(t *testing.T) {
		ss, err := NewShadowSlice[float64](6, true)
		if err != nil {
			t.Errorf("err should be nil")
		}

		err = ss.Cleanup()
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	})

	t.Run("without C allocation", func(t *testing.T) {
		ss, err := NewShadowSlice[float64](6, false)
		if err != nil {
			t.Errorf("err should be nil")
		}

		err = ss.Cleanup()
		if err == nil {
			t.Error("Expected error but got nil")
		}
	})
}

func TestBoundaryConditions(t *testing.T) {
	ss, err := NewShadowSlice[int](2, false)
	if err != nil {
		t.Errorf("err should be nil")
	}

	if _, ok := ss.Get(100); ok {
		t.Error("Expected false for out of bounds access")
	}

	ss.Modify(100, 42)

	for i := 0; i < 100; i++ {
		ss.Append(i)
	}

	for i := 0; i < 100; i++ {
		if val, ok := ss.Get(i); !ok || val != i {
			t.Errorf("Expected %d, got %v (ok: %v)", i, val, ok)
		}
	}
}

func TestConcurrentAccess(t *testing.T) {
	ss, err := NewShadowSlice[int](10, false)
	if err != nil {
		t.Errorf("err should be nil")
	}

	done := make(chan bool)

	go func() {
		for i := 0; i < 1000; i++ {
			ss.Append(i)
		}
		done <- true
	}()

	go func() {
		for i := 0; i < 1000; i++ {
			ss.Get(i % 100)
		}
		done <- true
	}()

	<-done
	<-done
}
