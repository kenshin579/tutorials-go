package main

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func Test_Parse_Eval(t *testing.T) {
	var node Node
	var result Number
	var p *Parser
	var parseOk, evalOk bool

	data := "3 + 2 - 1"

	p = new(Parser).Init(data)
	p.AddOperator('+', 1)
	p.AddOperator('-', 1)
	p.AddOperator('*', 2)
	p.AddOperator('/', 2)

	node, parseOk = p.Parse()
	require.True(t, parseOk)

	result, evalOk = node.Eval()
	require.True(t, evalOk)

	assert.Equal(t, 3, int(result))

}
