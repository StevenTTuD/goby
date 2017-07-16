package ast

import (
	"bytes"
	"github.com/goby-lang/goby/compiler/token"
)

type ScopeNode interface {
	Name() string
	String() string
}

type ClassStatement struct {
	Token          token.Token
	name           *Constant
	Body           *BlockStatement
	SuperClass     Expression
	SuperClassName string
}

func (cs *ClassStatement) statementNode() {}
func (cs *ClassStatement) TokenLiteral() string {
	return cs.Token.Literal
}
func (cs *ClassStatement) String() string {
	var out bytes.Buffer

	out.WriteString("class ")
	out.WriteString(cs.name.TokenLiteral())
	out.WriteString(" {\n")
	out.WriteString(cs.Body.String())
	out.WriteString("\n}")

	return out.String()
}
func (cs *ClassStatement) Name() string {
	return cs.name.Value
}
func (cs *ClassStatement) SetName(n *Constant) {
	cs.name = n
}

// ModuleStatement represents module node in AST
type ModuleStatement struct {
	Token      token.Token
	name       *Constant
	Body       *BlockStatement
	SuperClass *Constant
}

func (ms *ModuleStatement) statementNode() {}

// TokenLiteral returns token's literal
func (ms *ModuleStatement) TokenLiteral() string {
	return ms.Token.Literal
}
func (ms *ModuleStatement) String() string {
	var out bytes.Buffer

	out.WriteString("module ")
	out.WriteString(ms.name.TokenLiteral())
	out.WriteString(" {\n")
	out.WriteString(ms.Body.String())
	out.WriteString("\n}")

	return out.String()
}
func (ms *ModuleStatement) Name() string {
	return ms.name.Value
}
func (ms *ModuleStatement) SetName(n *Constant) {
	ms.name = n
}

type ReturnStatement struct {
	Token       token.Token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode() {}
func (rs *ReturnStatement) TokenLiteral() string {
	return rs.Token.Literal
}
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral() + " ")

	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}

	out.WriteString(";")

	return out.String()
}

type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (es *ExpressionStatement) statementNode() {}
func (es *ExpressionStatement) TokenLiteral() string {
	return es.Token.Literal
}
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}

	return ""
}

type DefStatement struct {
	Token          token.Token
	name           string
	Receiver       Expression
	Parameters     []Expression
	BlockStatement *BlockStatement
}

func (ds *DefStatement) statementNode() {}
func (ds *DefStatement) TokenLiteral() string {
	return ds.Token.Literal
}
func (ds *DefStatement) String() string {
	var out bytes.Buffer

	out.WriteString("def ")

	if ds.Receiver != nil {
		out.WriteString(ds.Receiver.String() + ".")
	}

	out.WriteString(ds.Name())
	out.WriteString("(")

	for i, param := range ds.Parameters {
		out.WriteString(param.String())
		if i != len(ds.Parameters)-1 {
			out.WriteString(", ")
		}
	}

	out.WriteString(") ")
	out.WriteString("{\n")
	out.WriteString(ds.BlockStatement.String())
	out.WriteString("\n}")

	return out.String()
}
func (ds *DefStatement) Name() string {
	return ds.name
}
func (ds *DefStatement) SetName(n string) {
	ds.name = n
}

// NextStatement represents "next" keyword
type NextStatement struct {
	Token token.Token
}

func (ns *NextStatement) statementNode() {}

// TokenLiteral returns token's literal
func (ns *NextStatement) TokenLiteral() string {
	return ns.Token.Literal
}
func (ns *NextStatement) String() string {
	return "next"
}

type WhileStatement struct {
	Token     token.Token
	Condition Expression
	Body      *BlockStatement
}

func (ws *WhileStatement) statementNode() {}
func (ws *WhileStatement) TokenLiteral() string {
	return ws.Token.Literal
}
func (ws *WhileStatement) String() string {
	var out bytes.Buffer

	out.WriteString("while ")
	out.WriteString(ws.Condition.String())
	out.WriteString(" do\n")
	out.WriteString(ws.Body.String())
	out.WriteString("\nend")

	return out.String()
}

type BlockStatement struct {
	Token      token.Token // {
	Statements []Statement
}

func (bs *BlockStatement) statementNode() {}
func (bs *BlockStatement) TokenLiteral() string {
	return bs.Token.Literal
}
func (bs *BlockStatement) String() string {
	var out bytes.Buffer

	for _, stmt := range bs.Statements {
		out.WriteString(stmt.String())
	}

	return out.String()
}
