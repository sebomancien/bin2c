package zhe

import (
	"fmt"
	"log"
	"os"
	"slices"
	"testing"
	"text/tabwriter"
	"time"
)

var resistors = []float32{
	1.0e0, 1.1e0, 1.2e0, 1.3e0, 1.5e0, 1.6e0, 1.8e0, 2.0e0, 2.2e0, 2.4e0, 2.7e0, 3.0e0, 3.3e0, 3.6e0, 3.9e0, 4.3e0, 4.7e0, 5.1e0, 5.6e0, 6.2e0, 6.8e0, 7.5e0, 8.2e0, 9.1e0,
	1.0e1, 1.1e1, 1.2e1, 1.3e1, 1.5e1, 1.6e1, 1.8e1, 2.0e1, 2.2e1, 2.4e1, 2.7e1, 3.0e1, 3.3e1, 3.6e1, 3.9e1, 4.3e1, 4.7e1, 5.1e1, 5.6e1, 6.2e1, 6.8e1, 7.5e1, 8.2e1, 9.1e1,
	1.0e2, 1.1e2, 1.2e2, 1.3e2, 1.5e2, 1.6e2, 1.8e2, 2.0e2, 2.2e2, 2.4e2, 2.7e2, 3.0e2, 3.3e2, 3.6e2, 3.9e2, 4.3e2, 4.7e2, 5.1e2, 5.6e2, 6.2e2, 6.8e2, 7.5e2, 8.2e2, 9.1e2,
	1.0e3, 1.1e3, 1.2e3, 1.3e3, 1.5e3, 1.6e3, 1.8e3, 2.0e3, 2.2e3, 2.4e3, 2.7e3, 3.0e3, 3.3e3, 3.6e3, 3.9e3, 4.3e3, 4.7e3, 5.1e3, 5.6e3, 6.2e3, 6.8e3, 7.5e3, 8.2e3, 9.1e3,
	1.0e4, 1.1e4, 1.2e4, 1.3e4, 1.5e4, 1.6e4, 1.8e4, 2.0e4, 2.2e4, 2.4e4, 2.7e4, 3.0e4, 3.3e4, 3.6e4, 3.9e4, 4.3e4, 4.7e4, 5.1e4, 5.6e4, 6.2e4, 6.8e4, 7.5e4, 8.2e4, 9.1e4,
	1.0e5, 1.1e5, 1.2e5, 1.3e5, 1.5e5, 1.6e5, 1.8e5, 2.0e5, 2.2e5, 2.4e5, 2.7e5, 3.0e5, 3.3e5, 3.6e5, 3.9e5, 4.3e5, 4.7e5, 5.1e5, 5.6e5, 6.2e5, 6.8e5, 7.5e5, 8.2e5, 9.1e5,
	1.0e6, 1.1e6, 1.2e6, 1.3e6, 1.5e6, 1.6e6, 1.8e6, 2.0e6, 2.2e6, 2.4e6, 2.7e6, 3.0e6, 3.3e6, 3.6e6, 3.9e6, 4.3e6, 4.7e6, 5.1e6, 5.6e6, 6.2e6, 6.8e6, 7.5e6, 8.2e6, 9.1e6,
}

var capacitors = []float32{
	1.0e-12, 1.2e-12, 1.5e-12, 1.8e-12, 2.2e-12, 2.7e-12, 3.3e-12, 3.9e-12, 4.7e-12, 5.6e-12, 6.8e-12, 8.2e-12,
	1.0e-11, 1.2e-11, 1.5e-11, 1.8e-11, 2.2e-11, 2.7e-11, 3.3e-11, 3.9e-11, 4.7e-11, 5.6e-11, 6.8e-11, 8.2e-11,
	1.0e-10, 1.2e-10, 1.5e-10, 1.8e-10, 2.2e-10, 2.7e-10, 3.3e-10, 3.9e-10, 4.7e-10, 5.6e-10, 6.8e-10, 8.2e-10,
	1.0e-09, 1.2e-09, 1.5e-09, 1.8e-09, 2.2e-09, 2.7e-09, 3.3e-09, 3.9e-09, 4.7e-09, 5.6e-09, 6.8e-09, 8.2e-09,
	1.0e-08, 1.2e-08, 1.5e-08, 1.8e-08, 2.2e-08, 2.7e-08, 3.3e-08, 3.9e-08, 4.7e-08, 5.6e-08, 6.8e-08, 8.2e-08,
	1.0e-07, 1.2e-07, 1.5e-07, 1.8e-07, 2.2e-07, 2.7e-07, 3.3e-07, 3.9e-07, 4.7e-07, 5.6e-07, 6.8e-07, 8.2e-07,
	1.0e-06, 1.2e-06, 1.5e-06, 1.8e-06, 2.2e-06, 2.7e-06, 3.3e-06, 3.9e-06, 4.7e-06, 5.6e-06, 6.8e-06, 8.2e-06,
}

func TestSolver(t *testing.T) {
	var err error
	s := NewSolver()
	s.AddVariable("R1", resistors)
	s.AddVariable("R2", resistors)
	s.AddVariable("C1", capacitors)

	s.AddVariable("R3", resistors)
	s.AddVariable("R4", resistors[:10])

	err = s.AddConstraint("Divider", "3.3*{1}/({0}+{1})", 1, 0.99, 1.01)
	if err != nil {
		log.Fatal(err)
	}
	err = s.AddConstraint("Frequency", "1/((({0}*{1})/({0}+{1}))*{2})", 1000, 900, 1100)
	if err != nil {
		log.Fatal(err)
	}

	start := time.Now()
	s.Solve(10)
	elapsed := time.Since(start)

	s.printSolutions()

	var total uint64 = 1
	for _, v := range s.variables {
		total *= uint64(len(v.values))
	}
	fmt.Printf("Elapsed time        : %s\n", elapsed)
	fmt.Printf("Total possibilities : %d\n", total)
	fmt.Printf("Processing speed    : %.2e possibilities/s\n", float64(total)/elapsed.Seconds())
	fmt.Printf("Solutions found     : %d solutions\n", s.nbSolution)
	fmt.Printf("Solutions ratio     : %.2e\n", float64(s.nbSolution)/float64(total))
}

func (s *Solver) printSolutions() {
	// Create a new tab writer
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)

	// Write header
	fmt.Fprintf(w, "#\tR1\tR2\tC1\tError\n")

	// Write rows
	i := 1
	for _, solution := range slices.Backward(s.solutions) {
		fmt.Fprintf(w, "%d\t%g\t%g\t%g\t%g\n", i, solution.values[0], solution.values[1], solution.values[2], solution.score)
		i++
	}

	// Flush the writer to output the content
	w.Flush()
}
