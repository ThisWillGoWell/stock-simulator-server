package effect

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/ThisWillGoWell/stock-simulator-server/src/objects"
	"github.com/ThisWillGoWell/stock-simulator-server/src/wires/sender"

	"github.com/ThisWillGoWell/stock-simulator-server/src/id"

	"github.com/ThisWillGoWell/stock-simulator-server/src/app/log"

	"github.com/ThisWillGoWell/stock-simulator-server/src/database"

	"github.com/ThisWillGoWell/stock-simulator-server/src/objects/notification"

	"github.com/ThisWillGoWell/stock-simulator-server/src/id/change"

	"github.com/ThisWillGoWell/stock-simulator-server/src/merge"

	"github.com/ThisWillGoWell/stock-simulator-server/src/wires"

	"github.com/ThisWillGoWell/stock-simulator-server/src/lock"
	"github.com/ThisWillGoWell/stock-simulator-server/src/utils"
)

var EffectLock = lock.NewLock("EffectLock")
var effects = make(map[string]*Effect)
var portfolioEffects = make(map[string]map[string]*Effect)
var portfolioEffectTags = make(map[string]map[string]*Effect)

func DeleteRequest(uuid string) error {
	EffectLock.Acquire("delete-effect")
	defer EffectLock.Release()
	e, ok := effects[uuid]
	if !ok {
		log.Log.Warnf("got effect delete for effect that does not exists")
		return nil
	}

	if err := database.Db.Execute(nil, []interface{}{e}); err != nil {
		return err
	}
	wires.EffectsDelete.Offer(e.Effect)
	return nil
}

func DeleteEffect(e *Effect) {

	delete(portfolioEffectTags[e.PortfolioUuid], e.Tag)
	delete(effects, e.Uuid)
	delete(portfolioEffects[e.PortfolioUuid], e.Uuid)
	if len(portfolioEffects[e.PortfolioUuid]) == 0 {
		delete(portfolioEffects, e.PortfolioUuid)
		delete(portfolioEffectTags, e.PortfolioUuid)
	}
	change.UnregisterChangeDetect(e.Effect)
	id.RemoveUuid(e.Uuid)

}

func getTaggedEffect(PortfolioUuid, tag string) *Effect {
	if tag == "" {
		return nil
	}
	if _, ok := portfolioEffectTags[PortfolioUuid]; !ok {
		return nil
	}
	if e, ok := portfolioEffectTags[PortfolioUuid][tag]; !ok {
		return nil
	} else {
		return e
	}
}

func newEffect(PortfolioUuid, title, effectType, tag string, innerEffect interface{}, duration time.Duration) (*Effect, *Effect, error) {
	EffectLock.EnableDebug()
	EffectLock.Acquire("make-effect")
	defer EffectLock.Release()

	// tagged effect
	preEffect := getTaggedEffect(PortfolioUuid, tag)
	effect := objects.Effect{
		Uuid:          id.SerialUuid(),
		PortfolioUuid: PortfolioUuid,
		Title:         title,
		Type:          effectType,
		Tag:           tag,
		InnerEffect:   innerEffect,
		Duration:      utils.Duration{Duration: duration},
		StartTime:     time.Now(),
	}
	var e *Effect
	var err error
	if e, err = MakeEffect(effect, true); err != nil {
		return nil, nil, err
	}
	return e, preEffect, nil
}

func MakeEffect(effect objects.Effect, lockAcquired bool) (*Effect, error) {
	if !lockAcquired {
		EffectLock.Acquire("make-effect")
		defer EffectLock.Release()
	}

	if s, ok := effect.InnerEffect.(string); ok {
		var err error
		if effect.InnerEffect, err = UnmarshalJsonEffect(effect.Type, s); err != nil {
			return nil, fmt.Errorf("failed to unmarhsal inner effect err=[%v]", err)
		}
	}
	newEffect := &Effect{
		Effect: effect,
	}

	if err := change.RegisterPublicChangeDetect(newEffect.Effect); err != nil {
		return nil, err
	}

	pEffects, ok := portfolioEffects[effect.PortfolioUuid]
	if !ok {
		pEffects = make(map[string]*Effect)
		portfolioEffects[effect.PortfolioUuid] = pEffects
	}
	if effect.Tag != "" {
		if _, portfolioExists := portfolioEffectTags[effect.PortfolioUuid]; !portfolioExists {
			// the portfolio map was deleted by the Delete Effect
			pEffects = make(map[string]*Effect)
			portfolioEffects[effect.PortfolioUuid] = pEffects
			portfolioEffectTags[effect.PortfolioUuid] = make(map[string]*Effect)
		}
		portfolioEffectTags[effect.PortfolioUuid][effect.Tag] = newEffect
	}

	pEffects[newEffect.Uuid] = newEffect
	effects[newEffect.Uuid] = newEffect
	id.RegisterUuid(newEffect.Uuid, newEffect)
	return newEffect, nil
}

//func UpdatePortfolioTag(PortfolioUuid, tag string, newEffect *Effect) {
//	EffectLock.Acquire("update portfolio effect tag")
//	defer EffectLock.Release()
//	tags, exists := portfolioEffectTags[PortfolioUuid]
//	if !exists {
//		panic("got tag: " + tag + " update for a portfolio: " + PortfolioUuid + " portfolio not found")
//	}
//	taggedEffect, foundTag := tags[tag]
//	if !foundTag {
//		panic("got tag: " + tag + " update for a portfolio: " + PortfolioUuid + " tag not found")
//	}
//	DeleteEffect(taggedEffect.Uuid, true)
//	portfolioEffects[PortfolioUuid][newEffect.Uuid]
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
	objects.Effect
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
					if err := DeleteRequest(uuid); err != nil {
						log.Log.Errorf("failed to clean effect err=[%v]", err)
					} else {
						notification.NotificationLock.Acquire("clean-effects")
						n := notification.EndEffectNotification(effect.PortfolioUuid, effect.Title)
						if dbErr := database.Db.Execute([]interface{}{n}, nil); dbErr != nil {
							log.Log.Errorf("failed to send end effect notification to %s id=%s err=[%v]", effect.PortfolioUuid, effect.Uuid, err)
							notification.DeleteNotification(n)
						} else {
							sender.SendNewObject(n.PortfolioUuid, n.Notification)
						}

						notification.NotificationLock.Release()

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
