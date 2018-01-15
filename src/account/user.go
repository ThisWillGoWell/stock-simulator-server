package account

import (
	"github.com/stock-simulator-server/src/utils"
	"errors"
	"github.com/stock-simulator-server/src/portfolio"
)

// keep the uuid to user
var userList = make(map[string]*User)
// keep the username to uuid list
var uuidList = make(map[string]string)
var userListLock = utils.NewLock("user-list")

type User struct {
	Username string
	password string
	DisplayName string
	Uuid string
	ActiveClients int64
}

func GetUser(username, password string) (*User, error) {
	userListLock.Acquire("get-user")
	defer userListLock.Release()
	userUuid, exists :=  uuidList[username]
	if ! exists{
		return nil, errors.New("user does not exist")
	}
	user := userList[userUuid]
	if user.password != password {
		return nil, errors.New("password is incorrect")
	}
	return user, nil
}

func NewUser(username, password string) *User{
	userListLock.Acquire("new-user")
	defer userListLock.Release()

	uuid := utils.PseudoUuid()
	for {
		// keep going util a unique uuid is found.. should really never have to retry
		_, exists := userList[uuid]
		if ! exists{
			uuidList[username] = uuid
			portfolio.NewPortfolio(uuid, username)
			userList[uuid]=&User{
				Username: username,
				DisplayName: username,
				password:password,
				Uuid: uuid,
			}

			return userList[uuid]
		}
		uuid = utils.PseudoUuid()
	}
	return nil
}