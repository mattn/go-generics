package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"log"
	"os"
	//"reflect"
	"strings"
)

func doGenerate(filename string) {
	var buf bytes.Buffer
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filename, nil, parser.ParseComments|parser.SpuriousErrors)
	if err != nil {
		log.Fatal(err)
	}

	typeMap := map[string]*ast.Ident{}
	funcMap := map[string]*ast.FuncDecl{}
	callMap := map[string][]*ast.CallExpr{}

	ast.Inspect(file, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.FuncDecl:
			funcMap[x.Name.Name] = x
		case *ast.TypeSpec:
			if ident, ok := x.Type.(*ast.Ident); ok {
				if ident.Name == "generics" {
					typeMap[x.Name.Name] = ident
					ident.Name = "interface{}"
				}
			}
		case *ast.CallExpr:
			if ident, ok := x.Fun.(*ast.Ident); ok {
				callMap[ident.Name] = append(callMap[ident.Name], x)
			}
		}
		return true
	})
	for _, v := range funcMap {
		if v.Type.Results != nil {
			found := false
			for _, w := range v.Type.Results.List {
				if ident, ok := w.Type.(*ast.Ident); ok {
					if _, ok := typeMap[ident.Name]; ok {
						found = true
						break
					}
				}
			}
			if found {
				for _, call := range callMap[v.Name.Name] {
					typeNames := []string{}
					for _, arg := range call.Args {
						name := ""
						switch a := arg.(type) {
						case *ast.BasicLit:
							name = strings.ToLower(a.Kind.String())
						case *ast.UnaryExpr:
							name = "*" + a.X.(*ast.CompositeLit).Type.(*ast.Ident).Name
						case *ast.CompositeLit:
							name = a.Type.(*ast.Ident).Name
						default:
						}
						if name == "int" {
							name = "int64"
						}
						if name == "float" {
							name = "float64"
						}
						typeNames = append(typeNames, name)
					}
					typeName := strings.Replace(strings.Join(typeNames, "_"), "*", "_", -1)
					if typeName == "" {
						continue
					}
					name := v.Name.Name + "_of_" + typeName
					call.Fun.(*ast.Ident).Name = name

					ret := &ast.FieldList{}
					for n, p := range v.Type.Params.List {
						for m, nm := range p.Names {
							f := &ast.Field{}
							f.Names = []*ast.Ident{{
								token.NoPos,
								nm.Name,
								ast.NewObj(ast.Fun, nm.Name),
							}}
							f.Type = &ast.Ident{
								p.Pos(),
								typeNames[n+m],
								ast.NewObj(ast.Fun, typeNames[n+m]),
							}
							ret.List = append(ret.List, f)
						}
					}

					fd := &ast.FuncDecl{
						Doc:  nil,
						Recv: nil,
						Name: &ast.Ident{
							token.NoPos,
							name,
							ast.NewObj(ast.Fun, name),
						},
						Type: &ast.FuncType{
							Func: token.NoPos,
							//Params:  &ast.FieldList{0, nil, 0},
							Params:  ret,
							Results: v.Type.Results,
						},
						Body: v.Body,
					}
					//rtrn := &ast.ReturnStmt{token.NoPos, nil}
					//fd.Body.List = append(fd.Body.List, rtrn)
					file.Decls = append(file.Decls, fd)
				}
				v.Body = nil
			}
		}
	}
	err = (&printer.Config{
		Mode:     printer.UseSpaces | printer.TabIndent,
		Tabwidth: 8,
	}).Fprint(&buf, fset, file)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(buf.String())
}

func main() {
	for _, arg := range os.Args[1:] {
		doGenerate(arg)
	}
}
