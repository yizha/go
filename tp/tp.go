package tp

import (
	"fmt"
)

const (
	// default EPSILON value for checking if a value is zero or not
	EPSILON = float32(1e-6)

	// default INFINITY value for calculating some minimal value
	INFINITY = float32(1e20)

	// default max iterations for optimizing the solution
	MAX_ITER = 100
)

type state struct {

	// 'static' variables
	epsilon, infinity float32
	maxIter           int

	// inputs, could be adjusted if supply/demand is unbalanced
	supply     []float32
	demand     []float32
	costMatrix [][]float32

	// balance flag
	//  -1: supply < demand
	//   0: supply == demand
	//   1: supply > demand
	balanced int

	// supply, demand size after adjustment if inputs are unbalanced
	sLen, dLen int

	// total amount to tranport from producers to consumers
	// including the dummy amount if inputs are unbalanced
	quatity float32

	// u, v matrix
	u, v []float32

	// optimization starting cell
	row, col int

	// When computing u,v
	// 0: init state (not scanned)
	// 1: scan candidate
	// 2: already scanned
	// -------------------
	// When finding the loop
	// 0: row/col is available
	// 1: row/col is occupied
	rowFlags, colFlags []int

	// loop (link-list) head
	loop *cell

	// solution
	cost float32
	flow [][]float32
}

type cell struct {
	// cell location
	row, col int

	// if this cell is the even one in the loop/chain
	even bool

	// minimum flow value of the even cell in the loop/chain
	loopEvenMinFlow float32

	// direction to next cell
	//  true:  horizontal
	//  false: vertical
	horizontal bool

	// pointers to prev/next cell
	prev, next *cell
}

func (c *cell) String() string {
	direction := "V"
	if c.horizontal {
		direction = "H"
	}
	sign := "+"
	if c.even {
		sign = "-"
	}
	return fmt.Sprintf("(%v,%v)/%v/%v", c.row, c.col, direction, sign)
}

func f32Min(a, b float32) float32 {
	if a < b {
		return a
	}
	return b
}

func f32Max(a, b float32) float32 {
	if a > b {
		return a
	}
	return b
}

func createState(s, d []float32, c [][]float32, maxIter int, epsilon, infinity float32) (*state, error) {
	sLen := len(s)
	if sLen < 1 {
		return nil, fmt.Errorf("not enough producers, need at least 1!")
	}
	dLen := len(d)
	if dLen < 1 {
		return nil, fmt.Errorf("not enough consumers, need at least 1! ")
	}
	if sLen != len(c) {
		return nil, fmt.Errorf("producer count doesn't match 1st dimension length of costMatrix!")
	}
	if dLen != len(c[0]) {
		return nil, fmt.Errorf("consumer count doesn't match 2nd dimension length of costMatrix!")
	}

	var sSum, dSum, quatity float32
	for i := 0; i < sLen; i++ {
		if s[i] < epsilon {
			return nil, fmt.Errorf("supply[%v]=%v is too small (<%v)!", i, s[i], epsilon)
		}
		sSum += s[i]
	}
	//if sSum <= epsilon {
	//	return nil, fmt.Errorf("total supply is ")
	//}
	for i := 0; i < dLen; i++ {
		if d[i] < epsilon {
			return nil, fmt.Errorf("demand[%v]=%v is too small (<%v)!", i, d[i], epsilon)
		}
		dSum += d[i]
	}
	diff := float32(0)
	balanced := 0
	if sSum > dSum {
		diff = sSum - dSum
		quatity = sSum
		balanced = 1
	} else if sSum < dSum {
		diff = dSum - sSum
		quatity = dSum
		balanced = -1
	} else {
		quatity = sSum
		balanced = 0
	}
	var supply, demand []float32
	var costMatrix [][]float32
	if diff > epsilon { // unbalanced
		if sSum > dSum {
			// copy supply
			supply = make([]float32, sLen)
			for i := 0; i < sLen; i++ {
				supply[i] = s[i]
			}
			// copy demand and add one more (0)
			demand = make([]float32, dLen+1)
			for i := 0; i < dLen; i++ {
				demand[i] = d[i]
			}
			d[dLen] = diff
			// copy cost matrix, append one column (0s)
			costMatrix = make([][]float32, sLen)
			for i := 0; i < sLen; i++ {
				costMatrix[i] = make([]float32, dLen+1)
				for j := 0; j < dLen; j++ {
					costMatrix[i][j] = c[i][j]
				}
				costMatrix[i][dLen] = 0
			}
			// fix demand size
			dLen += 1
		} else { // sSum < dSum
			// copy supply and add one more (diff)
			supply = make([]float32, sLen+1)
			for i := 0; i < sLen; i++ {
				supply[i] = s[i]
			}
			supply[sLen] = diff
			// copy demand
			demand = make([]float32, dLen)
			for i := 0; i < dLen; i++ {
				demand[i] = d[i]
			}
			// copy cost matrix and append one row (0s)
			costMatrix = make([][]float32, sLen+1)
			for i := 0; i < sLen; i++ {
				costMatrix[i] = make([]float32, dLen)
				for j := 0; j < dLen; j++ {
					costMatrix[i][j] = c[i][j]
				}
			}
			costMatrix[sLen] = make([]float32, dLen)
			// fix supply size
			sLen += 1
		}
	} else { // balanced
		supply = make([]float32, sLen)
		for i := 0; i < sLen; i++ {
			supply[i] = s[i]
		}
		demand = make([]float32, dLen)
		for i := 0; i < dLen; i++ {
			demand[i] = d[i]
		}
		costMatrix = make([][]float32, sLen)
		for i := 0; i < sLen; i++ {
			costMatrix[i] = make([]float32, dLen)
			for j := 0; j < dLen; j++ {
				costMatrix[i][j] = c[i][j]
			}
		}
	}

	// create emd-state struct
	flow := make([][]float32, sLen)
	for i := 0; i < sLen; i++ {
		flow[i] = make([]float32, dLen)
	}
	return &state{
		epsilon:  epsilon,
		infinity: infinity,
		maxIter:  maxIter,

		supply:     supply,
		demand:     demand,
		costMatrix: costMatrix,
		balanced:   balanced,

		sLen:     sLen,
		dLen:     dLen,
		quatity:  quatity,
		u:        make([]float32, sLen),
		v:        make([]float32, dLen),
		row:      -1,
		col:      -1,
		rowFlags: make([]int, sLen),
		colFlags: make([]int, dLen),
		loop:     nil,

		cost: float32(0),
		flow: flow,
	}, nil
}

