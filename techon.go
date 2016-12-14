package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/noonien/techon/parser"
	"github.com/noonien/techon/runner"
)

func main() {
	p := parser.NewParser(os.Stdin)
	prog, err := p.Parse()
	if err != nil {
		log.Fatal(err)
	}

	m := runner.NewMachine()
	err = m.Execute(prog)
	if err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(os.Stdout).Encode(m.Stack)
}
