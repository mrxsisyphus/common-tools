package bigdecimal_helper

import "github.com/shopspring/decimal"

var (
	//ZeroToTen 0-10
	ZeroToTen = []decimal.Decimal{
		decimal.NewFromInt(0),
		decimal.NewFromInt(1),
		decimal.NewFromInt(2),
		decimal.NewFromInt(3),
		decimal.NewFromInt(4),
		decimal.NewFromInt(5),
		decimal.NewFromInt(6),
		decimal.NewFromInt(7),
		decimal.NewFromInt(8),
		decimal.NewFromInt(9),
		decimal.NewFromInt(10),
	}

	Zero = ZeroToTen[0]
	One  = ZeroToTen[1]
	Two  = ZeroToTen[2]
	Ten  = ZeroToTen[10]
)
