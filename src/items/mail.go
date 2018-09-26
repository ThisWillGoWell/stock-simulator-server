package items

import "github.com/stock-simulator-server/src/order"

type MailItemType struct {

}

type MailItemParameters struct{
	To string `json:"to"`
	Text string `json:"string"`
	Amount int64 `json:"amount"`
}

func (MailItemType) GetName() string{
	return "Mail"
}

func (MailItemType) GetCost() int64{
	return 10
}

func (MailItemType) GetDescription() string{
	return "Send money and a message to another user"
}

func  (MailItemType) RequiredLevel() int64{
	return 1
}

func (MailItemType) GetActivateParameters() interface{}{
	return MailItemParameters{}
}

type MailItem struct {
	Type MailItemType
	UserUuid string
	Uuid string
	Used bool
}

func (it *MailItem) GetType() ItemType{
	return it.Type
}
func  (it *MailItem) GetUserUuid() string {
	return it.UserUuid
}
func  (it *MailItem) GetUuid() string {
	return it.Uuid
}

func (it *MailItem) HasBeenUsed() bool{
	return it.Used
}

func (it *MailItem) RequiredLevel() int64{
	return 1
}
func (it *MailItem) Activate(){
	order.MakeTransferOrder()
}