// 19269749 - Oliver Nagy
// Algorithm 1 QuickSort using Waitgroups and Channels
//O(n*log(n))
// Tested on 6 core processor

// Below are the results for Concurrent QuickSort using the Testharness with 3 milion numbers
// These are averages for 3 runs / each time I varied the max number of cores the programm could use.
/*

USING MAX 6 CORES
Sort time ~ 1.37s

USING MAX 5 CORES
Sort time ~ 1.43s

USING MAX 4 CORES
Sort time ~1.64s

USING MAX 3 CORES
Sort time ~ 2.0s


USING MAX 2 CORES
Sort time ~ 2.31s

USING 1 CORE
Sort time ~ 3.92s

*****Visual Representation of Varying MAX Cores********




       3.90s|x

       3.70s|

       3.40s|

       3.10s|

       2.80s|

       2.50s|
	             x
       2.20s|
	                x
       1.90s|

       1.60s|             x
                              x
       1.30s|                    x
			 1---2---3---4---5---6




*/

// Since the divide and conquer nature of quicksort means the recursive calls to Qsort don’t touch the same memory
// Therefore QuickSort would be suitable for concurrency
package main

import (
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"sort"
	"sync"
	"time"
)

// Main Function and Test Harness
func main() {

	runtime.GOMAXPROCS(6) // Limit Max cores

	const numRuns = 3
	rand.Seed(420) // adds determinism to Slice generation
	fmt.Println("Sorting 3 million numbers...")
	times := make([]time.Duration, numRuns)
	for i := 0; i < numRuns; i++ {
		slice := generateSlice(3000000) //takes 5-10 seconds to sort for optimised algorithms
		startTime := time.Now()

		//Insert Algorithm Here

		QuickSort(slice)

		time := time.Since(startTime)

		if sort.SliceIsSorted(slice, func(i, j int) bool { return slice[i] < slice[j] }) {
			fmt.Println("Sorted, Algorithm functional")
			times[i] = time
			fmt.Println("Run", i+1, "time:", times[i])
		} else {
			fmt.Println("Not sorted, Algorithm not functional")
			os.Exit(1) // terminate test if algo fails
		}
	}

	totalTime := time.Duration(0)
	for _, v := range times {
		totalTime += v
	}
	fmt.Println("Average time:", totalTime/numRuns)
}

// Generates Slice of size int for test Harness
func generateSlice(size int) []uint64 {
	slice := make([]uint64, size)
	for i := 0; i < size; i++ {
		slice[i] = rand.Uint64()
	}
	return slice
}

// parition function
func partition(arr []uint64, low, high int) int {
	index := low - 1
	pivotElement := arr[high]
	for i := low; i < high; i++ {
		if arr[i] <= pivotElement {
			index += 1
			arr[index], arr[i] = arr[i], arr[index]
		}
	}
	arr[index+1], arr[high] = arr[high], arr[index+1]
	return index + 1
}

// QuickSortRange Sorts the specified range within the array
func QuickSortRange(arr []uint64, low, high int) {
	if len(arr) <= 1 {
		return
	}

	if low < high {

		var wg sync.WaitGroup //var Waitgroup
		wg.Add(2)

		pivot := partition(arr, low, high)

		go func() {
			defer wg.Done() // defers decremenet of Waitgroup counter until the two halfs are sorted properly
			QuickSortRange(arr, low, pivot-1)
		}()
		go func() {
			defer wg.Done()
			QuickSortRange(arr, pivot+1, high)
		}()
		wg.Wait()
	}

}

// QuickSort Sorts the entire array
func QuickSort(arr []uint64) []uint64 {
	doneChannel := make(chan []uint64) //creates a channel
	go func() {
		QuickSortRange(arr, 0, len(arr)-1)
		doneChannel <- arr //sends sorted array to channel

	}()

	return <-doneChannel // returns channel
}
