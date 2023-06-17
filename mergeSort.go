// 19269749 - Oliver Nagy
// Algorithm 2 MergeSort using Waitgroups
// O(n*log(n))
// For this algorithm I had to change the test harness to 3 million

// These are averages for 3 runs / each time I varied the max number of cores the programm could use.
/*

USING MAX 6 CORES
Sort time ~ 2.4s

USING MAX 5 CORES
Sort time ~ 2.6s

USING MAX 4 CORES
Sort time ~ 2.9s

USING MAX 3 CORES
Sort time ~3.7s


USING MAX 2 CORES
Sort time ~ 4.7s

USING 1 CORE
Sort time ~ 7.84s

*****Visual Representation of Varying MAX Cores********



       7.84s |x
       7.40s |
       7.20s |
       6.80s |
       6.40s |
       6.00s |
       5.60s |
       5.20s |
       4.80s |     x
       4.40s |
       4.00s |
       3.60s |         x
       3.20s |              x
       2.80s |                 x
       2.40s |                      x
			   1---2---3---4---5---6



*/

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

	runtime.GOMAXPROCS(6) //varies max number of cores
	const numRuns = 3
	rand.Seed(420) // adds determinism to Slice generation
	fmt.Println("Sorting 3 million numbers...")
	times := make([]time.Duration, numRuns)
	for i := 0; i < numRuns; i++ {
		slice := generateSlice(3000000) //takes 5-10 seconds to sort for optimised algorithms
		startTime := time.Now()

		//Insert Algorithm Here

		Sort(slice)

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
func Mergesort(s []uint64) {
	len := len(s)

	if len > 1 {
		middle := len / 2 // finds middle of slice

		var wg sync.WaitGroup //var Waitgroup
		wg.Add(2)             // increase wg (waitgroup) by 2 for our two go functions below

		go func() {
			defer wg.Done()       // defers decremenet of Waitgroup counter until the two halfs are sorted properly
			Mergesort(s[:middle]) //sorts divided slice first half
		}()

		go func() {
			defer wg.Done()       // defers decremenet of Waitgroup counter until the two halfs are sorted properly
			Mergesort(s[middle:]) // sorts divided slice second half
		}()

		wg.Wait()        // Wait block exection of merge(s, middle) until counter is 0 or all go funcs finish above
		merge(s, middle) //merges two together
	}
}

func Sort(arr []uint64) []uint64 {
	Mergesort(arr)
	return arr
}

// Function below megres two half back together
func merge(s []uint64, middle int) {
	helper := make([]uint64, len(s)) //temp array
	copy(helper, s)                  //copy s into temp

	helperLeft := 0
	helperRight := middle
	current := 0
	high := len(s) - 1

	for helperLeft <= middle-1 && helperRight <= high {
		if helper[helperLeft] <= helper[helperRight] {
			s[current] = helper[helperLeft]
			helperLeft++
		} else {
			s[current] = helper[helperRight]
			helperRight++
		}
		current++
	}

	for helperLeft <= middle-1 {
		s[current] = helper[helperLeft]
		current++
		helperLeft++
	}
}