func (es *state) printSolution() {
	fmt.Println("[Solution]")
	cost := float32(0)
	epsilon := es.epsilon
	for i := 0; i < es.sLen; i++ {
		for j := 0; j < es.dLen; j++ {
			if es.flow[i][j] < epsilon {
				continue
			}
			//fmt.Printf("i=%v,j=%v,supply=%v,demand=%v,transport=%v\n",
			//i, j, es.supply[i], es.demand[j], es.flow[i][j])
			c := es.costMatrix[i][j]
			f := es.flow[i][j]
			fmt.Printf(" (%v,%v),cost=%v,flow=%v\n", i, j, c, f)
			cost += c * f
		}
	}
	fmt.Printf("cost=%v\n", cost)
	fmt.Println("")
}

// find the initial solution with the "Minimal Cost" method
func (es *state) findFeasibleSolution() error {
	//fmt.Println("[Finding feasible solution ...]")

	sLen, dLen := es.sLen, es.dLen
	epsilon, infinity := es.epsilon, es.infinity

	quatity := es.quatity
	//fmt.Printf("quatity=%v\n", quatity)
	flowCnt := 0

	for {
		// the least-cost (selected) row/column index
		si, sj := -1, -1
		var minCost = infinity
		// loop to find the least cost row/column
		for i := 0; i < sLen; i++ {
			// skip row if supply is "0"
			if es.supply[i] <= epsilon {
				continue
			}
			for j := 0; j < dLen; j++ {
				// skip column if demand is "0"
				if es.demand[j] <= epsilon {
					continue
				}
				cost := es.costMatrix[i][j]
				if cost < minCost {
					minCost = cost
					si = i
					sj = j
				} else if cost == minCost {
					// for same cost cell, choose the one which
					// transports more
					sq := f32Min(es.supply[si], es.demand[sj])
					q := f32Min(es.supply[i], es.demand[j])
					if q > sq {
						si = i
						sj = j
					}
				}
			}
		}
		// substract the selected quatity from supply/demand
		//fmt.Printf("got min cost cell at i=%v, j=%v, cost=%v\n", si, sj, es.costMatrix[si][sj])
		s := es.supply[si]
		d := es.demand[sj]
		q := float32(0)
		if d < s {
			q = d
			es.supply[si] = s - q
			es.demand[sj] = 0
		} else {
			q = s
			es.demand[sj] = d - q
			es.supply[si] = 0
		}
		quatity = quatity - q
		// save it to flow matrix
		es.flow[si][sj] = q
		flowCnt += 1

		if quatity <= epsilon {
			break
		}
	}

	//fmt.Println("")

	if sLen+dLen == flowCnt+1 {
		return nil
	} else {
		return fmt.Errorf("[findFeasibleSolution] sLen:%v + dLen:%v - 1 != flowCnt:%v", sLen, dLen, flowCnt)
	}
}

