package effect

import (
	"time"

	"github.com/stock-simulator-server/src/wires"

	"github.com/stock-simulator-server/src/lock"
	"github.com/stock-simulator-server/src/utils"
)

var EffectLock = lock.NewLock("EffectLock")
var effects = make(map[string]*Effect)
var portfolioEffects = make(map[string]map[string]*Effect)

// Calculate the total bonus  for a portfolio
func TotalTradeEffect(portfolioUuid string) (TradeEffect, []string) {
	EffectLock.Acquire("TotalBonus")
	effect, ok := portfolioEffects[portfolioUuid]
	totalEffect := TradeEffect{}
	uuids := make([]string, 0)
	if !ok {
		return totalEffect, uuids
	}
	for uuid, e := range effect {
		switch e.InnerEffect.(type) {
		case TradingType:
			totalEffect.Add(e.InnerEffect.(*TradeEffect))
			uuids = append(uuids, uuid)
		}
	}
	return totalEffect, uuids
}

func DeleteEffect(uuid string, lockAcquired bool) {
	if lockAcquired {
		EffectLock.Acquire("delete effect")
		defer EffectLock.Release()
	}

	if _, ok := effects[uuid]; !ok {
		panic("got delete for effect Not found" + uuid)
	}
	e := effects[uuid]
	wires.EffectsDelete.Offer(e)
	delete(effects, uuid)
	delete(portfolioEffects[e.PortfolioUuid], uuid)
	if len(portfolioEffects[e.PortfolioUuid]) == 0 {
		delete(portfolioEffects, e.PortfolioUuid)
	}
	utils.RemoveUuid(uuid)
}

func newEffect(portfolioUuid, title string, innerEffect interface{}, duration time.Duration) {
	uuid := utils.SerialUuid()
	e := makeEffect(uuid, portfolioUuid, title, innerEffect, duration)
	wires.EffectsNewObject.Offer(e)
}

func makeEffect(uuid, portfolioUuid, title string, innerEffect interface{}, duration time.Duration) *Effect {
	newEffect := &Effect{
		PortfolioUuid: portfolioUuid,
		Uuid:          uuid,
		Title:         title,
		Active:        true,
		StartTime:     time.Now(),
		Duration:      utils.Duration{Duration: duration},
		Type:          TradingType{},
		InnerEffect:   innerEffect,
	}

	peffects, ok := portfolioEffects[portfolioUuid]
	if !ok {
		peffects = make(map[string]*Effect)
		portfolioEffects[portfolioUuid] = effects
	}
	peffects[newEffect.Uuid] = newEffect
	effects[newEffect.Uuid] = newEffect
	utils.RegisterUuid(uuid, newEffect)
	return newEffect
}

type EffectType interface {
	Name() string
}

type Effect struct {
	PortfolioUuid string         `json:"portfolio_uuid"`
	Uuid          string         `json:"uuid"`
	Title         string         `json:"title"`
	Active        bool           `json:"active"`
	IsPublic      bool           `json:"public"`
	Duration      utils.Duration `json:"duration"`
	StartTime     time.Time      `json:"time"`
	Type          EffectType     `json:"effect"`
	InnerEffect   interface{}
}

func RunEffectCleaner() {
	for range time.Tick(time.Second) {
		EffectLock.Acquire("clean")
		for uuid, effect := range effects {
			if time.Since(effect.StartTime) > effect.Duration.Duration {
				DeleteEffect(uuid, true)
			}
		}
		EffectLock.Release()
	}
}
