package go_script

import (
	"fmt"
	"reflect"
	"strings"
	"text/scanner"
)

func Eval(expr string) Expr {
	scan := new(scanner.Scanner)
	scan.Init(strings.NewReader(expr))
	scan.Mode = scanner.ScanFloats | scanner.ScanChars | scanner.ScanIdents
	scan.Scan()
	binary := new(Binary)
	binary.Parse(scan)
	return binary
}

func Env(i interface{}, receiver string) map[string]interface{} {
	val := reflect.ValueOf(i)
	m := make(map[string]interface{})
	switch val.Kind() {
	case reflect.Ptr:
		env(m, receiver, val.Elem())
	default:
		env(m, receiver, val)
	}
	return m
}

func env(m map[string]interface{}, prefix string, value reflect.Value) {
	typ := value.Type()
	switch typ.Kind() {
	case reflect.Int:
		m[prefix] = int(value.Int())
	case reflect.Int64:
		m[prefix] = value.Int()
	case reflect.Float64:
		m[prefix] = value.Float()
	case reflect.Struct:
		n := value.NumField()
		for i := 0; i < n; i++ {
			val := value.Field(i)
			field := typ.Field(i)
			switch val.Kind() {
			case reflect.Int:
				m[fmt.Sprintf("%s.%s", prefix, field.Name)] = int(val.Int())
			case reflect.Int64:
				m[fmt.Sprintf("%s.%s", prefix, field.Name)] = val.Int()
			case reflect.Float64:
				m[fmt.Sprintf("%s.%s", prefix, field.Name)] = val.Float()
			case reflect.Struct:
				env(m, fmt.Sprintf("%s.%s", prefix, field.Name), val)
			case reflect.Ptr:
				env(m, fmt.Sprintf("%s.%s", prefix, field.Name), val.Elem())
			}
		}
	}
}
