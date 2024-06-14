package ast

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"

	"golang.org/x/tools/go/ast/astutil"
)

type instrumenter struct {
	traceImport string
	tracePkg    string
	traceFunc   string
}

/*
*

	在计算机科学中，抽象语法树（abstract syntax tree，AST）是源代码的抽象语法结构的树状表现形式，树上的每个节点都表示源代码中的一种结构。
	因为 Go 语言是开源编程语言，所以它的抽象语法树的操作包也和语言一起开放给了 Go 开发人员，我们可以基于 Go 标准库以及Go 实验工具库提供的 ast 相关包，
	快速地构建基于 AST 的应用，这里的 ast.instrumenter 就是一个应用 AST 的典型例子。
	一旦我们通过 ast 相关包解析 Go 源码得到相应的抽象语法树后，我们便可以操作这棵语法树，并按我们的逻辑在语法树中注入我们的 Trace 函数，最后我们再将修改后的抽象语法树转换为 Go 源码，
	就完成了整个自动注入的工作了。
*/
func New(traceImport, tracePkg, traceFunc string) *instrumenter {
	return &instrumenter{
		traceImport: traceImport,
		tracePkg:    tracePkg,
		traceFunc:   traceFunc,
	}
}

func hasFuncDecl(f *ast.File) bool {
	if len(f.Decls) == 0 {
		return false
	}
	// 遍历 AST 树中的顶层声明，如果整个源码都不包含函数声明，则无需注入操作，直接返回。
	for _, decl := range f.Decls {
		_, ok := decl.(*ast.FuncDecl)
		if ok {
			return true
		}
	}

	return false
}

/*
*
Instrument 首先通过 go/paser 的 ParserFile 函数对传入的 Go 源文件中的源码进行解析，并得到对应的抽象语法树 AST，然后向 AST 中导入 Trace 函数所在的包，
并向这个 AST 的所有函数声明注入 Trace 函数调用。
*/
func (a instrumenter) Instrument(filename string) ([]byte, error) {
	// 用于创建一个新的 token.FileSet 实例。token.FileSet 是 go/token 包中的一个类型，主要用于管理和跟踪源代码文件的位置信息。这个位置包括每个标记（token）的文件、行和列信息。
	fset := token.NewFileSet()
	// filename是要解析的 Go 源代码文件的路径 curAST是生成的AST语法树
	curAST, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("error parsing %s: %w", filename, err)
	}

	if !hasFuncDecl(curAST) { // 如果整个源码都不包含函数声明，则无需注入操作，直接返回。
		return nil, nil
	}

	// add import declaration 在AST上添加包导入语句
	astutil.AddImport(fset, curAST, a.traceImport)

	// inject code into each function declaration 向AST上的所有函数注入Trace函数
	a.addDeferTraceIntoFuncDecls(curAST)

	buf := &bytes.Buffer{}
	err = format.Node(buf, fset, curAST) // 将修改后的AST转换回Go源码
	if err != nil {
		return nil, fmt.Errorf("error formatting new code: %w", err)
	}
	return buf.Bytes(), nil
}

func (a instrumenter) addDeferTraceIntoFuncDecls(f *ast.File) {
	for _, decl := range f.Decls {
		fd, ok := decl.(*ast.FuncDecl)
		if ok {
			// inject code to fd => 如果是函数声明，则注入跟踪设施
			a.addDeferStmt(fd)
		}
	}
}

func (a instrumenter) addDeferStmt(fd *ast.FuncDecl) (added bool) {
	stmts := fd.Body.List

	// check whether "defer trace.Trace()()" has already exists => 判断"defer trace.Trace()()"语句是否已经存在
	for _, stmt := range stmts {
		ds, ok := stmt.(*ast.DeferStmt)
		if !ok {
			// 如果不是defer语句，则继续for循环
			continue
		}
		// it is a defer stmt => 如果是defer语句，则要进一步判断是否是defer trace.Trace()()
		ce, ok := ds.Call.Fun.(*ast.CallExpr)
		if !ok {
			continue
		}

		se, ok := ce.Fun.(*ast.SelectorExpr)
		if !ok {
			continue
		}

		x, ok := se.X.(*ast.Ident)
		if !ok {
			continue
		}
		if (x.Name == a.tracePkg) && (se.Sel.Name == a.traceFunc) {
			// already exist , return
			return false
		}
	}

	// not found "defer trace.Trace()()"
	// add one
	// 没有找到"defer trace.Trace()()"，注入一个新的跟踪语句
	// 在AST上构造一个defer trace.Trace()()
	ds := &ast.DeferStmt{
		Call: &ast.CallExpr{
			Fun: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X: &ast.Ident{
						Name: a.tracePkg,
					},
					Sel: &ast.Ident{
						Name: a.traceFunc,
					},
				},
			},
		},
	}

	newList := make([]ast.Stmt, len(stmts)+1)
	copy(newList[1:], stmts)
	newList[0] = ds // 注入新构造的defer语句
	fd.Body.List = newList
	return true
}
