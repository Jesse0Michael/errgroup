package errgroup

import (
	"context"
	"fmt"
	"time"
)

func TimeSinceAndWait(t time.Time) func() error {
	return func() error {
		fmt.Printf("func: time since start: %d seconds\n", int(time.Since(t).Seconds()))
		time.Sleep(time.Second)
		return nil
	}
}

func IsEven(i int) func() error {
	return func() error {
		time.Sleep(time.Duration(i*500) * time.Millisecond)
		if i%2 == 0 {
			fmt.Printf("number '%d' is even\n", i)
			return nil
		}

		fmt.Printf("number '%d' is not even\n", i)
		return fmt.Errorf("number '%d' is not even", i)
	}
}

func ExampleNewSizedErrGroup() {
	start := time.Now()
	g := NewSizedErrGroup(2)
	for i := 0; i < 5; i++ {
		g.Go(TimeSinceAndWait(start))
	}

	err := g.Wait()
	fmt.Printf("time since start: %d seconds\n", int(time.Since(start).Seconds()))
	fmt.Printf("err: %v\n", err)

	// Output:
	// func: time since start: 0 seconds
	// func: time since start: 0 seconds
	// func: time since start: 1 seconds
	// func: time since start: 1 seconds
	// func: time since start: 2 seconds
	// time since start: 3 seconds
	// err: <nil>
}

func ExampleNewSizedErrGroup_withFailure() {
	g := NewSizedErrGroup(0)
	for i := 0; i < 5; i++ {
		g.Go(IsEven(i))
		time.Sleep(100 * time.Millisecond)
	}

	err := g.Wait()
	fmt.Printf("err: %v\n", err)

	// Output:
	// number '0' is even
	// number '1' is not even
	// err: number '1' is not even
}

func ExampleWithContext_cancel() {
	start := time.Now()
	ctx, cancel := context.WithCancel(context.TODO())
	g, _ := WithContext(ctx, 2)
	for i := 0; i < 10; i++ {
		g.Go(TimeSinceAndWait(start))
	}

	go func() {
		time.Sleep(500 * time.Millisecond)
		cancel()
	}()

	err := g.Wait()
	fmt.Printf("err: %v\n", err)

	// Output:
	// func: time since start: 0 seconds
	// func: time since start: 0 seconds
	// err: wait group context cancelled
}