func (es *state) computeUV() error {
	//fmt.Println("[Computing U,V ...]")

	sLen, dLen := es.sLen, es.dLen
	epsilon := es.epsilon

	// reset row/col flags
	for i := 0; i < sLen; i++ {
		es.rowFlags[i] = 0
	}
	for i := 0; i < dLen; i++ {
		es.colFlags[i] = 0
	}
	// computed counts
	uComputedCnt, vComputedCnt := 0, 0
	// set u[0] = 0
	es.u[0] = float32(0)
	es.rowFlags[0] = 1
	uComputedCnt += 1
	//fmt.Println("set u[0]=0")
	//fmt.Println("")

	more2scan := false

	for {
		// scan rows
		more2scan = false
		for row := 0; row < sLen; row++ {
			if es.rowFlags[row] != 1 {
				continue
			}
			//fmt.Printf(">> scanning row #%v ...\n", row)
			for col := 0; col < dLen; col++ {
				if es.flow[row][col] <= epsilon || es.colFlags[col] > 1 {
					//fmt.Printf(">>>> skipping col #%v.\n", col)
					continue
				}
				v := es.costMatrix[row][col] - es.u[row]
				es.v[col] = v
				vComputedCnt += 1
				//fmt.Printf(">>>> computed v[%v]=%v\n", col, v)
				if es.colFlags[col] == 0 {
					es.colFlags[col] = 1
					more2scan = true
					//fmt.Printf(">>>> set colFlags[%v]=1\n", col)
				}
			}
			es.rowFlags[row] = 2
			//fmt.Printf(">> set rowFlags[%v]=2\n", row)
			//fmt.Println("")
		}

		//fmt.Printf("rowFlags=%v, colFlags=%v, uComputedCnt=%v, vComputedCnt=%v\n", es.rowFlags, es.colFlags, uComputedCnt, vComputedCnt)
		//fmt.Println("")
		if !more2scan || (uComputedCnt == sLen && vComputedCnt == dLen) {
			//fmt.Println("break!")
			break
		}

		// scan columns
		more2scan = false
		for col := 0; col < dLen; col++ {
			if es.colFlags[col] != 1 {
				continue
			}
			//fmt.Printf(">> scanning col #%v ...\n", col)
			for row := 0; row < sLen; row++ {
				if es.flow[row][col] <= epsilon || es.rowFlags[row] > 1 {
					//fmt.Printf(">>>> skipping row #%v.\n", row)
					continue
				}
				u := es.costMatrix[row][col] - es.v[row]
				es.u[row] = u
				uComputedCnt += 1
				//fmt.Printf(">>>> computed u[%v]=%v\n", row, u)
				if es.rowFlags[row] == 0 {
					es.rowFlags[row] = 1
					more2scan = true
					//fmt.Printf(">>>> set rowFlags[%v]=1\n", row)
				}
			}
			es.colFlags[col] = 2
			//fmt.Printf(">> set colFlags[%v]=2\n", col)
			//fmt.Println("")
		}

		//fmt.Printf("rowFlags=%v, colFlags=%v, uComputedCnt=%v, vComputedCnt=%v\n", es.rowFlags, es.colFlags, uComputedCnt, vComputedCnt)
		//fmt.Println("")
		if !more2scan || (uComputedCnt == sLen && vComputedCnt == dLen) {
			//fmt.Println("break!")
			break
		}
	}

	//fmt.Printf("computed u=%v\n", es.u)
	//fmt.Printf("computed v=%v\n", es.v)
	//fmt.Println("")

	if uComputedCnt != sLen || vComputedCnt != dLen {
		return fmt.Errorf("[computeUV()] U: %v/%v, V: %v/%v", uComputedCnt, sLen, vComputedCnt, dLen)
	} else {
		return nil
	}
}

