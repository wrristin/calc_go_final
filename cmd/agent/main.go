package main

import (
	"calc_service/internal/agent"
	"sync"
)

func main() {
	agent.InitGRPCClient()
	defer agent.CloseGRPCClient()

	var wg sync.WaitGroup
	for i := 0; i < agent.ComputingPower; i++ {
		wg.Add(1)
		go agent.Worker(&wg)
	}
	wg.Wait()
}
