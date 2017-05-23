package tp

import (
	"fmt"
	"testing"
)

type TestProblem struct {
	id             string
	supply, demand []float32
	costs          [][]float32
}

var testData = []*TestProblem{

	&TestProblem{
		id:     "balanced (supply == demand)",
		supply: []float32{300, 400, 500},
		demand: []float32{250, 350, 400, 200},
		costs: [][]float32{
			[]float32{3, 1, 7, 4},
			[]float32{2, 6, 5, 9},
			[]float32{8, 3, 3, 2},
		},
	},

	&TestProblem{
		id:     "unbalanced (supply > demand)",
		supply: []float32{300, 400, 570},
		demand: []float32{250, 350, 400, 200},
		costs: [][]float32{
			[]float32{3, 1, 7, 4},
			[]float32{2, 6, 5, 9},
			[]float32{8, 3, 3, 2},
		},
	},

	&TestProblem{
		id:     "unbalanced (supply < demand)",
		supply: []float32{300, 400, 500},
		demand: []float32{250, 350, 440, 280},
		costs: [][]float32{
			[]float32{3, 1, 7, 4},
			[]float32{2, 6, 5, 9},
			[]float32{8, 3, 3, 2},
		},
	},

	&TestProblem{
		id:     "balanced (degeneracy)",
		supply: []float32{300, 400, 500, 200},
		demand: []float32{300, 400, 500, 200},
		costs: [][]float32{
			[]float32{0, 2, 8, 4},
			[]float32{2, 0, 5, 9},
			[]float32{8, 5, 0, 3},
			[]float32{4, 9, 3, 0},
		},
	},

	&TestProblem{
		id:     "misc-1",
		supply: []float32{45, 90, 95, 75, 105},
		demand: []float32{120, 80, 50, 75, 85},
		costs: [][]float32{
			[]float32{6, 6, 9, 4, 10},
			[]float32{3, 2, 7, 5, 12},
			[]float32{8, 7, 5, 6, 4},
			[]float32{11, 12, 9, 5, 2},
			[]float32{4, 3, 4, 5, 11},
		},
	},

	&TestProblem{
		id:     "misc-2",
		supply: []float32{35, 50, 40},
		demand: []float32{45, 20, 30, 30},
		costs: [][]float32{
			[]float32{8, 6, 10, 9},
			[]float32{9, 12, 13, 7},
			[]float32{14, 9, 16, 5},
		},
	},
}

func printProblemAndSolution(p *TestProblem, solutionCost float32, flow [][]float32) {
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
