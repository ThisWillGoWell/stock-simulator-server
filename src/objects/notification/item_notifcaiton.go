package notification

const NewItemNotificationType = "new_item"
const UsedItemNotificationType = "used_item"

type ItemNotification struct {
	ItemType string `json:"item_type"`
	ItemUuid string `json:"item_uuid"`
}

func NewItemNotification(portfolioUuid, itemType, itemUuid string) *Notification {
	return NewNotification(portfolioUuid, NewItemNotificationType,
		&ItemNotification{
			ItemType: itemType,
			ItemUuid: itemUuid,
		})

}

func UsedItemNotification(portfolioUuid, itemUuid, itemType string) *Notification {
	return NewNotification(portfolioUuid, UsedItemNotificationType,
		&ItemNotification{
			ItemType: itemType,
			ItemUuid: itemUuid,
		})
}
