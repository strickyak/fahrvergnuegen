package fahrvergnuegen

import (
	"fmt"
	"log"
	"math"
)

type intop func(a, b int) int
type flop func(a, b float64) float64
type strop func(a, b string) string

func (o *Terp) BinaryOp(_i_ intop, _f_ flop, _s_ strop) {
	b := o.Data.Pop()
	a := o.Data.Pop()

	switch a_ := a.(type) {
	case int:
		switch b_ := b.(type) {
		case int:
			o.Data.Push(_i_(a_, b_))
			return
		case float64:
			o.Data.Push(_f_(float64(a_), b_))
			return
		}
	case float64:
		switch b_ := b.(type) {
		case int:
			o.Data.Push(_f_(a_, float64(b_)))
			return
		case float64:
			o.Data.Push(_f_(a_, b_))
			return
		}
	case string:
		switch b_ := b.(type) {
		case string:
			o.Data.Push(_s_(a_, b_))
			return
		}
	}
	log.Panicf("Cannot add with `+`: %v, %v", a, b)
}

func (o *Terp) JoyInit() {
	o.Prim["+"] = func(o *Terp) {
		o.BinaryOp(
			func(a, b int) int { return a + b },
			func(a, b float64) float64 { return a + b },
			func(a, b string) string { return a + b },
		)
	}
	o.Prim["-"] = func(o *Terp) {
		o.BinaryOp(
			func(a, b int) int { return a - b },
			func(a, b float64) float64 { return a - b },
			nil,
		)
	}
	o.Prim["*"] = func(o *Terp) {
		o.BinaryOp(
			func(a, b int) int { return a * b },
			func(a, b float64) float64 { return a * b },
			nil,
		)
	}
	o.Prim["/"] = func(o *Terp) {
		o.BinaryOp(
			func(a, b int) int { return a / b },
			func(a, b float64) float64 { return a / b },
			nil,
		)
	}
	o.Prim["%"] = func(o *Terp) {
		o.BinaryOp(
			func(a, b int) int { return a % b },
			func(a, b float64) float64 { return math.Mod(a, b) },
			nil,
		)
	}
	o.Prim["."] = func(o *Terp) {
		x := o.Pop()
		fmt.Printf("(%v) ", x)
	}
	o.Prim["emit"] = func(o *Terp) {
		x := o.Pop().(int)
		fmt.Printf("(%c) ", x)
	}
	o.Prim["dup"] = func(o *Terp) {
		x := o.Pop()
		o.Push(x)
		o.Push(x)
	}
	o.Prim["drop"] = func(o *Terp) {
		o.Pop()
	}
	o.Prim["swap"] = func(o *Terp) {
		y := o.Pop()
		x := o.Pop()
		o.Push(y)
		o.Push(x)
	}
	o.Prim["small"] = func(o *Terp) {
		y := o.Pop()
		o.Push(y.(int) < 2)
	}
	/*
		o.Prim["map"] = func(o *Terp) {
			fn := o.Pop()
			vec := o.Pop().([]Any)
			var z []Any
			for _, e := range vec {
				o.Push(e)
				o.Apply(fn)
				z = append(z, o.Pop())
			}
		}
	*/
}