func (es *state) isOptimal() bool {
	//fmt.Println("[Checking if current solution is optimal ...]")
	// find the base cell by computing the penalty for all no-flow cell
	sLen, dLen := es.sLen, es.dLen
	epsilon := es.epsilon

	es.row, es.col = -1, -1
	pMax := float32(-1)
	for i := 0; i < sLen; i++ {
		for j := 0; j < dLen; j++ {
			if es.flow[i][j] > epsilon {
				continue
			}
			p := es.u[i] + es.v[j] - es.costMatrix[i][j]
			//fmt.Printf("i=%v,j=%v,p=%v\n", i, j, p)
			if p > epsilon && p > pMax {
				es.row, es.col, pMax = i, j, p
			}
		}
	}
	var optimal bool = (es.row == -1)
	/*if optimal {
		fmt.Println("current solution is optimal!")
	} else {
		fmt.Printf("current solution is NOT optimal! Optimization starting cell: (%v, %v)\n", es.row, es.col)
	}
	fmt.Println("")*/
	return optimal
}

func (es *state) findLoop() error {
	//fmt.Println("[Finding a valid loop ...]")
	sLen, dLen := es.sLen, es.dLen
	epsilon, infinity := es.epsilon, es.infinity
	// reset row/col flags
	for i := 0; i < sLen; i++ {
		es.rowFlags[i] = 0
	}
	for i := 0; i < dLen; i++ {
		es.colFlags[i] = 0
	}
	head := &cell{
		row:             es.row,
		col:             es.col,
		horizontal:      true,
		even:            false,
		loopEvenMinFlow: infinity,
		prev:            nil,
		next:            nil,
	}

	curr := head
	cycle := 0
	step := 0

	for {

		if curr == head {
			cycle += 1
			if cycle == 1 {
				// starting horizontally
				//fmt.Printf("starting from cell at %v\n", head)
			} else if cycle == 2 {
				// try again starting vertically
				curr.horizontal = false
				//fmt.Printf("starting again from cell at %v\n", head)
			} else {
				// fail to found a valid loop
				head = nil
				break
			}
		}

		nexti, nextj := -1, -1
		if curr.horizontal {
			i := curr.row
			next := curr.next
			start := 0
			if next != nil {
				start = next.col + 1
			}
			if start < dLen {
				for j := start; j < dLen; j++ {
					if j == curr.col || es.flow[i][j] <= epsilon || es.colFlags[j] == 1 {
						continue
					}
					nexti, nextj = i, j
					es.rowFlags[i] = 1
					break
				}
			}
		} else { // vertical
			j := curr.col
			next := curr.next
			start := 0
			if next != nil {
				start = next.row + 1
			}
			if start < sLen {
				for i := 0; i < sLen; i++ {
					if i == curr.row || es.flow[i][j] <= epsilon || es.rowFlags[i] == 1 {
						continue
					}
					nexti, nextj = i, j
					es.colFlags[j] = 1
					break
				}
			}
		}

		if nexti >= 0 { // found next cell
			nextEven := !curr.even
			loopEvenMinFlow := curr.loopEvenMinFlow
			// new cell is the even one in the chain,
			// need to calculate the min flow
			if nextEven {
				nextFlow := es.flow[nexti][nextj]
				if nextFlow < loopEvenMinFlow {
					loopEvenMinFlow = nextFlow
				}
			}

			next := &cell{
				row:             nexti,
				col:             nextj,
				horizontal:      !curr.horizontal,
				even:            nextEven,
				loopEvenMinFlow: loopEvenMinFlow,
				prev:            curr,
				next:            nil,
			}
			curr.next = next
			curr = next
			step += 1

			//fmt.Printf("step #%v: found next cell %v\n", step, curr)

			if (head.horizontal && (curr.col == head.col)) || (!head.horizontal && (curr.row == head.row)) {
				// found a valid loop
				break
			}
		} else { // cannot find next cell, need to go back
			if curr != head { // cannot go back if it is already head
				curr = curr.prev
				if curr.horizontal {
					es.rowFlags[curr.row] = 0
				} else {
					es.colFlags[curr.col] = 0
				}
				step += 1

				//fmt.Printf("step #%v: go back to %v\n", step, curr)
			}
		}
	}

	if head != nil {
		// save loop even-cell min flow to head cell for easy access
		head.loopEvenMinFlow = curr.loopEvenMinFlow
		es.loop = head

		/*fmt.Println("found a valid loop:")
		p := head
		for p != nil {
			if p != head {
				fmt.Print(" -> ")
			}
			fmt.Printf("%v", p)
			//fmt.Printf("%v/%v/%v", p, es.flow[p.row][p.col], p.loopEvenMinFlow)
			p = p.next
		}
		fmt.Println("")
		fmt.Printf("loop even-cell min flow: %v\n", head.loopEvenMinFlow)
		fmt.Println("")*/

		return nil
	} else {
		es.loop = nil
		//fmt.Println("")
		return fmt.Errorf("[findLoop()] cannot find a valid loop starting from (%v,%v).", es.row, es.col)
	}
}

