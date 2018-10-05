package notification

const NewItemNotificationType = "new_item"
const UsedItemNotificationType = "used_item"

type ItemNotification struct {
	ItemType string `json:"item_type"`
	ItemUuid string `json:"item_uuid"`
}

func NewItemNotification(userUuid, itemType, itemUuid string) {
	NewNotification(userUuid, NewItemNotificationType,
		&ItemNotification{
			ItemType: itemType,
			ItemUuid: itemUuid,
		})

}

func UsedItemNotification(userUuid, itemUuid, itemType string) {
	NewNotification(userUuid, UsedItemNotificationType,
		&ItemNotification{
			ItemType: itemType,
			ItemUuid: itemUuid,
		})
}
