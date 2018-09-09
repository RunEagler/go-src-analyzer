package go_src_analyzer

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"go/parser"
	"go/token"
	"go/ast"
	"strings"
	"os"
)

type IdentType int

const (
	Var   IdentType = iota
	Chan
	Array
	Map
	Func
)

type Node struct {
	function Function
}

type Function struct {
	name       string
	startPos   int
	endPos     int
	isReceiver bool
	arguments  []Ident
	returns    []Ident
}

type Ident struct {
	name      string
	kind      string
	isPointer bool
	isChan    bool
	identType IdentType
}

func main() {

	//var filename = "./output/go/codes.csv"
	//
	//file, err := os.Open(filename)
	//if err != nil {
	//	fmt.Println(file)
	//}
	//
	//defer file.Close()
	//
	//r := csv.NewReader(file)
	//for i := 0; i < 1; i++ {
	//	line, _ := r.Read()
	//	if len(line) == 0 {
	//		break
	//	}
	//	URL := line[3]
	//	body := getRequest(URL)
	//	goFileName := URL[strings.LastIndex(URL, "/"):strings.LastIndex(URL, ".")]
	//	file, err := os.Create(fmt.Sprintf("./output/go_file/%s.go", goFileName))
	//	if err != nil {
	//		fmt.Println(err)
	//	}
	//	file.WriteString(body)
	//	defer file.Close()
	//
	//	codeAnalyze(body, goFileName)
	//}

	codeAnalyze("./output/go_file/session.go")

	fmt.Println("end")
}

func codeAnalyze(goFileName string) {

	//var functions []Function
	var comments []string

	fset := token.NewFileSet()
	expr, err := parser.ParseFile(fset, goFileName, nil, parser.Mode(0))
	if err != nil {
		fmt.Println(err)
	}

	analyze(goFileName, expr, fset)
	comments = analyzeComments(goFileName)
	fmt.Println(comments)

}

func analyze(goFileName string, astFile *ast.File, fileSet *token.FileSet) {

	var functions []Function
	ast.Inspect(astFile, func(n ast.Node) bool {
		switch t := n.(type) {
		case *ast.FuncDecl:
			functions = append(functions, extractFuncDecl(t, fileSet))
		default:
			fmt.Println(t)
		}
		return true
	})

	//printFunctions(functions)

	//defer resultFile.Close()
}

func analyzeComments(goFileName string) []string {

	var comments []string
	fset := token.NewFileSet()
	expr, err := parser.ParseFile(fset, goFileName, nil, parser.ParseComments)
	if err != nil {
		fmt.Println(err)
	}
	for _, commentGroup := range expr.Comments {
		var commentSentences string
		for _, comment := range commentGroup.List {
			commentSentences += comment.Text
		}
		commentSentences = strings.Replace(commentSentences, "//", "", -1)
		comments = append(comments, commentSentences)
	}
	return comments

}

func printFunctions(functions []Function) {

	for _, function := range functions {

		fmt.Println(fmt.Sprintf("%d:%d funcName=%s,isReceiver=%t", function.startPos, function.endPos, function.name, function.isReceiver))

		fmt.Println("==arguments==")
		for _, argument := range function.arguments {

			fmt.Println(fmt.Sprintf("name=%s,type=%s,isPointer=%t,identType=%s", argument.name, argument.kind, argument.isPointer, identTypeString(argument.identType)))

		}
		fmt.Println("==returns==")
		for _, ret := range function.returns {

			fmt.Println(fmt.Sprintf("name=%s,type=%s,isPointer=%t,identType=%s", ret.name, ret.kind, ret.isPointer, identTypeString(ret.identType)))

		}
	}

}

func getRequest(URL string) string {

	resp, _ := http.Get(URL)
	defer resp.Body.Close()
	byteArray, _ := ioutil.ReadAll(resp.Body)
	return string(byteArray)
}

func identTypeString(identType IdentType) string {

	var str string

	switch(identType) {
	case Var:
		str = "Var"
	case Array:
		str = "Array"
	case Map:
		str = "Map"
	case Chan:
		str = "Chan"
	case Func:
		str = "Func"
	}
	return str
}
func writeParser(goFilePath string, astFile *ast.File, fileSet *token.FileSet) {
	parserFile, err := os.Create(goFilePath)
	if err != nil {
		fmt.Println(err)
	}
	for _, d := range astFile.Decls {
		ast.Fprint(parserFile, fileSet, d, nil)
	}

	defer parserFile.Close()
}
