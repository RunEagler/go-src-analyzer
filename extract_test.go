package go_src_analyzer

import (
	"testing"
	"fmt"
	"go/ast"
	"go/token"
	"go/parser"
	"github.com/stretchr/testify/assert"
)

func TestExtractFunc(t *testing.T) {

	var actual []Function
	testCase := struct {
		expect []Function
	}{
		expect: []Function{
			{
				name:       "test1",
				startPos:   21,
				endPos:     23,
				isReceiver: false,
				arguments:  []Ident{},
				returns:    []Ident{},
			},
			{
				name:       "test2",
				startPos:   25,
				endPos:     28,
				isReceiver: false,
				arguments: []Ident{
					{
						name:      "arg1",
						kind:      "int",
						isPointer: false,
						identType: Var,
						isChan:    false,
					},
				},
				returns: []Ident{
					{
						name:      "err",
						kind:      "error",
						isPointer: false,
						identType: Var,
						isChan:    false,
					},
				},
			},
			{
				name:       "test3",
				startPos:   48,
				endPos:     50,
				isReceiver: false,
				arguments: []Ident{
					{
						name:      "arg1",
						kind:      "int",
						isPointer: false,
						identType: Var,
						isChan:    false,
					},
					{
						name:      "arg2",
						kind:      "string",
						isPointer: false,
						identType: Var,
						isChan:    false,
					},
					{
						name:      "arg3",
						kind:      "int",
						isPointer: false,
						identType: Array,
						isChan:    false,
					},
					{
						name:      "arg4",
						kind:      "sampleStruct",
						isPointer: false,
						identType: Var,
						isChan:    false,
					},
					{
						name:      "arg5",
						kind:      "int",
						isPointer: true,
						identType: Var,
						isChan:    false,
					},
					{
						name:      "arg6",
						kind:      "string",
						isPointer: true,
						identType: Array,
						isChan:    false,
					}, {
						name:      "arg7",
						kind:      "int",
						isPointer: false,
						identType: Var,
						isChan:    true,
					},
					{
						name:      "arg8",
						kind:      "constType",
						isPointer: false,
						identType: Var,
						isChan:    false,
					},
					{
						name:      "arg9",
						kind:      "constType",
						isPointer: false,
						identType: Array,
						isChan:    false,
					},
					{
						name:      "arg10",
						kind:      "map[int]chan string",
						isPointer: false,
						identType: Map,
						isChan:    false,
					},
					{
						name:      "arg11",
						kind:      "map[constType]*int",
						isPointer: false,
						identType: Map,
						isChan:    false,
					},
					{
						name:      "arg12",
						kind:      "map[struct]string",
						isPointer: true,
						identType: Map,
						isChan:    false,
					},
					{
						name:      "arg13",
						kind:      "func(int) error",
						isPointer: false,
						identType: Func,
						isChan:    false,
					},
					{
						name:      "arg14",
						kind:      "*func(string)",
						isPointer: false,
						identType: Array,
						isChan:    false,
					},
					{
						name:      "arg15",
						kind:      "map[int]func(struct)",
						isPointer: true,
						identType: Map,
						isChan:    false,
					},
					{
						name:      "arg16",
						kind:      "string",
						isPointer: true,
						identType: Var,
						isChan:    true,
					},
					{
						name:      "arg17",
						kind:      "map[chan *int]string",
						isPointer: false,
						identType: Map,
						isChan:    false,
					},
					{
						name:      "arg18",
						kind:      "map[string]map[*func][]int",
						isPointer: false,
						identType: Map,
						isChan:    false,
					},
				},
				returns: []Ident{},
			},
		},
	}

	testFilePath := "./testdata/func.src.go"
	testFileParserPath := "./output/func.src.parser"
	fset := token.NewFileSet()
	astFile, err := parser.ParseFile(fset, testFilePath, nil, parser.Mode(0))
	if err != nil {
		fmt.Println(err)
	}
	writeParser(testFileParserPath, astFile, fset)
	ast.Inspect(astFile, func(n ast.Node) bool {
		switch t := n.(type) {
		case *ast.FuncDecl:
			actual = append(actual, extractFuncDecl(t, fset))
		default:
			fmt.Println(t)
		}
		return true
	})

	equalFunctions(t, testCase.expect, actual)
}

func equalFunctions(t *testing.T, expect, actual []Function) {

	assert.Equal(t, len(expect), len(actual))

	for i, expectFunction := range expect {
		actualFunction := actual[i]
		assert.Equal(t, expectFunction, actualFunction)
	}

}
