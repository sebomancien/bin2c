package zhe

import (
	"log"
	"slices"
	"sync"

	"github.com/sebomancien/tools/pkg/expression"
)

type variable struct {
	name   string
	values []float32
}

type constraint struct {
	name   string
	exp    expression.Operation
	target float32
	min    float32
	max    float32
	weight float32
}

type solution struct {
	values []float32
	score  float32
}

type Solver struct {
	variables   []variable
	constraints []constraint
	solutions   []solution
}

func NewSolver() *Solver {
	return &Solver{}
}

func (s *Solver) Solve(maxSolution int) {
	if len(s.variables) == 0 || len(s.constraints) == 0 {
		log.Fatal("No variables or constraints")
	}

	s.solutions = nil

	// Create a channel to receive the solutions
	channel := make(chan solution)
	defer close(channel)

	// Start a goroutine for each number
	var wg sync.WaitGroup
	for _, value := range s.variables[0].values {
		wg.Add(1)
		go s.solveRoutine(&wg, channel, value)
	}

	// Start a goroutine to close the progress channel when all tasks are done
	go func() {
		for solution := range channel {
			s.insertSolution(solution, maxSolution)
		}
	}()

	wg.Wait()
}

func (s *Solver) solveRoutine(wg *sync.WaitGroup, channel chan<- solution, values ...float32) {
	defer wg.Done()

	inputs := make([]float32, len(s.variables))
	outputs := make([]float32, len(s.constraints))
	copy(inputs, values)

	s.solve(inputs, outputs, len(values), channel)
}

func (s *Solver) solve(inputs []float32, outputs []float32, depth int, channel chan<- solution) {
	if depth == len(inputs)-1 {
		for _, value := range s.variables[depth].values {
			inputs[depth] = value
			ok, score := s.evaluate(inputs, outputs)
			if ok {
				values := make([]float32, len(inputs))
				copy(values, inputs)
				channel <- solution{
					values: values,
					score:  score,
				}
			}
		}
	} else {
		for _, value := range s.variables[depth].values {
			inputs[depth] = value
			s.solve(inputs, outputs, depth+1, channel)
		}
	}
}

func (s *Solver) evaluate(inputs []float32, outputs []float32) (bool, float32) {
	// Evaluate all constraints and check their validity
	var err error
	for i, c := range s.constraints {
		outputs[i], err = c.exp.Evaluate(inputs)
		if err != nil {
			log.Fatal(err)
		}
		if outputs[i] < c.min || outputs[i] > c.max {
			return false, 0
		}
	}

	// Compute the solution score
	var score float32 = 0
	for i, c := range s.constraints {
		var diff float32
		if c.target > outputs[i] {
			diff = c.target - outputs[i]
		} else {
			diff = outputs[i] - c.target
		}
		if c.target != 0 {
			diff /= c.target
		}
		score += diff * c.weight
	}

	return true, score
}

func (s *Solver) insertSolution(solution solution, maxSolution int) {
	// Find the index of the solution, starting with the first one (worst one)
	i := 0
	for ; i < len(s.solutions); i++ {
		if solution.score >= s.solutions[i].score {
			break
		}
	}

	if len(s.solutions) >= maxSolution && i == 0 {
		// If the solution buffer is full and this one is the worst, nothing to do
		return
	}

	// Insert the new solution
	s.solutions = slices.Insert(s.solutions, i, solution)

	// Trim the number of solution by removing the first one (worst one)
	if len(s.solutions) > maxSolution {
		s.solutions = slices.Delete(s.solutions, 0, 1)
	}
}

func (s *Solver) AddVariable(name string, values []float32) {
	s.variables = append(s.variables, variable{
		name:   name,
		values: values,
	})
}

func (s *Solver) AddConstraint(name string, formula string, target float32, min float32, max float32) error {
	exp, err := expression.Parse(formula)
	if err != nil {
		return err
	}

	s.constraints = append(s.constraints, constraint{
		name:   name,
		exp:    exp,
		target: target,
		min:    min,
		max:    max,
		weight: 1,
	})

	return nil
}
