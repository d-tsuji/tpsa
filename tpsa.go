package tpsa

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type TPSA struct {
	TPSAConfig
	Temperatures []float64
	Tours        [][]int
	Matrix       Matrix
	Size         int
}

type TPSAConfig struct {
	MinTemp      float64
	MaxTemp      float64
	Thread       int
	Period       int
	MaxIteration int
	DataFileName string
}

type Matrix [][]float64

func NewTPSA(tpsaConfig TPSAConfig) *TPSA {
	return &TPSA{
		TPSAConfig: tpsaConfig,
	}
}

func (t *TPSA) initialize() error {
	if err := t.loadData(); err != nil {
		return err
	}
	if err := t.makeInitialData(); err != nil {
		return err
	}
	return nil
}

func (t *TPSA) Solve() error {
	if err := t.initialize(); err != nil {
		return err
	}

	iteration := 0
	for iteration < t.MaxIteration {
		wg := &sync.WaitGroup{}
		for i := 0; i < t.Thread; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				t.sa(i)
			}(i)
		}
		wg.Wait()

		t.exchangeSolutions(iteration)
		iteration++
	}

	t.printSolution()

	return nil
}

func (t *TPSA) sa(threadNumber int) {
	cnt := 0
	n := t.Size
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	for cnt < t.Period {
		for i := 0; i < n-2; i++ {
			for j := i + 2; j < n; j++ {
				current := t.getCost(&t.Tours[threadNumber], i, i+1) + t.getCost(&t.Tours[threadNumber], j, j+1)
				next := t.getCost(&t.Tours[threadNumber], i, j) + t.getCost(&t.Tours[threadNumber], i+1, j+1)
				p := random.Float64()
				q := t.exchange(next-current, t.Temperatures[threadNumber])
				force := p <= q
				if next < current || force {
					for k := 0; k < (j-i)/2; k++ {
						t.flip(threadNumber, i+1+k, j-k)
					}
				}
			}
		}
		cnt++
	}
}

func (t *TPSA) getCost(tour *[]int, n int, m int) float64 {
	return t.Matrix[(*tour)[n%t.Size]][(*tour)[m%t.Size]]
}

func (t *TPSA) getTotalCost(tour *[]int) float64 {
	ret := 0.0
	for i := 0; i < len(*tour); i++ {
		ret += t.getCost(tour, i, i+1)
	}
	return ret
}

func (t *TPSA) flip(threadNumber, i, j int) {
	t.Tours[threadNumber][i], t.Tours[threadNumber][j] = t.Tours[threadNumber][j], t.Tours[threadNumber][i]
}

func (t *TPSA) exchange(cost, temp float64) float64 {
	return math.Exp(-cost / temp)
}

func (t *TPSA) exchangeSolution(cur, next int) {
	deltaTemperature := t.Temperatures[next] - t.Temperatures[cur]
	deltaValue := t.getTotalCost(&(t.Tours)[next]) - t.getTotalCost(&(t.Tours)[cur])

	p := rand.Float64()
	q := math.Exp(-deltaTemperature * deltaValue / (t.Temperatures[next] * t.Temperatures[cur]))

	if deltaTemperature*deltaValue < 0 || p <= q {
		t.Tours[next], t.Tours[cur] = t.Tours[cur], t.Tours[next]
	}
}

func (t *TPSA) exchangeSolutions(count int) {
	if count%2 == 0 {
		for i := 0; i < t.Thread-1; i += 2 {
			t.exchangeSolution(i, i+1)
		}
	} else {
		for i := 1; i < t.Thread-1; i += 2 {
			t.exchangeSolution(i, i+1)
		}
	}
}

func (t *TPSA) loadData() error {
	tspData, err := os.Open("./data/" + t.DataFileName)
	if err != nil {
		return err
	}
	defer tspData.Close()

	xys := make([]Point, 0)

	rowCount := 0
	scanner := bufio.NewScanner(tspData)
	for scanner.Scan() {
		line := scanner.Text()
		nums := strings.Split(line, "\t")
		x, _ := strconv.ParseFloat(nums[0], 64)
		y, _ := strconv.ParseFloat(nums[1], 64)
		xys = append(xys, Point{x, y})
		rowCount++
	}
	t.Size = rowCount

	t.Matrix = make(Matrix, t.Size)
	for i := 0; i < t.Size; i++ {
		t.Matrix[i] = make([]float64, t.Size)
	}
	for i := 0; i < t.Size; i++ {
		for j := i + 1; j < t.Size; j++ {
			cost := dist(xys[i], xys[j])
			t.Matrix[i][j] = cost
			t.Matrix[j][i] = cost
		}
	}

	return nil
}

func (t *TPSA) makeInitialData() error {
	rand.Seed(time.Now().UnixNano())

	t.Temperatures = make([]float64, t.Thread)
	t.Tours = make([][]int, t.Thread)
	interval := (t.MaxTemp - t.MinTemp) / float64(t.Thread-1)
	for i := 0; i < t.Thread; i++ {
		t.Temperatures[t.Thread-1-i] = interval*float64(i) + t.MinTemp
	}

	tour := make([]int, 0)
	for i := 0; i < t.Size; i++ {
		tour = append(tour, i)
	}
	rand.Shuffle(len(tour), func(i, j int) { tour[i], tour[j] = tour[j], tour[i] })
	for i := 0; i < t.Thread; i++ {
		tmpArray := make([]int, t.Size)
		copy(tmpArray, tour)
		t.Tours[i] = tmpArray
	}

	return nil
}

func (t *TPSA) printSolution() {
	fmt.Printf("Data(%v)\n", t.DataFileName)
	fmt.Printf("TPSA solution  : %v\n", t.getTotalCost(&(t.Tours[t.Thread-1])))
	tspData, err := os.Open("./data/ans/" + strings.ReplaceAll(t.DataFileName, ".tsp", "") + ".opt.tour")
	if err != nil {
		panic(err)
	}
	defer tspData.Close()

	ansTour := make([]int, 0)
	scanner := bufio.NewScanner(tspData)
	for scanner.Scan() {
		num := scanner.Text()
		x, _ := strconv.Atoi(num)
		ansTour = append(ansTour, x-1)
	}
	fmt.Printf("Exact solution : %v\n", t.getTotalCost(&ansTour))
}
