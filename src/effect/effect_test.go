package effect

import (
	"testing"

	"github.com/stock-simulator-server/src/utils"

	"github.com/gotestyourself/gotestyourself/assert"
)

func TestMakeEffect(t *testing.T) {
	NewBaseTradeEffect("portfolio")
	e, _ := TotalTradeEffect("portfolio")
	assert.Assert(t, *e.SellFeeAmount == BaseSellFell)
	utils.PrintJson(e)
	printAllEffects()
	UpdateBaseProfit("portfolio", 1.3)
	printAllEffects()
}

func printAllEffects() {
	for _, e := range effects {
		utils.PrintJson(e)
	}
}
