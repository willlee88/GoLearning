package main

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"
)

func main() {
	g, ctx := errgroup.WithContext(context.Background())

	g.Go(func() error {
		select {
		case <-time.After(200 * time.Millisecond):
			fmt.Println("A done")
			return nil
		case <-ctx.Done():
			fmt.Println("A canceled")
			return ctx.Err()
		}
	})
	g.Go(func() error {
		time.Sleep(50 * time.Millisecond)
		return fmt.Errorf("B failed")
	})

	if err := g.Wait(); err != nil {
		fmt.Println("group:", err)
	}
}
