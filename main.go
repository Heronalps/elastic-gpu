package main

import (
	"fmt"
	"log"

	"github.com/heronalps/elastic-gpu/querygpu"
)

func main() {
	// request.Request()
	// request.Update("racelab", "image-clf-train", -1)
	if value, err := querygpu.Query("racelab", "avg_over_time(namespace_gpu_utilization[1m])"); err != nil {
		log.Printf("%v", err)
	} else {
		fmt.Println(value)
	}
}
