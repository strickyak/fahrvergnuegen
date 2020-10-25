// +build main

/*
  echo ' 3 8 + ! ' | go run fahrvergnuegen.go
*/
package main

import (
	"flag"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime/debug"

	"github.com/chzyer/readline"

	fahr "github.com/strickyak/fahrvergnuegen"
)

var FlagI = flag.Bool("i", false, "run interactive shell even if command line scripts are given")

// var FlagE = flag.Bool("e", false, "don't catch errors")

func main() {
	flag.Parse()
	terp := fahr.NewTerp(emit)
	for _, filename := range flag.Args() {
		bb, err := ioutil.ReadFile(filename)
		if err != nil {
			log.Fatalf("Cannot read file %q: %v", filename, err)
		}
		terp.RunProgram(string(bb), filename)
	}

	if *FlagI || flag.NArg() == 0 {

		home := os.Getenv("HOME")
		if home == "" {
			home = "."
		}
		rl, err := readline.NewEx(&readline.Config{
			Prompt:          " ok ",
			HistoryFile:     filepath.Join(home, ".fahrvergnuegen.history"),
			InterruptPrompt: " *SIGINT* ",
			EOFPrompt:       " *EOF* ",
			// AutoComplete:    completer,
			// HistorySearchFold:   true,
			// FuncFilterInputRune: filterInput,
		})
		if err != nil {
			log.Fatalf("*** ERROR *** in readline: %v", err)
		}
		defer rl.Close()

		for {
			os.Stderr.Write([]byte{'\n'})
			line, err := rl.Readline()
			if err == readline.ErrInterrupt {
				if len(line) == 0 {
					break
				} else {
					continue
				}
			} else if err == io.EOF {
				break
			}
			TryRunProgram(terp, line, "*stdin*")
		}
	}
}

func TryRunProgram(terp *fahr.Terp, line string, filename string) {
	defer func() {
		r := recover()
		if r != nil {
			log.Printf("*** ERROR *** %v <<<<<<<<<<<<<<<<<<<<", r)
			debug.PrintStack()
			log.Printf(">>>>>>>>>>>>>>>>>>>>")
		}
	}()
	terp.RunProgram(line, filename)
}

func emit(r rune) {
	buf := []byte(string([]rune{r}))
	_, err := os.Stdout.Write(buf)
	if err != nil {
		log.Fatalf("cannot emit: %v", err)
	}
}
