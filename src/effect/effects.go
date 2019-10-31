package effect

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/ThisWillGoWell/stock-simulator-server/src/log"

	"github.com/ThisWillGoWell/stock-simulator-server/src/database"

	"github.com/sirupsen/logrus"

	"github.com/ThisWillGoWell/stock-simulator-server/src/notification"

	"github.com/ThisWillGoWell/stock-simulator-server/src/change"

	"github.com/ThisWillGoWell/stock-simulator-server/src/merge"

	"github.com/ThisWillGoWell/stock-simulator-server/src/wires"

	"github.com/ThisWillGoWell/stock-simulator-server/src/lock"
	"github.com/ThisWillGoWell/stock-simulator-server/src/utils"
)

var EffectLock = lock.NewLock("EffectLock")
var effects = make(map[string]*Effect)
var portfolioEffects = make(map[string]map[string]*Effect)
var portfolioEffectTags = make(map[string]map[string]*Effect)

const EffectIdType = "effect"

func deleteEffect(uuid string, lockAcquired bool) error {
	if !lockAcquired {
		EffectLock.Acquire("delete-effect")
		defer EffectLock.Release()
	}
	if _, ok := effects[uuid]; !ok {
		log.Log.Warnf("got effect delete for effect that does not exists")
		return nil
	}
	e := effects[uuid]
	wires.EffectsDelete.Offer(e)
	dbErr := database.Db.DeleteEffect(e)
	delete(portfolioEffectTags[e.PortfolioUuid], e.Tag)
	delete(effects, uuid)
	delete(portfolioEffects[e.PortfolioUuid], uuid)
	if len(portfolioEffects[e.PortfolioUuid]) == 0 {
		delete(portfolioEffects, e.PortfolioUuid)
		delete(portfolioEffectTags, e.PortfolioUuid)
	}
	change.UnregisterChangeDetect(e)
	notification.EndEffectNotification(e.PortfolioUuid, e.Title)
	utils.RemoveUuid(uuid)

	return dbErr
}

func newEffect(portfolioUuid, title, effectType, tag string, innerEffect interface{}, duration time.Duration) (e *Effect, err error) {
	uuid := utils.SerialUuid()
	if e, err = MakeEffect(uuid, portfolioUuid, title, effectType, tag, innerEffect, duration, time.Now()); err != nil {
		logrus.Errorf("making effect err=%v", err)
		return nil, fmt.Errorf("failed to make effect")
	}

	if err := database.Db.WriteEffect(e); err != nil {
		_ = deleteEffect(uuid, false)
		return nil, fmt.Errorf("failed to make effect")
	}

	wires.EffectsNewObject.Offer(e)
	notification.NewEffectNotification(portfolioUuid, title)
	return e, nil
}

func getTaggedEffect(portfolioUuid, tag string) *Effect {
	pTagEffects, portfolioExists := portfolioEffectTags[portfolioUuid]
	if !portfolioExists {
		panic("portfolio not fond in update base: " + portfolioUuid)
	}
	effect, tagExists := pTagEffects[tag]
	if !tagExists {
		panic("portfolio: " + portfolioUuid + " does not have a base tag to update")
	}
	return effect
}

func MakeEffect(uuid, portfolioUuid, title, effectType, tag string, innerEffect interface{}, duration time.Duration, startTime time.Time) (*Effect, error) {
	EffectLock.Acquire("make-effect")
	defer EffectLock.Release()
	newEffect := &Effect{
		PortfolioUuid: portfolioUuid,
		Uuid:          uuid,
		Title:         title,
		StartTime:     startTime,
		Duration:      utils.Duration{Duration: duration},
		Type:          effectType,
		InnerEffect:   innerEffect,
		Tag:           tag,
	}

	if err := change.RegisterPublicChangeDetect(newEffect); err != nil {
		return nil, err
	}

	pEffects, ok := portfolioEffects[portfolioUuid]
	if !ok {
		pEffects = make(map[string]*Effect)
		portfolioEffects[portfolioUuid] = pEffects
	}
	if tag != "" {
		oldEffect, tagExists := portfolioEffectTags[portfolioUuid][tag]
		if tagExists {
			_ = deleteEffect(oldEffect.Uuid, true)
		}
		if _, portfolioExists := portfolioEffectTags[portfolioUuid]; !portfolioExists {
			// the portfolio map was deleted by the Delete Effect
			pEffects = make(map[string]*Effect)
			portfolioEffects[portfolioUuid] = pEffects
			portfolioEffectTags[portfolioUuid] = make(map[string]*Effect)
		}
		portfolioEffectTags[portfolioUuid][tag] = newEffect
	}

	pEffects[newEffect.Uuid] = newEffect
	effects[newEffect.Uuid] = newEffect
	utils.RegisterUuid(uuid, newEffect)
	return newEffect, nil
}

