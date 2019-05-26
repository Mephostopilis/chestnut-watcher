package parser

import (
	"testing"
)

func TestIntegerLiteralExpression(t *testing.T) {

}

func TestInterLiteral(t *testing.T, il Expression, value int64) {
	integ, ok := il.(*IntegerLiteralExpression)
	if !ok {
		t.Errorf("il not *ast")
		return false
	}
	
}