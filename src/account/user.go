package account

import (
	"errors"
	"github.com/stock-simulator-server/src/portfolio"
	"github.com/stock-simulator-server/src/utils"
)

// keep the uuid to user
var userList = make(map[string]*User)

// keep the username to uuid list
var uuidList = make(map[string]string)
var userListLock = utils.NewLock("user-list")

type User struct {
	UserName      string
	password      string
	DisplayName   string `json:"display_name" change:"-"`
	Uuid          string
	Active        bool   `json:"active" change:"-"`
	ActiveClients int64
	Lock 		  *utils.Lock
}

func GetUser(username, password string) (*User, error) {
	userListLock.Acquire("get-user")
	defer userListLock.Release()
	userUuid, exists := uuidList[username]
	if !exists {
		return nil, errors.New("user does not exist")
	}
	user := userList[userUuid]
	if user.password != password {
		return nil, errors.New("password is incorrect")
	}
	user.Active = true
	utils.SubscribeUpdateInputs.Offer(user)
	return user, nil
}

func NewUser(username, password string) (*User, error){
	userListLock.Acquire("new-user")
	defer userListLock.Release()
	_, userNameExists := uuidList[username]
	if userNameExists{
		return nil, errors.New("username already exists")
	}
	uuid := utils.PseudoUuid()
	for {
		// keep going util a unique uuid is found.. should really never have to retry
		_, exists := userList[uuid]
		if !exists {
			uuidList[username] = uuid
			portfolio.NewPortfolio(uuid, username)
			userList[uuid] = &User{
				UserName:    username,
				DisplayName: username,
				password:    password,
				Uuid:        uuid,
				Lock: 		 utils.NewLock("user"),
				Active: 	 true,
			}
			utils.SubscribeUpdateInputs.Offer(userList[uuid])
			return userList[uuid], nil
		}
		uuid = utils.PseudoUuid()
	}
	panic("exited new user without making a new user")

	return nil, nil
}

func (user *User) LogoutUser(){
	user.Lock.Acquire("logout")
	defer user.Lock.Release()
	user.ActiveClients -= 1
	if user.ActiveClients < 0{
		user.ActiveClients = 0
	}
	if user.ActiveClients == 0{
		user.Active = false
	}
	utils.SubscribeUpdateInputs.Offer(user)
}

func GetAllUsers() map[string]string{

	return uuidList
}
func (user *User)GetId() string {
	return user.Uuid
}
func (user *User) GetType() string {
	return "user"
}

