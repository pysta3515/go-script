package go_script

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
	"text/scanner"
)

type Expr interface {
	Eval() Expr
	ToInt(env map[string]interface{}) int
	ToInt64(env map[string]interface{}) int64
	ToFloat(env map[string]interface{}) float64
	ToBool(env map[string]interface{}) bool
}

var op = map[string]bool{
	"+":   true,
	"-":   true,
	"*":   true,
	"/":   true,
	">":   true,
	">=":  true,
	"<":   true,
	"<=":  true,
	"==":  true,
	"!=":  true,
	"&&":  true,
	"and": true,
	"or":  true,
	"||":  true,
	"if":  true,
	"in":  true,
}

// Binary 抽象一个二元或多远运算
type Binary struct {
	OP   string
	Args []Expr
}

func (binary *Binary) Eval() Expr {
	return nil
}

func (binary *Binary) ToInt(env map[string]interface{}) int {
	switch binary.OP {
	case "+":
		res := 0
		for _, exp := range binary.Args {
			res += exp.ToInt(env)
		}
		return res
	case "-":
		res := binary.Args[0].ToInt(env)
		for _, exp := range binary.Args[1:] {
			res -= exp.ToInt(env)
		}
		return res
	case "*":
		res := 1
		for _, exp := range binary.Args {
			res *= exp.ToInt(env)
		}
		return res
	case "/":
		res := binary.Args[0].ToInt(env)
		for _, exp := range binary.Args[1:] {
			res /= exp.ToInt(env)
		}
		return res
	case "if":
		if binary.Args[0].ToBool(env) {
			return binary.Args[1].ToInt(env)
		}
		return binary.Args[2].ToInt(env)
	}
	panic(fmt.Sprintf("invalid method: %s", binary.OP))
}

func (binary *Binary) ToInt64(env map[string]interface{}) int64 {
	switch binary.OP {
	case "+":
		res := int64(0)
		for _, exp := range binary.Args {
			res += exp.ToInt64(env)
		}
		return res
	case "-":
		res := binary.Args[0].ToInt64(env)
		for _, exp := range binary.Args[1:] {
			res -= exp.ToInt64(env)
		}
		return res
	case "*":
		res := int64(1)
		for _, exp := range binary.Args {
			res *= exp.ToInt64(env)
		}
		return res
	case "/":
		res := binary.Args[0].ToInt64(env)
		for _, exp := range binary.Args[1:] {
			res /= exp.ToInt64(env)
		}
		return res
	case "if":
		if binary.Args[0].ToBool(env) {
			return binary.Args[1].ToInt64(env)
		}
		return binary.Args[2].ToInt64(env)
	}
	panic(fmt.Sprintf("invalid method: %s", binary.OP))
}

func (binary *Binary) ToFloat(env map[string]interface{}) float64 {
	switch binary.OP {
	case "+":
		res := float64(0)
		for _, exp := range binary.Args {
			res += exp.ToFloat(env)
		}
		return res
	case "-":
		res := binary.Args[0].ToFloat(env)
		for _, exp := range binary.Args[1:] {
			res -= exp.ToFloat(env)
		}
		return res
	case "*":
		res := float64(1)
		for _, exp := range binary.Args {
			res *= exp.ToFloat(env)
		}
		return res
	case "/":
		res := binary.Args[0].ToFloat(env)
		for _, exp := range binary.Args[1:] {
			res /= exp.ToFloat(env)
		}
		return res
	case "if":
		if binary.Args[0].ToBool(env) {
			return binary.Args[1].ToFloat(env)
		}
		return binary.Args[2].ToFloat(env)
	}
	panic(fmt.Sprintf("invalid method: %s", binary.OP))
}

func (binary *Binary) ToBool(env map[string]interface{}) bool {
	switch binary.OP {
	case ">":
		return binary.Args[0].ToFloat(env) > binary.Args[1].ToFloat(env)
	case ">=":
		return binary.Args[0].ToFloat(env) >= binary.Args[1].ToFloat(env)
	case "<":
		return binary.Args[0].ToFloat(env) < binary.Args[1].ToFloat(env)
	case "<=":
		return binary.Args[0].ToFloat(env) <= binary.Args[1].ToFloat(env)
	case "+", "-", "*", "/":
		return binary.ToFloat(env) != 0
	case "==":
		return binary.Args[0].ToFloat(env) == binary.Args[1].ToFloat(env)
	case "!=":
		return binary.Args[0].ToFloat(env) != binary.Args[1].ToFloat(env)
	case "if":
		if binary.Args[0].ToBool(env) {
			return binary.Args[1].ToBool(env)
		}
		return binary.Args[2].ToBool(env)
	case "and", "&&":
		for _, exp := range binary.Args {
			if !exp.ToBool(env) {
				return false
			}
		}
		return true
	case "or", "||":
		for _, exp := range binary.Args {
			if exp.ToBool(env) {
				return true
			}
		}
		return false
	case "in":
		arg0 := binary.Args[0].ToFloat(env)
		for _, arg := range binary.Args[1:] {
			if arg0 == arg.ToFloat(env) {
				return true
			}
		}
		return false
	}
	panic(fmt.Sprintf("invalid method: %s", binary.OP))
}

// Atom 抽象一个最简单的常量表达式
type Atom struct {
	typ  rune // -4 float -3 int -2 var
	text string
}

func (node *Atom) Eval() Expr {
	panic("atom usage error")
}

