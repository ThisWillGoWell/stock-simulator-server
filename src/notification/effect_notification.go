package notification

const NewEffectNotificationType = "new_effect"
const EndEffectNotificationType = "end_effect"

type EffectNotification struct {
	EffectTitle string `json:"title"`
}

func NewEffectNotification(portfolioId, effectTitle string) error {
	_, err := NewNotification(portfolioId, NewEffectNotificationType, &EffectNotification{effectTitle})
	return err
}
func EndEffectNotification(portfolioId, effectTitle string) error {
	_, err := NewNotification(portfolioId, EndEffectNotificationType, &EffectNotification{effectTitle})
	return err
}