func (es *state) applyOptimization() {
	p := es.loop
	q := p.loopEvenMinFlow
	for p != nil {
		row, col := p.row, p.col
		flow := es.flow[row][col]
		if p.even {
			es.flow[row][col] = flow - q
		} else {
			es.flow[row][col] = flow + q
		}
		p = p.next
	}
}

func (es *state) solve() error {
	if err := es.findFeasibleSolution(); err != nil {
		return err
	}
	//es.printSolution()

	maxIter := es.maxIter
	iter := 0
	//optimal := false
	for {
		if err := es.computeUV(); err != nil {
			return err
		}
		if es.isOptimal() {
			//optimal = true
			break
		}
		if err := es.findLoop(); err != nil {
			return err
		}
		es.applyOptimization()
		//es.printSolution()
		iter += 1
		if iter > maxIter {
			break
		}
	}

	//fmt.Printf("is optimal solution: %v\n", optimal)
	//fmt.Printf("iteration ran: %v\n", iter)
	//fmt.Println()

	return nil
}

func (es *state) getSolution() (float32, [][]float32) {
	epsilon, balanced, sLen, dLen := es.epsilon, es.balanced, 0, 0
	if balanced < 0 {
		sLen, dLen = es.sLen-1, es.dLen
	} else if balanced > 0 {
		sLen, dLen = es.sLen, es.dLen-1
	} else {
		sLen, dLen = es.sLen, es.dLen
	}
	cost := float32(0)
	flow := make([][]float32, sLen)
	for i := 0; i < sLen; i++ {
		flow[i] = make([]float32, dLen)
		for j := 0; j < dLen; j++ {
			f := es.flow[i][j]
			if f < epsilon {
				continue
			}
			flow[i][j] = f
			cost += f * es.costMatrix[i][j]
		}
	}
	return cost, flow
}

// Solves a given transportation problem and returns the (optimal) solution.
//
//  supply, demand: positive float32 array/slice.
//  costs: 2-D matrix, row size should match supply length, column
//         size should match demand length.
//
//  opts is for optional args:
//   opts[0]: MAX_ITER, max iterations to run when optimizing the
//            solution, default to 100.
//   opts[1]: EPSILON, used to tell if a float32 value is zero or not,
//            default to 1e-6.
//   opts[2]: INFINITY, used as the ceiling value when find some min
//            value, default to 1e20.
//   if you need to use a non-default EPSILON (opt[1]), you must also
//   set MAX_ITER (opt[0]), similarly to set a different INFINITY
//   (opt[2]) from the default value, you must also set both MAX_ITER
//   (opt[0]) and EPSILON (opt[1]).
//
//  returns
//   cost: the solution cost (float32) or -1 if there is error
//   flow: the flow matrix (float32) or nil if there is error
//    err: nil or an error object if something goes wrong.
func Solve(supply, demand []float32, costs [][]float32, opts ...float32) (cost float32, flow [][]float32, err error) {
	maxIter, epsilon, infinity := MAX_ITER, EPSILON, INFINITY
	optsLen := len(opts)
	if optsLen > 0 {
		maxIter := int(opts[0])
		if maxIter < 1 {
			return float32(-1), nil, fmt.Errorf("invalid maxIter: %v", opts[0])
		}
		if optsLen > 1 {
			epsilon := opts[1]
			if epsilon > float32(1e-3) {
				return float32(-1), nil, fmt.Errorf("Given epsilon is too big (>1e-3): %v", opts[1])
			}
			if optsLen > 2 {
				infinity := opts[2]
				if infinity < float32(1e10) {
					return float32(-1), nil, fmt.Errorf("Given infinity is too small (<1e10): %v", opts[2])
				}
			}
		}
	}

	es, err := createState(supply, demand, costs, maxIter, epsilon, infinity)
	if err != nil {
		return float32(-1), nil, err
	}
	if err := es.solve(); err != nil {
		return float32(-1), nil, err
	}
	//es.printSolution()
	cost, flow = es.getSolution()
	return
}
