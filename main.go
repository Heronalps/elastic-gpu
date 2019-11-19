package main

import (
	"fmt"
	"log"
	"time"

	"github.com/heronalps/elastic-gpu/queryprom"
	"github.com/heronalps/elastic-gpu/request"
)

func main() {
	// GPU utilization is in percentage
	upperBound := 65.0
	lowerBound := 25.0
	upperBoundCount := 0
	lowerBoundCount := 0
	zeroCount := 0
	var startGPU int64 = 1

	for {
		// Keep querying the GPU utilization
		fmt.Println("========")
		t := time.Now()
		fmt.Println(t.Format("2006-01-02 15:04:05"))

		gpuUtil := query("racelab", "avg_over_time(namespace_gpu_utilization[1m])")
		fmt.Printf("GPU utilization is %f \n", gpuUtil)
		gpuNum := request.QueryGPU("racelab", "image-clf-train")
		fmt.Printf("Current number of GPU is %v \n", gpuNum)

		if gpuUtil == 0.0 {
			fmt.Printf("No valid GPU utilization data reading \n")
			zeroCount++
			upperBoundCount = 0
			lowerBoundCount = 0
			// When GPU utilization is greater than upper bound, Add one GPU
		} else if gpuUtil > upperBound && gpuNum < 6 {
			upperBoundCount++

			// When GPU utilization is less than lower bound, dismiss one GPU
			// Keep one GPU at least
		} else if gpuUtil < lowerBound && gpuNum > startGPU {
			lowerBoundCount++

		}
		// 5 * 6 = 30 seconds being above upper bound leads to scaling up
		if upperBoundCount >= 6 {
			fmt.Printf("Adding additional GPU... \n")
			request.Update("racelab", "image-clf-train", 1)

			upperBoundCount = 0
			fmt.Println("Waiting for deployment to take effect...")
			time.Sleep(60 * time.Second)
		}
		// 5 secs * 36 = 3 minute being above lower bound leads to scaling down
		if lowerBoundCount >= 36 {
			fmt.Printf("Remove extra GPU... \n")
			request.Update("racelab", "image-clf-train", -1)

			lowerBoundCount = 0
			fmt.Println("Waiting for deployment to take effect...")
			time.Sleep(60 * time.Second)
		}
		// Idling for ten minutes leads to falling back to starting GPU number
		if zeroCount >= 120 && gpuNum > startGPU {
			fmt.Printf("Fall back to %v GPU ... \n", startGPU)
			request.Set("racelab", "image-clf-train", startGPU)
			fmt.Println("Waiting for deployment to take effect...")
			time.Sleep(60 * time.Second)
		}
		time.Sleep(5 * time.Second)
	}
}

func query(namespace string, queryStr string) float64 {
	value, err := queryprom.Query(namespace, queryStr)
	if err != nil {
		log.Printf("%v", err)
	}
	return value
}
