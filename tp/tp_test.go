package tp

import (
	"fmt"
	"testing"
)

type TestProblem struct {
	id             string
	supply, demand []float64
	costs          [][]float64
}

var testData = []*TestProblem{

	&TestProblem{
		id:     "balanced (supply == demand)",
		supply: []float64{300, 400, 500},
		demand: []float64{250, 350, 400, 200},
		costs: [][]float64{
			[]float64{3, 1, 7, 4},
			[]float64{2, 6, 5, 9},
			[]float64{8, 3, 3, 2},
		},
	},

	&TestProblem{
		id:     "unbalanced (supply > demand)",
		supply: []float64{300, 400, 570},
		demand: []float64{250, 350, 400, 200},
		costs: [][]float64{
			[]float64{3, 1, 7, 4},
			[]float64{2, 6, 5, 9},
			[]float64{8, 3, 3, 2},
		},
	},

	&TestProblem{
		id:     "unbalanced (supply < demand)",
		supply: []float64{300, 400, 500},
		demand: []float64{250, 350, 440, 280},
		costs: [][]float64{
			[]float64{3, 1, 7, 4},
			[]float64{2, 6, 5, 9},
			[]float64{8, 3, 3, 2},
		},
	},

	&TestProblem{
		id:     "balanced (degeneracy)",
		supply: []float64{300, 400, 500, 200},
		demand: []float64{300, 400, 500, 200},
		costs: [][]float64{
			[]float64{0, 2, 8, 4},
			[]float64{2, 0, 5, 9},
			[]float64{8, 5, 0, 3},
			[]float64{4, 9, 3, 0},
		},
	},

	&TestProblem{
		id:     "misc-1",
		supply: []float64{45, 90, 95, 75, 105},
		demand: []float64{120, 80, 50, 75, 85},
		costs: [][]float64{
			[]float64{6, 6, 9, 4, 10},
			[]float64{3, 2, 7, 5, 12},
			[]float64{8, 7, 5, 6, 4},
			[]float64{11, 12, 9, 5, 2},
			[]float64{4, 3, 4, 5, 11},
		},
	},

	&TestProblem{
		id:     "misc-2",
		supply: []float64{35, 50, 40},
		demand: []float64{45, 20, 30, 30},
		costs: [][]float64{
			[]float64{8, 6, 10, 9},
			[]float64{9, 12, 13, 7},
			[]float64{14, 9, 16, 5},
		},
	},
}

func printProblemAndSolution(p *TestProblem, solutionCost float64, flow [][]float64) {
	sLen, dLen := len(p.supply), len(p.demand)
	fmt.Printf("Problem [%v]\n", p.id)
	fmt.Printf(" Supply: %v\n", p.supply)
	fmt.Printf(" Demand: %v\n", p.demand)
	fmt.Printf(" Cost Matrix: \n")
	for i := 0; i < sLen; i++ {
		fmt.Printf("  [")
		for j := 0; j < dLen; j++ {
			if j > 0 {
				fmt.Print(", ")
			}
			fmt.Printf("%v", p.costs[i][j])
		}
		fmt.Println("]")
	}
	fmt.Println(" Solution Flow: ")
	for i := 0; i < sLen; i++ {
		for j := 0; j < dLen; j++ {
			f := flow[i][j]
			if f == 0 {
				continue
			}
			fmt.Printf("  (%v,%v),flow=%v\n", i, j, f)
		}
	}
	fmt.Printf(" Solution Cost=%v\n", solutionCost)
	fmt.Println("")
}

func TestTP(t *testing.T) {
	max := len(testData)
	for i := 0; i < max; i++ {
		tp := testData[i]
		p, err := CreateProblem(tp.supply, tp.demand, tp.costs)
		if err != nil {
			t.Error(fmt.Sprintf("failed to create the problem %v", tp.id), err)
		} else {
			err = p.Solve()
			if err != nil {
				t.Error(fmt.Sprintf("failed to solve the problem %v", tp.id), err)
			} else {
				cost, flow := p.GetSolution()
				printProblemAndSolution(tp, cost, flow)
			}
		}
	}
}
