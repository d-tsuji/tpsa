# TPSA

Temperature Parallel Simulated Annealing

## Installation
```
$ go get github.com/d-tsuji/tpsa
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

Input data refers to a file in the data directory. Stores some data in [TSPLIB](http://elib.zib.de/pub/mp-testdata/tsp/tsplib/tsp/index.html).

### Output

The result of the solution by TPSA and the result of the exact solution by TSPLIB are output as follows.

```
Data(krod100.tsp)
TPSA solution  : 21294.290821490347
Exact solution : 21294.290821490355
```
