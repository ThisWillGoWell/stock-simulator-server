package notification

const NewEffectNotificationType = "new_effect"
const EndEffectNotificationType = "end_effect"

type EffectNotification struct {
	EffectTitle string `json:"title"`
}

func NewEffectNotification(portfolioId, effectTitle string) *Notification {
	return NewNotification(portfolioId, NewEffectNotificationType, &EffectNotification{effectTitle})
}
func EndEffectNotification(portfolioId, effectTitle string) *Notification {
	return NewNotification(portfolioId, EndEffectNotificationType, &EffectNotification{effectTitle})
}
