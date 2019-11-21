package effect

import (
	"testing"
	"time"

	"github.com/ThisWillGoWell/stock-simulator-server/src/wires"

	"github.com/ThisWillGoWell/stock-simulator-server/src/utils"

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

func TestTempEffect(t *testing.T) {
	RunEffectCleaner()
	wires.ConnectWires()
	deletes := wires.EffectsDelete.GetBufferedOutput(1)
	NewTaxModifier("0", "test", time.Second*5, 0)
	for {
		select {
		case <-deletes:
			return
		default:

		}
	}

}
