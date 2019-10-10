package main

import (
	"log"

	"github.com/d-tsuji/tpsa"
)

func main() {
	t := tpsa.NewTPSA(tpsa.TPSAConfig{
		MinTemp:      0.0,
		MaxTemp:      100.0,
		Thread:       16,
		Period:       32,
		MaxIteration: 100,
		DataFileName: "krod100.tsp",
	})

	if err := t.Solve(); err != nil {
		log.Fatal(err)
	}
}
