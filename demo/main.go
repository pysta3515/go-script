package main

import (
	"fmt"

	go_script "github.com/pysta3515/go-script"
)

type Sale struct {
	SaleID    int64
	SalePrice float64
	Shop      *Shop
}

func (sale *Sale) Env() map[string]interface{} {
	return go_script.Env(sale, "sale")
}

type Shop struct {
	ShopID int
}

func main() {
	id := 10
	env := go_script.Env(id, "id")
	e := "(+  (- id 1) 2 (+ 3 4))"
	fmt.Println(go_script.Eval(e).ToInt(env))
	sale := Sale{
		1,
		2.2,
		&Shop{
			3,
		},
	}
	e = "(if (>= sale.SaleID 1) 10.1 0)"
	fmt.Println(go_script.Eval(e).ToFloat(sale.Env()))
}
