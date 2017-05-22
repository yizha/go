package tp

import (
	//	"fmt"
	"testing"
)

func TestTP(t *testing.T) {
	s := []float32{300, 400, 500}
	d := []float32{250, 350, 400, 200}
	c := [][]float32{
		[]float32{3, 1, 7, 4},
		[]float32{2, 6, 5, 9},
		[]float32{8, 3, 3, 2},
	}
	for i := 0; i < 10; i++ {
		//sLen, dLen := len(s), len(d)
		_, _, err := Solve(s, d, c)
		if err != nil {
			t.Error("failed to solve givn TP problem:", err)
			//fmt.Println(err)
			//os.Exit(1)
		} else {
			/*fmt.Println("Solution: ")
			for i := 0; i < sLen; i++ {
				for j := 0; j < dLen; j++ {
					f := flow[i][j]
					if f == 0 {
						continue
					}
					fmt.Printf(" (%v,%v),cost=%v,flow=%v\n", i, j, c[i][j], f)
				}
			}
			fmt.Printf("solution cost=%v\n", cost)
			fmt.Println("")*/
		}
	}
}
