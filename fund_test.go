package funding

import (
	"sync"
	"testing"
)

const WORKERS = 10

func BenchmarkWithdrawals(b *testing.B) {
	// Skip N = 1
	if b.N < WORKERS {
		return
	}

	// Add as many dollars as we have iterations this run
	fund := NewFund(b.N)

	// Casually assume b.N divides cleanly
	dollarsPerFounder := b.N / WORKERS

	// WaitGroup structs don't need to be initialized
	// (their "zero value" is ready to use).
	// So, we just declare one and then use it.
	var wg sync.WaitGroup

	for i := 0; i < WORKERS; i++ {
		// Use waitgroup abstraction of the semaphore tool to manage goroutines
		wg.Add(1)

		// Spawn off a founder worker, as a closure
		go func() {
			// At the end of the clojure call Done on the WaitGroup semaphore
			defer wg.Done()

			for i := 0; i < dollarsPerFounder; i++ {
				fund.Withdraw(1)
			}

		}()
	}

	// Wait for all the workers to finish
	wg.Wait()

	if fund.Balance() != 0 {
		b.Error("Balance wasn't zero:", fund.Balance())
	}
}