//func UpdatePortfolioTag(portfolioUuid, tag string, newEffect *Effect) {
//	EffectLock.Acquire("update portfolio effect tag")
//	defer EffectLock.Release()
//	tags, exists := portfolioEffectTags[portfolioUuid]
//	if !exists {
//		panic("got tag: " + tag + " update for a portfolio: " + portfolioUuid + " portfolio not found")
//	}
//	taggedEffect, foundTag := tags[tag]
//	if !foundTag {
//		panic("got tag: " + tag + " update for a portfolio: " + portfolioUuid + " tag not found")
//	}
//	DeleteEffect(taggedEffect.Uuid, true)
//	portfolioEffects[portfolioUuid][newEffect.Uuid]
//
//}

type EffectType interface {
	Name() string
}

// ticket charge on a stock bought
// $5 per trade both sides of sale
// used to be % of trade
// fee on money managed
//
//

type Effect struct {
	PortfolioUuid string         `json:"portfolio_uuid"`
	Uuid          string         `json:"uuid"`
	Title         string         `json:"title" change:"-"`
	Duration      utils.Duration `json:"duration"`
	StartTime     time.Time      `json:"time"`
	Type          string         `json:"type"`
	InnerEffect   interface{}    `json:"-" change:"inner"`
	Tag           string         `json:"tag"`
}

type e2 struct {
	PortfolioUuid string         `json:"portfolio_uuid"`
	Uuid          string         `json:"uuid"`
	Title         string         `json:"title"`
	Duration      utils.Duration `json:"duration"`
	StartTime     time.Time      `json:"time"`
	Type          string         `json:"type"`
	Tag           string         `json:"tag"`
}

func (u *Effect) MarshalJSON() ([]byte, error) {

	return merge.Json(e2{
		PortfolioUuid: u.PortfolioUuid,
		Uuid:          u.Uuid,
		Title:         u.Title,
		Duration:      u.Duration,
		StartTime:     u.StartTime,
		Type:          u.Type,
		Tag:           u.Tag,
	}, u.InnerEffect)
}

func RunEffectCleaner() {
	go func() {
		for range time.Tick(time.Second) {
			EffectLock.Acquire("clean")
			for uuid, effect := range effects {
				if effect.Duration.Duration != 0 && time.Since(effect.StartTime) > effect.Duration.Duration {
					if err := deleteEffect(uuid, true); err != nil {
						log.Log.Errorf("failed to clean effect err=%v", err)
					}
				}
			}
			EffectLock.Release()
		}
	}()
}

func GetAllEffects() []*Effect {
	EffectLock.Acquire("get-all")
	defer EffectLock.Release()
	effectsSlice := make([]*Effect, len(effects))
	i := 0
	for _, e := range effects {
		effectsSlice[i] = e
		i++
	}
	return effectsSlice
}
func (*Effect) GetType() string {
	return EffectIdType
}

func (e *Effect) GetId() string {
	return e.Uuid
}
func UnmarshalJsonEffect(effectType, jsonStr string) (interface{}, error) {
	var innerEffect interface{}
	switch effectType {
	case TradeEffectType:
		innerEffect = &TradeEffect{}
	}
	err := json.Unmarshal([]byte(jsonStr), &innerEffect)
	if err != nil {
		return nil, err
	}
	return innerEffect, nil
}
