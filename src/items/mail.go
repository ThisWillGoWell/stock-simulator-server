package items

import (
	"github.com/pkg/errors"
	"github.com/stock-simulator-server/src/account"
	"github.com/stock-simulator-server/src/order"
	"github.com/stock-simulator-server/src/utils"
)

const mailItemType = "mail"

type MailItemType struct {
}

type MailItemParameters struct {
	To     string `json:"to"`
	Text   string `json:"string"`
	Amount int64  `json:"amount"`
}

func (MailItemType) GetName() string {
	return "Mail"
}

func (MailItemType) GetType() string {
	return mailItemType
}

func (MailItemType) GetCost() int64 {
	return 10
}

func (MailItemType) GetDescription() string {
	return "Send money and a message to another user"
}

func (MailItemType) RequiredLevel() int64 {
	return 1
}

func (MailItemType) GetActivateParameters() interface{} {
	return MailItemParameters{}
}

type MailItem struct {
	Type       MailItemType       `json:"type"`
	UserUuid   string             `json:"user_uuid"`
	Uuid       string             `json:"uuid"`
	Used       bool               `json:"used"`
	Result     MailItemParameters `json:"result,omitempty"`
	UpdateChan chan interface{}   `json:"-"`
}

func newMailItem(userUuid string) *MailItem {
	uuid := utils.SerialUuid()
	item := &MailItem{
		UserUuid: userUuid,
		Uuid:     uuid,
		Used:     false,
	}
	utils.RegisterUuid(uuid, item)
	return item
}
func (it *MailItem) GetType() string {
	return ItemIdentifiableType
}
func (it *MailItem) GetId() string {
	return it.Uuid
}

func (it *MailItem) GetItemType() ItemType {
	return it.Type
}
func (it *MailItem) GetUserUuid() string {
	return it.UserUuid
}
func (it *MailItem) GetUuid() string {
	return it.Uuid
}

func (it *MailItem) HasBeenUsed() bool {
	return it.Used
}

func (it *MailItem) SetUserUuid(uuid string) {
	it.UserUuid = uuid
}

func (it *MailItem) RequiredLevel() int64 {
	return 1
}

func (it *MailItem) View() interface{} {
	return it.Result
}

func (it *MailItem) GetUpdateChan() chan interface{} {
	return it.UpdateChan
}

func (it *MailItem) Activate(parameters interface{}) (interface{}, error) {
	mailParams, ok := parameters.(MailItemParameters)
	if !ok {
		return nil, errors.New("incorrect parameters type, how?")
	}
	senderId := account.UserList[it.UserUuid].PortfolioId
	receiver, ok := account.UserList[mailParams.To]
	if !ok {
		return nil, errors.New("giver user id not found")
	}

	to := order.MakeTransferOrder(senderId, receiver.PortfolioId, mailParams.Amount)
	result := <-to.ResponseChannel
	if !result.Success {
		return nil, errors.New(result.Err)
	}
	it.Used = true
	TransferItem(it.UserUuid, mailParams.To, it.Uuid)
	return nil, nil
}
