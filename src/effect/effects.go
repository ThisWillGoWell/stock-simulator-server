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
var portfolioEffectTags = make(map[string]map[string]*Effect)

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
		switch e.Type.(type) {
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

func newEffect(portfolioUuid, title, effectType, tag string, innerEffect interface{}, duration time.Duration) {
	uuid := utils.SerialUuid()
	e := MakeEffect(uuid, portfolioUuid, title, tag, effectType, innerEffect, duration)
	wires.EffectsNewObject.Offer(e)
}

func MakeEffect(uuid, portfolioUuid, title, effectType, tag string, innerEffect interface{}, duration time.Duration) *Effect {
	newEffect := &Effect{
		PortfolioUuid: portfolioUuid,
		Uuid:          uuid,
		Title:         title,
		Active:        true,
		StartTime:     time.Now(),
		Duration:      utils.Duration{Duration: duration},
		Type:          effectType,
		InnerEffect:   innerEffect,
		Tag:           tag,
	}

	pEffects, ok := portfolioEffects[portfolioUuid]
	if !ok {
		pEffects = make(map[string]*Effect)
		portfolioEffects[portfolioUuid] = effects
	}
	if tag != "" {
		oldEffect, tagExists := portfolioEffectTags[portfolioUuid][tag]
		if tagExists {
			DeleteEffect(oldEffect.Uuid, true)
		}
		if _, portfolioExists := portfolioEffectTags[portfolioUuid]; !portfolioExists{
			portfolioEffectTags[portfolioUuid] =
		}
		portfolioEffectTags[portfolioUuid][tag] = newEffect
	}
	pEffects[newEffect.Uuid] = newEffect
	effects[newEffect.Uuid] = newEffect

	utils.RegisterUuid(uuid, newEffect)
	return newEffect
}

func UpdatePortfolioTag(portfolioUuid, tag string, newEffect *Effect) {
	EffectLock.Acquire("update portfolio effect tag")
	defer EffectLock.Release()
	tags, exists := portfolioEffectTags[portfolioUuid]
	if !exists {
		panic("got tag: " + tag + " update for a portfolio: " + portfolioUuid + " portfolio not found")
	}
	taggedEffect, foundTag := tags[tag]
	if !foundTag {
		panic("got tag: " + tag + " update for a portfolio: " + portfolioUuid + " tag not found")
	}
	DeleteEffect(taggedEffect.Uuid, true)
	portfolioEffects[portfolioUuid][newEffect.Uuid]

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
	Type          string         `json:"type"`
	InnerEffect   interface{}
	Tag           string `json:"tag"`
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
