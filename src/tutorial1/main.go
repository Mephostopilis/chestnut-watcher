package main

func evalPrefixExpression(operator string, right int64) int64 {

}

func evalInfixExpression(left int64, operator string, right int64) int64 {
	
}

func Eval(exp Expression) int64 {
	switch node := exp.(type) {
	case *IntegerLiteralExpression:
		return node.Value
	case *PrefixExpression:
		rightV := Eval(node.Right)
		
	}
}

func calc(input string) int64 {
	return 0
}

func main() {

}
