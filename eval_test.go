package go_script

import (
	"fmt"
	"testing"
)

type Sale struct {
	SaleID    int
	SalePrice float64
	Shop      *Shop
}

func (sale *Sale) Env() map[string]interface{} {
	return map[string]interface{}{
		"sale.SaleID":      sale.SaleID,
		"sale.SalePrice":   sale.SalePrice,
		"sale.Shop.ShopID": sale.Shop.ShopID,
	}
}

type Shop struct {
	ShopID int
}

func TestEval(t *testing.T) {
	sale := Sale{
		1,
		2.2,
		&Shop{
			3,
		},
	}
	var e = "(if (== sale.Shop.ShopID 2.3) 5 0"
	//
	fmt.Println(Eval(e).ToInt(sale.Env()))
}

func BenchmarkEval(b *testing.B) {
	sale := Sale{
		1,
		2.2,
		&Shop{
			3,
		},
	}
	env := sale.Env()
	rules := []string{
		"(if (> sale.SaleID 3) 5 0",
		"(if (> sale.Shop.ShopID 3) 5 0",
	}
	scripts := make([]Expr, 0, len(rules))
	for _, rule := range rules {
		scripts = append(scripts, Eval(rule))
	}
	for i := 0; i < b.N; i++ {
		for _, script := range scripts {
			script.ToInt(env)
		}
	}
}
