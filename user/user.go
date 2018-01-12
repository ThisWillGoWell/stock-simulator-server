package user

import (
	"stock-server/utils"
	"errors"
	"stock-server/exchange"
)

// keep the uuid to user
var userList = make(map[string]*User)
// keep the username to uuid list
var uuidList = make(map[string]string)
var userListLock = utils.NewLock()

type User struct {
	username string
	password string
	displayName string
	uuid string

}

func getUser(username, password string) (*User, error) {
	userListLock.Acquire()
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
	userListLock.Acquire()
	defer userListLock.Release()

	uuid := utils.PseudoUuid()
	for {
		// keep going util a unique uuid is found.. should really never have to retry
		_, exists := userList[uuid]
		if ! exists{
			uuidList[username] = uuid
			userList[uuid]=&User{
				username: username,
				displayName: username,
				password:password,
				uuid: uuid,
			}

			return userList[uuid]
		}
		uuid = utils.PseudoUuid()
	}
	return nil
}