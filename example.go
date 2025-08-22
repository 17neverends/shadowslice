package main

import (
	"fmt"
	"sync"
)

func main() {
	ss, err := NewShadowSlice[int](4, true)

	if err != nil {
		fmt.Println("Bug")
	}

	ss.Append(1)
	fmt.Printf("After append: %s\n", ss.String())

	ss.Modify(0, 2)
	fmt.Printf("After modify: %s\n", ss.String())

	val, ok := ss.Get(0)

	if !ok {
		fmt.Println("Bug")
	}

	fmt.Printf("Get after modify: %v\n", val)

	wg := sync.WaitGroup{}
	iter := 10
	wg.Add(iter)

	for i := 0; i < iter; i++ {
		go func(val int) {
			defer wg.Done()
			ss.Append(val)
		}(i)
	}

	wg.Wait()

	fmt.Printf("Wg wait: %s\n", ss.String())
}
