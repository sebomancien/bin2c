package expression

import (
	"fmt"
	"strings"
)

// Operation interface
type Operation interface {
	Evaluate(inputs []float32) (float32, error)
	Simplify() Operation
	String() string
}

// Constant
type Const struct {
	value float32
}

func NewConst(value float32) *Const {
	return &Const{value: value}
}

func (op *Const) Evaluate(inputs []float32) (float32, error) {
	return op.value, nil
}

func (op *Const) Simplify() Operation {
	return op
}

func (op *Const) String() string {
	return fmt.Sprint(op.value)
}

// Variable
type Var struct {
	index uint8
}

func NewVar(index uint8) *Var {
	return &Var{index: index}
}

func (op *Var) Evaluate(inputs []float32) (float32, error) {
	return inputs[op.index], nil
}

func (op *Var) Simplify() Operation {
	return op
}

func (op *Var) String() string {
	return fmt.Sprintf("{%d}", op.index)
}

// Addition
type Add struct {
	terms []Operation
}

func NewAdd(terms ...Operation) *Add {
	return &Add{terms: terms}
}

func (op *Add) Evaluate(inputs []float32) (float32, error) {
	sum := float32(0)
	for _, operation := range op.terms {
		val, err := operation.Evaluate(inputs)
		if err != nil {
			return 0, err
		}
		sum += val
	}
	return sum, nil
}

func (op *Add) Simplify() Operation {
	var sum float32 = 0
	var terms []Operation

	for _, t := range op.terms {
		t = t.Simplify()
		if add, ok := t.(*Add); ok {
			// Simplify all terms and combine sub-additions
			terms = append(terms, add.terms...)
		} else if cst, ok := t.(*Const); ok {
			// Combine all constants together
			sum += cst.value
		} else {
			terms = append(terms, t)
		}
	}

	// Add back the constant, when required
	if sum != 0 {
		cst := NewConst(sum)
		if len(terms) == 0 {
			return cst
		}
		terms = append(terms, cst)
	}

	return NewAdd(terms...)
}

func (op *Add) String() string {
	strs := make([]string, len(op.terms))
	for i, v := range op.terms {
		strs[i] = fmt.Sprint(v)
	}
	return fmt.Sprintf("(%s)", strings.Join(strs, "+"))
}

// Substraction
type Sub struct {
	left  Operation
	right Operation
}

func NewSub(left Operation, right Operation) *Sub {
	return &Sub{left: left, right: right}
}

func (op *Sub) Evaluate(inputs []float32) (float32, error) {
	den, err := op.left.Evaluate(inputs)
	if err != nil {
		return 0, err
	}

	num, err := op.right.Evaluate(inputs)
	if err != nil {
		return 0, err
	}

	return num - den, nil
}

func (op *Sub) Simplify() Operation {
	op.left = op.left.Simplify()
	op.right = op.right.Simplify()

	left, ok := op.left.(*Const)
	if !ok {
		return op
	}
	right, ok := op.right.(*Const)
	if !ok {
		return op
	}

	return NewConst(left.value - right.value)
}

func (op *Sub) String() string {
	return fmt.Sprintf("(%s-%s)", op.left, op.right)
}

// Multiply
type Mul struct {
	terms []Operation
}

func NewMul(terms ...Operation) *Mul {
	return &Mul{terms: terms}
}

func (op *Mul) Evaluate(inputs []float32) (float32, error) {
	sum := float32(1)
	for _, operation := range op.terms {
		val, err := operation.Evaluate(inputs)
		if err != nil {
			return 0, err
		}
		sum *= val
	}
	return sum, nil
}

func (op *Mul) Simplify() Operation {
	var product float32 = 1
	var terms []Operation

	for _, t := range op.terms {
		t = t.Simplify()
		if add, ok := t.(*Mul); ok {
			// Simplify all terms and combine sub-multiplications
			terms = append(terms, add.terms...)
		} else if cst, ok := t.(*Const); ok {
			// Combine all constants together
			product *= cst.value
		} else {
			terms = append(terms, t)
		}
	}

	// Multiply back the constant, when required
	if product != 1 {
		cst := NewConst(product)
		if len(terms) == 0 {
			return cst
		}
		terms = append(terms, cst)
	}

	return NewMul(terms...)
}

func (op *Mul) String() string {
	strs := make([]string, len(op.terms))
	for i, v := range op.terms {
		strs[i] = fmt.Sprint(v)
	}
	return fmt.Sprintf("(%s)", strings.Join(strs, "x"))
}

// Division
type Div struct {
	num Operation
	den Operation
}

func NewDiv(num Operation, den Operation) *Div {
	return &Div{num: num, den: den}
}

func (op *Div) Evaluate(inputs []float32) (float32, error) {
	den, err := op.den.Evaluate(inputs)
	if err != nil {
		return 0, err
	} else if den == 0 {
		return 0, fmt.Errorf("division by zero")
	}

	num, err := op.num.Evaluate(inputs)
	if err != nil {
		return 0, err
	}

	return num / den, nil
}

func (op *Div) Simplify() Operation {
	op.num = op.num.Simplify()
	op.den = op.den.Simplify()

	num, ok := op.num.(*Const)
	if !ok {
		return op
	}
	den, ok := op.den.(*Const)
	if !ok {
		return op
	}

	return NewConst(num.value / den.value)
}

func (op *Div) String() string {
	return fmt.Sprintf("(%s)/(%s)", op.num, op.den)
}
