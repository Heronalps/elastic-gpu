package main

import (
	"fmt"
	"log"

	"github.com/heronalps/elastic-gpu/queryprom"
	"github.com/heronalps/elastic-gpu/request"
)

func main() {
	request.Request()
	// request.Update("racelab", "image-clf-train", 1)

}

func query() {
	if value, err := queryprom.Query("racelab", "avg_over_time(namespace_gpu_utilization[1m])"); err != nil {
		log.Printf("%v", err)
	} else {
		fmt.Println(value)
	}
}
