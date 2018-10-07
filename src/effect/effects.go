package effect

import (
	"time"
)

type EffectType interface {
	Name() string
	Cost() int64
	RequiredLevel() int64
}

type Effect interface {
	Active() bool
	Duration() time.Duration
	StartTime() time.Time
	Type() EffectType
	Start()
	PortfolioUuid() string
	GetUuid() string
}
