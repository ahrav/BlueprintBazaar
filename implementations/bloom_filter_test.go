package implementations

import (
	"fmt"
	"sync"
	"testing"
)

// TestBasicBloomFilter tests the Add and Test methods of BasicBloomFilter.
func TestBasicBloomFilter(t *testing.T) {
	// Define test cases
	testCases := []struct {
		value       []byte
		expectFound bool
	}{
		{[]byte("hello"), true},
		{[]byte("world"), true},
		{[]byte("not-added"), false},
	}

	bf := NewBasicBloomFilter(WithCapacity(1000), WithFalsePositiveRate(0.01))

	falsePositives := 0
	totalNotAdded := 0

	for _, tc := range testCases {
		found, err := bf.Test(tc.value)
		if err != nil {
			t.Errorf("Test failed for value %s: %v", tc.value, err)
			continue
		}

		if !tc.expectFound && found {
			falsePositives++
		}

		if !tc.expectFound {
			totalNotAdded++
		}
	}

	// Calculate the false positive rate.
	if totalNotAdded > 0 {
		fpRate := float64(falsePositives) / float64(totalNotAdded)
		if fpRate > 0.01 { // Adjust this threshold based on your acceptable false positive rate
			t.Errorf("False positive rate too high: got %v, want at most 0.01", fpRate)
		}
	}
}

func TestBasicBloomFilterConcurrent(t *testing.T) {
	const (
		numElementsToAdd = 100_000
		numRandomTests   = 100_000
	)

	bf := NewBasicBloomFilter(WithCapacity(10_000_000), WithFalsePositiveRate(0.01))
	var wg1 sync.WaitGroup

	for i := 0; i < numElementsToAdd; i++ {
		wg1.Add(1)
		go func(i int) {
			defer wg1.Done()
			value := []byte(fmt.Sprintf("value-%d", i))
			err := bf.Add(value)
			if err != nil {
				t.Errorf("Add failed for value %s: %v", value, err)
			}
		}(i)
	}
	wg1.Wait()

	var wg sync.WaitGroup
	for i := 0; i < numElementsToAdd; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			value := []byte(fmt.Sprintf("value-%d", i))
			found, err := bf.Test(value)
			if err != nil || !found {
				t.Errorf("Test failed for added value %s; expected true, got %v, error: %v", value, found, err)
			}
		}(i)
	}

	falsePositives := 0
	for i := 0; i < numRandomTests; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			value := []byte(fmt.Sprintf("random-value-%d", i))
			found, _ := bf.Test(value)
			if found {
				falsePositives++
			}
		}(i)
	}
	wg.Wait()

	// Evaluate false positive rate.
	fpRate := float64(falsePositives) / float64(numRandomTests)
	if fpRate > 0.01 {
		t.Errorf("False positive rate too high: got %v, want at most 0.01", fpRate)
	}
}

// BenchmarkBasicBloomFilterAdd benchmarks the Add method of the BasicBloomFilter.
func BenchmarkBasicBloomFilterAdd(b *testing.B) {
	bf := NewBasicBloomFilter(WithCapacity(1000000), WithFalsePositiveRate(0.01))

	b.ResetTimer() // Reset the timer to exclude the setup time
	for n := 0; n < b.N; n++ {
		value := []byte(fmt.Sprintf("value-%d", n))
		bf.Add(value)
	}
}

// BenchmarkBasicBloomFilterTest benchmarks the Test method of the BasicBloomFilter.
func BenchmarkBasicBloomFilterTest(b *testing.B) {
	bf := NewBasicBloomFilter(WithCapacity(1000000), WithFalsePositiveRate(0.01))

	for i := 0; i < 100000; i++ {
		value := []byte(fmt.Sprintf("value-%d", i))
		bf.Add(value)
	}

	b.ResetTimer() // Reset the timer to exclude the setup time
	for n := 0; n < b.N; n++ {
		value := []byte(fmt.Sprintf("value-%d", n))
		bf.Test(value)
	}
}
