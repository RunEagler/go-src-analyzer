package go_src_analyzer

import (
	"go/ast"
	"fmt"
	"go/token"
	"strings"
)

func extractSpecs(specs []ast.Spec) {

	for _, spec := range specs {
		extractSpec(spec.(*ast.TypeSpec))
	}
}

func extractSpec(spec ast.Spec) {

	var idents []Ident
	switch t := spec.(type) {
	case *ast.ValueSpec:

	case *ast.TypeSpec:

		extractTypeSpec(t)
	}
	fmt.Println(idents)

}

func extractValueSpec(valueSpec *ast.ValueSpec) {

	var name string
	var kind string
	var isPointer bool

	for _, ident := range valueSpec.Names {
		name = ident.Name
	}

	if valueSpec.Type != nil {
		switch t := valueSpec.Type.(type) {
		case *ast.StarExpr:
			kind, isPointer = extractStarExpr(t)
		case *ast.ArrayType:
			kind, isPointer = extractArrayType(t)
		}
	}
	fmt.Println(name, kind, isPointer)
}

func extractTypeSpec(spec *ast.TypeSpec) {

	var name string
	var idents []Ident

	name = spec.Name.Name
	switch t := spec.Type.(type) {
	case *ast.StructType:
		idents = extractFieldList(t.Fields.List)
	}
	fmt.Println(name, idents)
}

func extractStructType(structType *ast.StructType) []Ident {

	var idents []Ident

	switch t := interface{}(structType).(type) {
	case []*ast.Field:
		idents = extractFieldList(t)
	default:
		fmt.Println(t)
	}

	return idents

}

func extractFuncDecl(function *ast.FuncDecl, fileSet *token.FileSet) Function {

	var funcName string
	var arguments []Ident
	var returns []Ident
	var isReceiver bool

	startPos := fileSet.Position(function.Body.Lbrace).Line
	endPos := fileSet.Position(function.Body.Rbrace).Line

	if function.Name.Obj != nil {
		funcName = function.Name.Obj.Name
	} else {
		funcName = function.Name.Name
	}
	if function.Recv != nil {
		isReceiver = true
	}

	arguments, returns = extractFuncType(function.Type)

	return Function{
		startPos:   startPos,
		endPos:     endPos,
		isReceiver: isReceiver,
		name:       funcName,
		arguments:  arguments,
		returns:    returns,
	}
}

func extractFieldList(fields []*ast.Field) []Ident {

	var idents []Ident
	for _, result := range fields {
		var ident Ident
		ident = extractField(result)
		idents = append(idents, ident)
	}
	return idents
}

func extractField(field *ast.Field) Ident {

	var kind string
	var name string
	var isPointer bool
	var identType IdentType

	for _, fieldName := range field.Names {
		name = fieldName.Name
	}
	switch t := interface{}(field.Type).(type) {
	case *ast.Ident:
		kind = t.Name
	case *ast.SelectorExpr:
		kind = extractSelectorExpr(t)
	case *ast.Ellipsis:
		kind, isPointer = extractEllipsis(t)
	case *ast.StarExpr:
		kind, isPointer = extractStarExpr(t)
	case *ast.ArrayType:
		identType = Array
		kind, isPointer = extractArrayType(t)
	case *ast.ChanType:
		identType = Chan
	case *ast.MapType:
		kind = extractMapType(t)
	default:
		fmt.Println(t)
	}
	return Ident{
		name:      name,
		kind:      kind,
		isPointer: isPointer,
		identType: identType,
	}
}

func extractEllipsis(ellipsis *ast.Ellipsis) (string, bool) {

	var kind string
	var isPointer bool
	switch t := ellipsis.Elt.(type) {
	case *ast.StarExpr:
		kind, isPointer = extractStarExpr(t)
	default:
		fmt.Println(kind)
	}
	return "..." + kind, isPointer
}

func extractChanType(chanType *ast.ChanType) (string, bool) {

	var kind string
	switch t := chanType.Value.(type) {
	case *ast.Ident:
		kind = t.Name
	case *ast.StructType:
		kind = "struct"

	}
	return kind, false
}

