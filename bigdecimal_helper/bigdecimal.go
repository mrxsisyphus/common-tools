package bigdecimal_helper

import "github.com/shopspring/decimal"

var (
	//ZeroToTen 0-10
	ZeroToTen []*decimal.Decimal

	Zero = ZeroToTen[0]
	One  = ZeroToTen[1]
	Two  = ZeroToTen[2]
	Ten  = ZeroToTen[10]
)

func init() {
	d0 := decimal.NewFromInt(0)
	d1 := decimal.NewFromInt(1)
	d2 := decimal.NewFromInt(2)
	d3 := decimal.NewFromInt(3)
	d4 := decimal.NewFromInt(4)
	d5 := decimal.NewFromInt(5)
	d6 := decimal.NewFromInt(6)
	d7 := decimal.NewFromInt(7)
	d8 := decimal.NewFromInt(8)
	d9 := decimal.NewFromInt(9)
	d10 := decimal.NewFromInt(10)
	ZeroToTen = []*decimal.Decimal{
		&d0,
		&d1,
		&d2,
		&d3,
		&d4,
		&d5,
		&d6,
		&d7,
		&d8,
		&d9,
		&d10,
	}
}
