package elseless

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/printer"
	"go/token"

	"github.com/gostaticanalysis/comment"
	"github.com/gostaticanalysis/comment/passes/commentmap"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const doc = "elseless finds unnecessary else"

var Analyzer = &analysis.Analyzer{
	Name: "elseless",
	Doc:  doc,
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
		commentmap.Analyzer,
	},
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.IfStmt)(nil),
	}

	var err error
	inspect.Preorder(nodeFilter, func(n ast.Node) {
		ifstmt, _ := n.(*ast.IfStmt)
		if ifstmt == nil || ifstmt.Else == nil || ifstmt.Body == nil {
			return
		}

		for _, stmt := range ifstmt.Body.List {
			switch stmt := stmt.(type) {
			case *ast.ReturnStmt:
				err = report(pass, ifstmt)
				return
			case *ast.BranchStmt:
				if stmt.Tok == token.BREAK ||
					stmt.Tok == token.CONTINUE {
					err = report(pass, ifstmt)
					return
				}
			}
		}
	})

	if err != nil {
		return nil, err
	}

	return nil, nil
}

func report(pass *analysis.Pass, ifstmt *ast.IfStmt) error {
	cmaps := pass.ResultOf[commentmap.Analyzer].(comment.Maps)

	pos, end := ifstmt.Pos(), ifstmt.End()
	elsestmt := ifstmt.Else
	ifstmt.Else = nil

	var buf bytes.Buffer
	var ifnode interface{} = ifstmt
	if comments := cmaps.Comments(ifstmt); len(comments) > 0 {
		ifnode = &printer.CommentedNode{
			Node:     ifstmt,
			Comments: comments,
		}
	}
	if err := format.Node(&buf, pass.Fset, ifnode); err != nil {
		return err
	}

	fmt.Fprint(&buf, ";")

	var elsenode interface{} = elsestmt
	if comments := cmaps.Comments(elsestmt); len(comments) > 0 {
		elsenode = &printer.CommentedNode{
			Node:     elsestmt,
			Comments: comments,
		}
	}
	if err := format.Node(&buf, pass.Fset, elsenode); err != nil {
		return err
	}

	fix := analysis.SuggestedFix{
		Message: "remove else",
		TextEdits: []analysis.TextEdit{{
			Pos:     pos,
			End:     end,
			NewText: buf.Bytes(),
		}},
	}

	pass.Report(analysis.Diagnostic{
		Pos:            elsestmt.Pos(),
		End:            elsestmt.End(),
		Message:        "unnecessary else",
		SuggestedFixes: []analysis.SuggestedFix{fix},
	})

	return nil
}
