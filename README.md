# Sized Error Group

Sized Error Group is used to work through a queue of work using an [WaitGroup](https://golang.org/pkg/sync/#WaitGroup) with limits on the amount of work that will be done concurrently


Example:
``` golang
func TimeSinceAndWait(t time.Time) func() error {
	return func() error {
		fmt.Printf("func: time since start: %d seconds\n", int(time.Since(t).Seconds()))
		time.Sleep(time.Second)
		return nil
	}
}

func ExampleNewSizedErrGroup() {
	start := time.Now()
	g := NewSizedErrGroup(2) // Create a sized error group 
	for i := 0; i < 5; i++ {
		g.Go(TimeSinceAndWait(start)) // Initialize work to be done
	}

	err := g.Wait() // Wait for the work to be complete or for an error to return
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
```
