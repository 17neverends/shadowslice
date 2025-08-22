# ShadowSlice

A high-performance, hybrid-memory-managed slice implementation for Go. It's designed to minimize garbage collection pauses and provide efficient amortized appends through a proactive two-buffer strategy.

## Concept

The standard Go slice causes GC pauses when growing, as it requires allocating a new larger array and copying all elements. ShadowSlice solves this with:

1. **Two-Buffer System**:
   - **Main Slice**: Active slice for operations
   - **Growth Slice**: Pre-allocated larger slice (2x capacity)

2. **Proactive Copying**: Elements are gradually copied to the growth slice during appends, amortizing the copy cost

3. **Instant Switching**: When main slice fills, switch to pre-prepared growth slice is immediate

4. **C Allocation Option**: Can use C's malloc/free to bypass Go GC entirely

## Features

- **Amortized O(1) appends** with consistent performance
- **Reduced GC pressure** through smart memory management
- **Thread-safe** operations with RWMutex
- **Generic** implementation works with any type


# API Reference

## Constructor

| Method | Description | Parameters | Return Value |
| :--- | :--- | :--- | :--- |
| `NewShadowSlice[T any](initSize int, cAllocEnabled bool) *ShadowSlice[T]` | Creates a new ShadowSlice with specified initial size and allocation mode. | `initSize` - initial capacity of the main slice<br>`cAllocEnabled` - if `true`, uses C malloc allocation | Pointer to the created `ShadowSlice` instance |

## Core Methods

| Method | Description | Parameters | Return Value |# API Reference

## Constructor

| Method | Description | Parameters | Return Value |
| :--- | :--- | :--- | :--- |
| `NewShadowSlice[T any](initSize int, cAllocEnabled bool) *ShadowSlice[T]` | Creates a new ShadowSlice with specified initial size and allocation mode. | `initSize` - initial capacity of the main slice<br>`cAllocEnabled` - if `true`, uses C malloc allocation | Pointer to the created `ShadowSlice` instance |

## Core Methods

| Method | Description | Parameters | Return Value |
| :--- | :--- | :--- | :--- |
| `Append(val T)` | Appends a value, handling internal growth logic. | `val` - value to append | - |
| `Get(idx int) (T, bool)` | Retrieves value by index. | `idx` - element index | `(value, true)` on success, `(zeroValue, false)` on error |
| `Modify(idx int, newValue T)` | Updates value at specified index. | `idx` - element index<br>`newValue` - new value | - |
| `Cleanup() error` | Frees C-allocated memory. Required when using C allocation. | - | `error` if called without C allocation enabled |

## Utility Methods

| Method | Description | Parameters | Return Value |
| :--- | :--- | :--- | :--- |
| `Len() (int, int)` | Returns lengths of internal slices. | - | `(mainSliceLen, growSliceLen)` |
| `Cap() (int, int)` | Returns capacities of internal slices. | - | `(mainSliceCap, growSliceCap)` |
| `String() string` | Returns debug information about internal state. | - | String with state description |# API Reference

## Constructor

| Method | Description | Parameters | Return Value |
| :--- | :--- | :--- | :--- |
| `NewShadowSlice[T any](initSize int, cAllocEnabled bool) *ShadowSlice[T]` | Creates a new ShadowSlice with specified initial size and allocation mode. | `initSize` - initial capacity of the main slice<br>`cAllocEnabled` - if `true`, uses C malloc allocation | Pointer to the created `ShadowSlice` instance |

## Core Methods

| Method | Description | Parameters | Return Value |
| :--- | :--- | :--- | :--- |
| `Append(val T)` | Appends a value, handling internal growth logic. | `val` - value to append | - |
| `Get(idx int) (T, bool)` | Retrieves value by index. | `idx` - element index | `(value, true)` on success, `(zeroValue, false)` on error |
| `Modify(idx int, newValue T)` | Updates value at specified index. | `idx` - element index<br>`newValue` - new value | - |
| `Cleanup() error` | Frees C-allocated memory. Required when using C allocation. | - | `error` if called without C allocation enabled |

## Utility Methods

| Method | Description | Parameters | Return Value |
| :--- | :--- | :--- | :--- |
| `Len() (int, int)` | Returns lengths of internal slices. | - | `(mainSliceLen, growSliceLen)` |
| `Cap() (int, int)` | Returns capacities of internal slices. | - | `(mainSliceCap, growSliceCap)` |
| `String() string` | Returns debug information about internal state. | - | String with state description |
| :--- | :--- | :--- | :--- |
| `Append(val T)` | Appends a value, handling internal growth logic. | `val` - value to append | - |
| `Get(idx int) (T, bool)` | Retrieves value by index. | `idx` - element index | `(value, true)` on success, `(zeroValue, false)` on error |
| `Modify(idx int, newValue T)` | Updates value at specified index. | `idx` - element index<br>`newValue` - new value | - |
| `Cleanup() error` | Frees C-allocated memory. Required when using C allocation. | - | `error` if called without C allocation enabled |

## Utility Methods

| Method | Description | Parameters | Return Value |
| :--- | :--- | :--- | :--- |
| `Len() (int, int)` | Returns lengths of internal slices. | - | `(mainSliceLen, growSliceLen)` |
| `Cap() (int, int)` | Returns capacities of internal slices. | - | `(mainSliceCap, growSliceCap)` |
| `String() string` | Returns debug information about internal state. | - | String with state description |

# Overall: why named shadow?
![2d0da3b0eb24078e7434676b9eb63d58](https://github.com/user-attachments/assets/0eb7cfe5-9921-4714-b0ec-ed4373119312)
