# TPSA

[![Go Report Card](https://goreportcard.com/badge/github.com/d-tsuji/tpsa)](https://goreportcard.com/report/github.com/d-tsuji/tpsa)

This repository that solves the Traveling Salesman Problem (TSP) using **Temperature Parallel Simulated Annealing** (TPSA). TPSA is a concurrent algorithm for the simulated annealing method, which is published in the following paper.

[Temperature Parallel Simulated Annealing Algorithm and Its Evaluation](http://id.nii.ac.jp/1001/00013940/)

The feature of this repository is that it uses channel and goroutine in Go for concurrency (without using OpenMP etc...), so you can understand the greatness of concurrency mechanism in Go.

This content was presented at the [Go Conference 2019 Autumn](https://gocon.jp/) in Japan.

![tpsa](https://user-images.githubusercontent.com/24369487/81168922-e4a2cc80-8fd2-11ea-9c4d-1ab99b36e361.gif)

The above is a visualization of the process of solving the TSP using TPSA.

## Installation
```
$ git clone https://github.com/d-tsuji/tpsa.git
$ go run cmd/tpsa/main.go
```

## Usage

```go
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
```

### Input

Input data refers to a file in the data directory. Stores some data in TSPLIB.

http://elib.zib.de/pub/mp-testdata/tsp/tsplib/tsp/index.html

### Output

The result of the solution by TPSA and the result of the exact solution by TSPLIB are output as follows.

```
Data(krod100.tsp)
TPSA solution  : 21294.290821490347
Exact solution : 21294.290821490355
```

*Note: The solver viewer is not included in this repository.*
