package main

import (
	"fmt"
	"log"

	"github.com/lo1kx/math_parser/lib/parser"
)

func main() {
	p := parser.NewExpressionParser("(4/2+12)/7")

	err := p.Validate()
	if err != nil {
		fmt.Printf("validation: %s", err)
		return
	}

	err = p.Parse()
	if err != nil {
		log.Fatal(err)
	}

	res, err := p.Root.Eval()
	fmt.Println("Good: ", res, err)
}
