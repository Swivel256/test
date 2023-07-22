package exchange

import (
	"github.com/ginvmbot/newstrade/pkg/fixedpoint"
	"github.com/ginvmbot/newstrade/pkg/types"
)

type Config struct {
	MinProfit  fixedpoint.Value
	MaxProfit  fixedpoint.Value
	StopLoss   fixedpoint.Value
	StepProfit fixedpoint.Value
	Market     types.Market
}
