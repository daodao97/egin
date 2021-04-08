package parser

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

type StructInfo struct {
	Name   string
	Fields []StructField
	Funcs  []StructFunc
	Doc    []string
}

type FuncParam struct {
	Name string
	Type string
}

type StructFunc struct {
	Name        string
	Doc         []string
	Params      []FuncParam
	ResultCount int
}

type StructField struct {
	Name string
	Tags map[string]string
	Type string
}

type VarInfo struct {
	Name string
	Type string
}

type FuncInfo struct {
	Name string
}

// getComment 获取注释信息，来自AST标准库的summary方法
func getComment(group *ast.CommentGroup) (list []string) {
	for _, comment := range group.List {
		// 注释信息会以 // 参数名，开始，我们实际使用时不需要，去掉
		text := strings.TrimSpace(strings.TrimPrefix(comment.Text, "//"))
		list = append(list, text)
	}

	return list
}

func getStruct(f *ast.File) (result []StructInfo) {
	for _, item := range f.Decls {
		obj, ok := item.(*ast.GenDecl)
		if !ok || len(obj.Specs) != 1 {
			continue
		}
		s, ok := obj.Specs[0].(*ast.TypeSpec)
		if !ok {
			continue
		}
		name := s.Name.Name
		body, ok := s.Type.(*ast.StructType)
		var docs []string
		if obj.Doc != nil {
			docs = getComment(obj.Doc)
		}
		if ok {
			var structInfo StructInfo
			structInfo.Name = name
			structInfo.Funcs = getStructFuncDoc(name, f)
			structInfo.Fields = getStructFieldTag(body.Fields.List)
			structInfo.Doc = docs
			result = append(result, structInfo)
		}
	}

	return result
}

func getStructFieldTag(fields []*ast.Field) (result []StructField) {
	for _, v := range fields {
		tags := make(map[string]string)
		if len(v.Names) > 0 && v.Tag != nil {
			tagParts := strings.Split(strings.Trim(v.Tag.Value, "`"), " ")
			for _, v := range tagParts {
				tag := strings.Split(v, ":")
				if len(tag) == 2 {
					tags[strings.TrimSpace(tag[0])] = strings.TrimSpace(strings.Trim(tag[1], "\""))
				}
			}
			result = append(result, StructField{
				Name: v.Names[0].Name,
				Tags: tags,
				// Type: v.Type.(*ast.SelectorExpr).Sel.Name,
			})
		}
	}

	return result
}

func getStructFuncDoc(structName string, f *ast.File) (result []StructFunc) {

	for _, item := range f.Decls {
		fun, ok := item.(*ast.FuncDecl)
		if !ok || (fun.Recv != nil && fun.Recv.List[0].Type.(*ast.Ident).Name != structName) {
			continue
		}
		var docs []string
		if fun.Doc != nil {
			docs = getComment(fun.Doc)
		}
		funcName := fun.Name.Name
		var paramsName []FuncParam
		params := fun.Type.Params.List
		resultCount := 0
		if fun.Type.Results != nil {
			resultCount = len(fun.Type.Results.List)
		}
		for _, v := range params {
			var ptype string
			if pt, ok := v.Type.(*ast.StarExpr); ok {
				ptReal := pt.X.(*ast.SelectorExpr)
				ptype = ptReal.X.(*ast.Ident).Name + "." + ptReal.Sel.Name
			}
			if pt, ok := v.Type.(*ast.Ident); ok {
				ptype = pt.Name
			}
			paramsName = append(paramsName, FuncParam{
				Name: v.Names[0].Name,
				Type: ptype,
			})
		}

		result = append(result, StructFunc{
			Name:        funcName,
			Doc:         docs,
			Params:      paramsName,
			ResultCount: resultCount,
		})
	}

	return result
}

func getType(v interface{}) (ptype string) {
	if pt, ok := v.(*ast.SelectorExpr); ok {
		ptype = pt.X.(*ast.Ident).Name + "." + pt.Sel.Name
	}
	if pt, ok := v.(*ast.Ident); ok {
		ptype = pt.Name
	}
	return ptype
}

func getVarsInfo(f *ast.File) (vars []VarInfo) {
	for _, item := range f.Decls {
		obj, ok := item.(*ast.GenDecl)
		if !ok {
			continue
		}
		if obj.Tok.String() != "var" {
			continue
		}
		if len(obj.Specs) == 0 {
			continue
		}
		body, ok := obj.Specs[0].(*ast.ValueSpec)
		if !ok {
			continue
		}
		var varType string
		if body.Values == nil {
			varType = body.Type.(*ast.Ident).Name
		} else {
			varType = getType(body.Values[0].(*ast.CompositeLit).Type)
		}
		vars = append(vars, VarInfo{
			Name: body.Names[0].Name,
			Type: varType,
		})
	}
	return vars
}

func getFuncInfo(f *ast.File) (result []FuncInfo) {
	for _, item := range f.Decls {
		obj, ok := item.(*ast.FuncDecl)
		if ok {
			result = append(result, FuncInfo{
				Name: obj.Name.Name,
			})
		}
	}
	return
}

func FileStructInfo(fileName string) (result []StructInfo, err error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, fileName, nil, parser.ParseComments)
	if err != nil {
		return result, err
	}
	return getStruct(f), nil
}

func FileVarInfo(fileName string) (result []VarInfo, err error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, fileName, nil, parser.ParseComments)
	if err != nil {
		return result, err
	}
	return getVarsInfo(f), nil
}

func FileFunInfo(fileName string) (result []FuncInfo, err error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, fileName, nil, parser.ParseComments)
	if err != nil {
		return result, err
	}
	return getFuncInfo(f), nil
}
