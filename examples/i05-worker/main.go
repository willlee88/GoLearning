package main

import (
	"fmt"
	"sync"
	"time"
)

type Job struct {
	MatchID string
	Winner  string
}

func worker(id int, jobs <-chan Job, wg *sync.WaitGroup) {
	defer wg.Done()
	for j := range jobs {
		// pretend DB write
		time.Sleep(20 * time.Millisecond)
		fmt.Printf("worker=%d settled match=%s winner=%s\n", id, j.MatchID, j.Winner)
	}
}

func main() {
	jobs := make(chan Job, 8)
	var wg sync.WaitGroup
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go worker(i, jobs, &wg)
	}
	for i := 0; i < 5; i++ {
		jobs <- Job{MatchID: fmt.Sprintf("m%d", i), Winner: "Ada"}
	}
	close(jobs)
	wg.Wait()
	fmt.Println("drained")
}
