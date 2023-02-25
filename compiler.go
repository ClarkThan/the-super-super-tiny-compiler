package main

import (
	"fmt"
	"os"
	"strings"
)

const (
	Name       string = "Name"
	Paren      string = "Paren"
	Number     string = "Number"
	String     string = "String"
	Expression string = "Expression"
	Program    string = "Program"
)

// data structure for Tokenizer phrase
type Token struct {
	typ, val string
}

// data structure for Parser and CodeGen phrase
type Node interface {
	Type() string
}

type Literal struct {
	typ, val string
}

func (n Literal) Type() string {
	return n.typ
}

type Expr struct {
	typ, name string
	params    []Node
}

func (n Expr) Type() string {
	return n.typ
}

type Ast struct {
	typ  string
	body []Node
}

func (n Ast) Type() string {
	return n.typ
}

func main() {
	src := `(add 3 (sub 4 (len "foo")))` // add(3, sub(4, len("foo")));
	if len(os.Args) > 1 {
		src = os.Args[1]
	}
	code := Compile(src)
	fmt.Println(code)
}

func Compile(src string) string {
	tokens := Tokenizer(src)
	// fmt.Println(tokens)
	// [{Paren (} {Name add} {Number 3} {Paren (} {Name sub} {Number 4} {Paren (} {Name len} {String foo} {Paren )} {Paren )} {Paren )}]
	ast := Parser(tokens)
	// fmt.Println(ast)
	// {Program [
	// 		{Expression add [
	// 			{Number 3}
	// 			{Expression sub [
	// 				{Number 4}
	// 				{Expression len [
	// 					{String foo}
	// 				]}
	// 			]}
	// 		]}
	// ]}
	code := CodeGen(ast)
	return code
}

func Tokenizer(src string) (tokens []Token) {
	if len(src) == 0 {
		return
	}

	chars := []rune(src)
	idx, length := 0, len(chars)
	chars = append(chars, '\n')
	for idx < length {
		c := chars[idx]
		if c == '(' || c == ')' {
			tokens = append(tokens, Token{Paren, string(c)})
			idx++
			continue
		}
		if isNum(c) {
			numIdx := idx + 1
			for c = chars[numIdx]; isNum(c); c = chars[numIdx] {
				numIdx++
			}
			tokens = append(tokens, Token{Number, string(chars[idx:numIdx])})
			idx = numIdx
			continue
		}
		if isSpace(c) {
			idx++
			for c = chars[idx]; isSpace(c); c = chars[idx] {
				idx++
			}
			continue
		}
		if c == '"' {
			strIdx := idx + 1
			for c = chars[strIdx]; c != '"'; c = chars[strIdx] {
				strIdx++
			}
			tokens = append(tokens, Token{String, string(chars[idx+1 : strIdx])})
			idx = strIdx + 1 // skip the right `"`
			continue
		}
		if isAlpha(c) {
			alphaIdx := idx + 1
			for c = chars[alphaIdx]; isAlpha(c); c = chars[alphaIdx] {
				alphaIdx++
			}
			tokens = append(tokens, Token{Name, string(chars[idx:alphaIdx])})
			idx = alphaIdx
			continue
		}
		panic(fmt.Sprintf(`Unknown char[%d]: %c`, idx, c))
	}

	return
}

func Parser(tokens []Token) Ast {
	t := &Tracer{tokens: tokens}

	var stmts []Node
	for t.idx < len(tokens) {
		stmts = append(stmts, t.walk())
	}

	return Ast{
		typ:  Program,
		body: stmts,
	}
}

func CodeGen(n Node) (code string) {
	switch n.Type() {
	case Number:
		n := n.(Literal)
		return n.val
	case String:
		n := n.(Literal)
		return fmt.Sprintf("%q", n.val)
	case Expression:
		n := n.(Expr)
		var args []string
		for _, para := range n.params {
			args = append(args, CodeGen(para))
		}
		return fmt.Sprintf("%s(%s)", n.name, strings.Join(args, ", "))
	case Program:
		n := n.(Ast)
		var lines []string
		for _, b := range n.body {
			lines = append(lines, CodeGen(b)+";")
		}
		return strings.Join(lines, "\n")
	default:
		return ""
	}
}

type Tracer struct {
	tokens []Token
	idx    int
}

func (t *Tracer) walk() Node {
	tok := t.tokens[t.idx]
	if tok.typ == Number || tok.typ == String {
		t.idx++
		return Literal{tok.typ, tok.val}
	}
	if tok.val == "(" {
		t.idx += 2
		expr := Expr{
			typ:  Expression,
			name: t.tokens[t.idx-1].val,
		}
		for tok = t.tokens[t.idx]; tok.val != ")"; tok = t.tokens[t.idx] {
			expr.params = append(expr.params, t.walk())
		}
		t.idx++ // skip the closing parentheses: `)`
		return expr
	}
	panic(fmt.Sprintf("Syntax error - token[%d]: %s", t.idx, tok))
}

func isNum(c rune) bool {
	if c >= '0' && c <= '9' {
		return true
	}
	return false
}

func isAlpha(c rune) bool {
	if (c >= 'a' && c < 'z') || (c >= 'A' && c < 'Z') {
		return true
	}
	return false
}

func isSpace(c rune) bool {
	return c == ' ' || c == '\t' || c == '\n'
}
