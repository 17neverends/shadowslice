package calloc

import (
	"testing"
	"unsafe"
)

func TestCreateSlice(t *testing.T) {
	t.Run("create slice of int", func(t *testing.T) {
		slice := CreateSlice[int](5)
		defer FreeMemory(unsafe.Pointer(&slice[0]))

		if len(slice) != 5 {
			t.Errorf("Expected length 5, got %d", len(slice))
		}
		if cap(slice) != 5 {
			t.Errorf("Expected capacity 5, got %d", cap(slice))
		}

		for i := 0; i < 5; i++ {
			if slice[i] != 0 {
				t.Errorf("Expected zero value at index %d, got %d", i, slice[i])
			}
		}
	})

	t.Run("create slice of string", func(t *testing.T) {
		slice := CreateSlice[string](3)
		defer FreeMemory(unsafe.Pointer(&slice[0]))

		if len(slice) != 3 {
			t.Errorf("Expected length 3, got %d", len(slice))
		}

		for i := 0; i < 3; i++ {
			if slice[i] != "" {
				t.Errorf("Expected empty string at index %d, got '%s'", i, slice[i])
			}
		}
	})

	t.Run("create slice of struct", func(t *testing.T) {
		type TestStruct struct {
			Value int
			Name  string
		}

		slice := CreateSlice[TestStruct](2)
		defer FreeMemory(unsafe.Pointer(&slice[0]))

		if len(slice) != 2 {
			t.Errorf("Expected length 2, got %d", len(slice))
		}

		for i := 0; i < 2; i++ {
			if slice[i].Value != 0 {
				t.Errorf("Expected Value 0 at index %d, got %d", i, slice[i].Value)
			}
			if slice[i].Name != "" {
				t.Errorf("Expected empty Name at index %d, got '%s'", i, slice[i].Name)
			}
		}
	})

	t.Run("create zero length slice", func(t *testing.T) {
		slice := CreateSlice[int](0)
		if slice != nil {
			t.Errorf("Expected nil slice for zero length, got %v", slice)
		}
	})

	t.Run("create slice and modify values", func(t *testing.T) {
		slice := CreateSlice[int](3)
		defer FreeMemory(unsafe.Pointer(&slice[0]))

		slice[0] = 10
		slice[1] = 20
		slice[2] = 30

		if slice[0] != 10 || slice[1] != 20 || slice[2] != 30 {
			t.Errorf("Expected [10, 20, 30], got %v", slice)
		}
	})
}

func TestFreeMemory(t *testing.T) {
	t.Run("free nil pointer", func(t *testing.T) {
		FreeMemory(nil)
	})

	t.Run("free valid pointer", func(t *testing.T) {
		slice := CreateSlice[int](5)

		ptr := unsafe.Pointer(&slice[0])

		FreeMemory(ptr)
	})
}

func TestCreateSliceEdgeCases(t *testing.T) {
	t.Run("large slice", func(t *testing.T) {
		slice := CreateSlice[byte](10000)
		defer FreeMemory(unsafe.Pointer(&slice[0]))

		if len(slice) != 10000 {
			t.Errorf("Expected length 10000, got %d", len(slice))
		}

		for i := range slice {
			slice[i] = byte(i % 256)
		}

		for i := range slice {
			if slice[i] != byte(i%256) {
				t.Errorf("Expected %d at index %d, got %d", i%256, i, slice[i])
			}
		}
	})

	t.Run("slice of pointers", func(t *testing.T) {
		slice := CreateSlice[*int](4)
		defer FreeMemory(unsafe.Pointer(&slice[0]))

		if len(slice) != 4 {
			t.Errorf("Expected length 4, got %d", len(slice))
		}

		for i := 0; i < 4; i++ {
			if slice[i] != nil {
				t.Errorf("Expected nil pointer at index %d, got %v", i, slice[i])
			}
		}
	})
}

func TestMemoryIsolation(t *testing.T) {
	t.Run("independent slices", func(t *testing.T) {
		slice1 := CreateSlice[int](3)
		defer FreeMemory(unsafe.Pointer(&slice1[0]))

		slice2 := CreateSlice[int](3)
		defer FreeMemory(unsafe.Pointer(&slice2[0]))

		slice1[0] = 100
		slice2[0] = 200

		if slice1[0] != 100 {
			t.Errorf("Expected slice1[0] = 100, got %d", slice1[0])
		}
		if slice2[0] != 200 {
			t.Errorf("Expected slice2[0] = 200, got %d", slice2[0])
		}
	})
}

func BenchmarkCreateSlice(b *testing.B) {
	b.Run("small slice", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			slice := CreateSlice[int](10)
			FreeMemory(unsafe.Pointer(&slice[0]))
		}
	})

	b.Run("medium slice", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			slice := CreateSlice[int](1000)
			FreeMemory(unsafe.Pointer(&slice[0]))
		}
	})

	b.Run("large slice", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			slice := CreateSlice[byte](100000)
			FreeMemory(unsafe.Pointer(&slice[0]))
		}
	})
}

func BenchmarkCreateSliceVsMake(b *testing.B) {
	b.Run("calloc", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			slice := CreateSlice[int](1000)
			FreeMemory(unsafe.Pointer(&slice[0]))
		}
	})

	b.Run("make", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			slice := make([]int, 1000)
			_ = slice
		}
	})
}
