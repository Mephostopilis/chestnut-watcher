package main

import (
	"fmt"
	"gotutorial/tutorial1/parser"
)

func evalPrefixExpression(operator string, right int64) int64 {
	return 0
}

func evalInfixExpression(left int64, operator string, right int64) int64 {
	return 0
}

func Eval(exp parser.Expression) int64 {
	switch node := exp.(type) {
	case *parser.IntegerLiteralExpression:
		return node.Value
	case *parser.PrefixExpression:
		// rightV := Eval(node.Right)
	}
	return 0
}

func calc(input string) int64 {
	return 0
}

func main() {
	var n = calc("12*12")
	fmt.Println("%d", n)
}
