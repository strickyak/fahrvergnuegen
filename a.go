package fahrvergnuegen

import (
	"fmt"
	"log"
)

type (
	Terp struct {
		Emit func(r rune)
		Data Stack
		Prim map[Word]func(*Terp)
		Defs map[Word]Any
	}

	Stack struct {
		V []Any
	}

	Frame struct {
		V map[string]Any
	}

	Any interface{}
)

func (terp *Terp) Push(x Any) {
	terp.Data.Push(x)
}
func (terp *Terp) Pop() Any {
	return terp.Data.Pop()
}

func (terp *Stack) Push(x Any) {
	terp.V = append(terp.V, x)
}

func (terp *Stack) Len() int {
	return len(terp.V)
}

func (terp *Stack) Pop() Any {
	z := terp.V[len(terp.V)-1]
	terp.V = terp.V[:len(terp.V)-1]
	return z
}

func (terp *Stack) String() string {
	return fmt.Sprintf("%v", terp.V)
}

func NewTerp(emit func(r rune)) *Terp {
	terp := &Terp{
		Emit: emit,
		Prim: make(map[Word]func(*Terp)),
		Defs: make(map[Word]Any),
	}
	/*
		terp.Prim[Word("+")] = func(o *Terp) {
			b := o.Data.Pop()
			a := o.Data.Pop()

			switch a_ := a.(type) {
			case int:
				switch b_ := b.(type) {
				case int:
					o.Data.Push(a_ + b_)
					return
				case float64:
					o.Data.Push(float64(a_) + b_)
					return
				}
			case float64:
				switch b_ := b.(type) {
				case int:
					o.Data.Push(a_ + float64(b_))
					return
				case float64:
					o.Data.Push(a_ + b_)
					return
				}
			}
			log.Panicf("Cannot add with `+`: %v, %v", a, b)
		}
	*/
	terp.Prim[Word("!")] = func(o *Terp) {
		s := o.Data.Pop().(string)
		for _, r := range s {
			o.Emit('<')
			o.Emit(r)
			o.Emit('>')
		}
	}
	terp.JoyInit()
	return terp
}

func (terp *Terp) RunProgram(s string, filename string) {
	tt := Tokenize(s, filename)

	for i, t := range tt {
		_ = i
		// log.Printf("Step[%d]: %v", i, t)
		terp.Step(t)
	}

	if terp.Data.Len() > 0 {
		log.Printf("Final Stack: %v", terp.Data)
	}
}

func (terp *Terp) Step(t Tok) {
	switch a := t.X.(type) {
	case nil:
		return
	case Word:
		prim, ok := terp.Prim[a]
		if !ok {
			def, dok := terp.Defs[a]
			if !dok {
				log.Panicf("op not defined: %q", a)
			}
			log.Panicf("TODO: Def %q === %v", a, def)
		}
		prim(terp)
		return

	default:
		terp.Data.Push(t.X)
		return
	}
}
