package fahrvergnuegen

import (
	"log"
	"strconv"
	"strings"
	"text/scanner"
)

type Word string

type Tok struct {
	scanner.Position
	Text string
	X    Any
}

func Tokenize(src string, filename string) []Tok {
	var zz []Tok
	var s scanner.Scanner
	s.Init(strings.NewReader(src))
	s.Filename = filename

	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		t := s.TokenText()
		i, ierr := strconv.ParseInt(t, 10, 64)
		f, ferr := strconv.ParseFloat(t, 64)

		var x Any
		switch {
		case ierr == nil:
			x = int(i)
		case ferr == nil:
			x = f
		case t[0] == '\'' || t[0] == '"' || t[0] == '`':
			t1, uerr := strconv.Unquote(t)
			if uerr != nil {
				log.Panicf("Bad str token: %v: %q", uerr, t)
			}
			x = t1

		default:
			x = Word(t)
		}

		zz = append(zz, Tok{
			Position: s.Position,
			Text:     s.TokenText(),
			X:        x,
		})
	}
	zz = append(zz, Tok{
		X: nil,
	})
	return zz
}
