package main

import (
	"calc_service/agent/internal/worker"
	"log"
	"os"
	"strconv"
	"sync"
)

func main() {
	// Получаем количество горутин из переменной окружения
	computingPower := 10 // Значение по умолчанию
	if cp, err := strconv.Atoi(os.Getenv("COMPUTING_POWER")); err == nil {
		computingPower = cp
	}

	var wg sync.WaitGroup
	for i := 0; i < computingPower; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			worker.Start()
		}()
	}
	log.Println("Agent started with", computingPower, "workers")
	wg.Wait() // Ожидаем завершения всех горутин
}
