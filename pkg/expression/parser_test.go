package expression

import (
	"testing"
)

func TestParser(t *testing.T) {
	values := []float32{4.2}
	exp := "15+(222-76)/7.5+7*{0}"
	expected := expression(values)

	op, err := Parse(exp)
	if err != nil {
		t.Error(err)
		return
	}

	calculated, err := op.Evaluate(values)
	if err != nil {
		t.Error(err)
	}
	if calculated != expected {
		t.Error("mismatch between calculated and expected value")
	}
}

func expression(values []float32) float32 {
	return 15 + (222-76)/7.5 + 7*values[0]
}