func (node *Atom) ToInt(_ map[string]interface{}) int {
	switch node.typ {
	case -3: // int
		i, err := strconv.Atoi(node.text)
		if err != nil {
			panic(err)
		}
		return i
	case -4:
		f, err := strconv.ParseFloat(node.text, 64)
		if err != nil {
			panic(err)
		}
		return int(f)
	}
	panic(fmt.Sprintf("invalid value [%s] with type: %s", node.text, scanner.TokenString(node.typ)))

}

func (node *Atom) ToInt64(_ map[string]interface{}) int64 {
	switch node.typ {
	case -3: // int
		i, err := strconv.ParseInt(node.text, 10, 0)
		if err != nil {
			panic(err)
		}
		return i
	case -4:
		f, err := strconv.ParseFloat(node.text, 64)
		if err != nil {
			panic(err)
		}
		return int64(f)
	}
	panic(fmt.Sprintf("invalid value [%s] with type: %s", node.text, scanner.TokenString(node.typ)))
}

func (node *Atom) ToFloat(_ map[string]interface{}) float64 {
	switch node.typ {
	case -4:
		f, err := strconv.ParseFloat(node.text, 64)
		if err != nil {
			panic(err)
		}
		return f
	case -3: // int
		i, err := strconv.Atoi(node.text)
		if err != nil {
			panic(err)
		}
		return float64(i)
	}
	panic(fmt.Sprintf("invalid value [%s] with type: %s", node.text, scanner.TokenString(node.typ)))
}

func (node *Atom) ToBool(_ map[string]interface{}) bool {
	switch node.typ {
	case -3: // int
		i, err := strconv.Atoi(node.text)
		if err != nil {
			panic(err)
		}
		return i != 0
	case -4:
		f, err := strconv.ParseFloat(node.text, 64)
		if err != nil {
			panic(err)
		}
		return f != 0
	}
	panic(fmt.Sprintf("invalid value [%s] with type: %s", node.text, scanner.TokenString(node.typ)))
}

// Var 抽象一个变量引用
type Var struct {
	Name string
	Val  interface{}
}

func (v *Var) ToInt(env map[string]interface{}) int {
	switch v := env[v.Name].(type) {
	case int:
		return v
	case int64:
		return int(v)
	case float64:
		return int(v)
	}
	panic(fmt.Sprintf("invalid env[%s] type: %s", v.Name, reflect.TypeOf(v)))
}

func (v *Var) ToInt64(env map[string]interface{}) int64 {
	switch v := env[v.Name].(type) {
	case int64:
		return v
	case int:
		return int64(v)
	case float64:
		return int64(v)
	}
	panic(fmt.Sprintf("invalid env[%s] type: %s", v.Name, reflect.TypeOf(v)))
}

func (v *Var) ToFloat(env map[string]interface{}) float64 {
	switch v := env[v.Name].(type) {
	case float64:
		return v
	case int:
		return float64(v)
	case int64:
		return float64(v)
	}
	panic(fmt.Sprintf("invalid env[%s] type: %s", v.Name, reflect.TypeOf(v)))
}

func (v *Var) ToBool(env map[string]interface{}) bool {
	switch v := env[v.Name].(type) {
	case int:
		return v != 0
	case int64:
		return v != 0
	case float64:
		return v != 0
	}
	panic(fmt.Sprintf("invalid env[%s] type: %s", v.Name, reflect.TypeOf(v)))
}

func (v *Var) Eval() Expr {
	return nil
}

func (binary *Binary) Parse(scan *scanner.Scanner) {
	buf := bytes.Buffer{}
	scan.Scan()
	switch next, o := scan.Peek(), scan.TokenText(); o {
	case ">":
		if next == '=' {
			scan.Scan()
			binary.OP = ">="
		} else {
			binary.OP = ">"
		}
	case "<":
		if next == '=' {
			scan.Scan()
			binary.OP = "<="
		} else {
			binary.OP = "<"
		}
	case "=":
		if next == '=' {
			scan.Scan()
			binary.OP = "=="
		}
	case "!":
		if next == '=' {
			scan.Scan()
			binary.OP = "!="
		}
	case "&":
		if next == '&' {
			scan.Scan()
			binary.OP = "&&"
		}
	case "|":
		if next == '|' {
			scan.Scan()
			binary.OP = "||"
		}
	default:
		if !op[o] {
			panic(fmt.Sprintf("invalid method: %s", o))
		}
		binary.OP = o
	}

	for token := scan.Scan(); token != scanner.EOF; token = scan.Scan() {
		switch token {
		case 40: // (
			b := new(Binary)
			b.Parse(scan)
			binary.Args = append(binary.Args, b)
		case scanner.Float: // float
			if name := buf.String(); name != "" {
				buf.Reset()
				binary.Args = append(binary.Args, &Var{name, nil})
			}
			binary.Args = append(binary.Args, &Atom{token, scan.TokenText()})
		case scanner.Int: // int
			if name := buf.String(); name != "" {
				buf.Reset()
				binary.Args = append(binary.Args, &Var{name, nil})
			}
			binary.Args = append(binary.Args, &Atom{token, scan.TokenText()})
		case scanner.Ident: // var
			buf.WriteString(scan.TokenText())
		case 41: // )
			if name := buf.String(); name != "" {
				binary.Args = append(binary.Args, &Var{name, nil})
			}
			return
		case 46:
			buf.WriteByte('.')
		default:
			fmt.Println(token, scan.TokenText(), scanner.TokenString(token))
			panic("panic")
		}
	}
}