func extractMapType(mapType *ast.MapType) (string) {

	var kind string
	var valueKind string
	var keyKind string
	//var isKeyPointer bool
	//var isValuePointer bool

	keyKind = extractMapKeyType(mapType.Key)
	valueKind = extractMapValueType(mapType.Value)

	kind = fmt.Sprintf("map[%s]%s", keyKind, valueKind)

	return kind
}
func extractMapKeyType(mapKeyType ast.Expr) (string) {
	var kind string
	var isPointer bool

	switch t := mapKeyType.(type) {
	case *ast.Ident:
		kind = t.Name
	case *ast.StarExpr:
		kind, isPointer = extractStarExpr(t)
	case *ast.SelectorExpr:
		kind = "struct"
	case *ast.MapType:
		kind = extractMapType(t)
	}
	if isPointer {
		kind = "*" + kind
	}

	return kind
}
func extractMapValueType(mapValueType ast.Expr) (string) {

	var kind string
	var isPointer bool

	switch t := mapValueType.(type) {
	case *ast.Ident:
		kind = t.Name
	case *ast.StarExpr:
		kind, isPointer = extractStarExpr(t)
	case *ast.SelectorExpr:
		kind = extractSelectorExpr(t)
	case *ast.MapType:
		kind = extractMapType(t)
	case *ast.ArrayType:
		kind, isPointer = extractArrayType(t)
	case *ast.ChanType:
		kind, isPointer = extractChanType(t)
	case *ast.FuncType:
		_, _ = extractFuncType(t)
		kind = "func"
	}
	if isPointer {
		kind = "*" + kind
	}

	return kind
}

func extractExpr(expr ast.Expr) (string, bool) {
	var kind string
	var isPointer bool

	switch t := expr.(type) {
	case *ast.Ident:
		kind = t.Name
	case *ast.StarExpr:
		kind, isPointer = extractStarExpr(t)
	case *ast.SelectorExpr:
		kind = extractSelectorExpr(t)
	case *ast.MapType:
		kind = extractMapValueType(t)
	case *ast.ArrayType:
		kind, isPointer = extractArrayType(t)
	case *ast.ChanType:
		kind, isPointer = extractChanType(t)
	case *ast.StructType:
		//kind = extractStructType(t)
	case *ast.InterfaceType:
		kind = extractInterfaceType(t)
	case *ast.FuncType:
		_, _ = extractFuncType(t)
		kind = "func"

	}
	return kind, isPointer
}

func extractInterfaceType(interfaceType *ast.InterfaceType) string {

	var kind string

	return kind
}

func extractFuncType(funcType *ast.FuncType) ([]Ident, []Ident) {

	var arguments []Ident
	var returns []Ident

	arguments = extractFieldList(funcType.Params.List)
	if funcType.Results != nil {
		returns = extractFieldList(funcType.Results.List)
	}
	return arguments, returns
}

func extractStarExpr(starExpr *ast.StarExpr) (string, bool) {

	var kind string
	var isPointer bool
	kind, isPointer = extractExpr(starExpr.X)

	return kind, isPointer
}
func extractSelectorExpr(selectExpr *ast.SelectorExpr) string {

	var kind string
	var x string
	var sel string
	switch t := selectExpr.X.(type) {
	case *ast.Ident:
		x = t.Name
	default:
		fmt.Println(t)
	}
	switch t := interface{}(selectExpr.Sel).(type) {
	case *ast.Ident:
		sel = t.Name
	default:
		fmt.Println(t)
	}
	kind = strings.Join([]string{x, sel}, ".")
	return kind
}

func extractArrayType(arrayType *ast.ArrayType) (string, bool) {

	var kind string
	var isPointer bool

	switch t := arrayType.Elt.(type) {
	case *ast.StarExpr:
		kind, isPointer = extractStarExpr(t)
	default:
		fmt.Println(t)
	}
	return kind, isPointer

}
