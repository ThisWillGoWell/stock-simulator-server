package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ThisWillGoWell/stock-simulator-server/src/objects"

	"github.com/ThisWillGoWell/stock-simulator-server/src/id"

	"github.com/ThisWillGoWell/stock-simulator-server/src/id/change"

	"github.com/ThisWillGoWell/stock-simulator-server/src/wires/sender"

	"github.com/ThisWillGoWell/stock-simulator-server/src/wires"

	"unicode"

	"github.com/ThisWillGoWell/stock-simulator-server/src/lock"
	"github.com/ThisWillGoWell/stock-simulator-server/src/web/session"
	"github.com/ThisWillGoWell/stock-simulator-server/src/wires/duplicator"
)

// keep the uuid to user
var UserList = make(map[string]*User)

// keep the username to uuid list
var uuidList = make(map[string]string)
var UserListLock = lock.NewLock("user-list")

/*
User Object
Represents a unique individual of the system
*/
type User struct {
	objects.User
	Lock           *lock.Lock                    `json:"-"`
	UserUpdateChan *duplicator.ChannelDuplicator `json:"-"`
	Sender         *sender.Sender                `json:"-"`
}

/**
Return a user provided the username and Password
If the Password is correct return user, else return err
*/
func GetUser(username, password string) (*User, error) {
	UserListLock.Acquire("get-user")
	defer UserListLock.Release()
	userUuid, exists := uuidList[username]
	if !exists {
		return nil, errors.New("user does not exist")
	}
	user := UserList[userUuid]

	if !comparePasswords(user.Password, password) {
		return nil, errors.New("password is incorrect")
	}
	user.Active = true
	wires.UsersUpdate.Offer(user)
	return user, nil
}

func RenewUser(sessionToken string) (*User, error) {
	userId, err := session.GetUserId(sessionToken)
	if err != nil {
		return nil, err
	}
	UserListLock.Acquire("renew-user")
	defer UserListLock.Release()
	user, exists := UserList[userId]
	if !exists {
		return nil, errors.New("user found in session list but not in current users")
	}
	user.Active = true
	wires.UsersUpdate.Offer(user)
	return user, nil
}

func deleteUser(uuid string, lockAquired bool) {
	if !lockAquired {
		UserListLock.Acquire("delete-user")
		defer UserListLock.Release()
	}
	u, ok := UserList[uuid]
	if !ok {
		return
	}
	change.UnregisterChangeDetect(u)
	delete(UserList, uuid)
	delete(uuidList, u.DisplayName)
	u.Sender.Stop()

}


func MakeUser(uModel objects.User) (*User, error){
	_, userNameExists := uuidList[uModel.DisplayName]
	if userNameExists {
		return nil, errors.New("username already exists")
	}
	if uModel.Config == nil {
		uModel.Config = make(map[string]interface{})
		if uModel.ConfigStr != ""{
			if err := json.Unmarshal([]byte(uModel.ConfigStr), &uModel.Config); err != nil {
				return nil, fmt.Errorf("failed to unmarhsal provided config err=[%v]", err)
			}
		}
	}

	u := &User{
		User:           uModel,
		Lock:           lock.NewLock("user-" + uModel.Uuid),
		UserUpdateChan: duplicator.MakeDuplicator("user-" + uModel.Uuid),
		Sender:         sender.NewSender(uModel.Uuid, uModel.PortfolioId),
	}

	if err := change.RegisterPublicChangeDetect(u); err != nil {
		return nil, err
	}
	uuidList[uModel.DisplayName] = u.Uuid
	UserList[u.Uuid] = u
	id.RegisterUuid(u.Uuid, UserList[u.Uuid])
	return UserList[u.Uuid], nil
}

/**
Logout the user and decrement the active client count
*/
func (user *User) LogoutUser() {
	user.Lock.Acquire("logout")
	defer user.Lock.Release()
	user.ActiveClients -= 1
	if user.ActiveClients < 0 {
		user.ActiveClients = 0
	}
	if user.ActiveClients == 0 {
		user.Active = false
	}
	wires.UsersUpdate.Offer(user)
}

func (user *User) GetId() string {
	return user.Uuid
}
func (user *User) GetType() string {
	return "user"
}

/**
Turn the user map into a list so they can be sent to a rx client
*/
func GetAllUsers() []*User {
	UserListLock.Acquire("get all users")
	defer UserListLock.Release()
	lst := make([]*User, len(UserList))
	i := 0
	for _, val := range UserList {
		lst[i] = val
		i += 1
	}
	return lst
}

func (user *User) SetConfig(config map[string]interface{}) {
	user.Config = config
	configBytes, _ := json.Marshal(config)
	user.ConfigStr = string(configBytes)
	wires.UsersUpdate.Offer(user)
}

func (user *User) SetPassword(pass string) error {
	if len(pass) < minPasswordLength {
		return errors.New("password too short")
	}
	hashedPassword := hashAndSalt(pass)
	user.Password = hashedPassword
	wires.UsersUpdate.Offer(user)
	return nil
}

func (user *User) SetDisplayName(displayName string) error {
	if !isAllowedCharacterDisplayName(displayName) {
		return errors.New("display name contains invalid character")
	}
	if len(displayName) > maxDisplayNameLength {
		return errors.New("display name too long")
	}
	if len(displayName) < minDisplayNameLength {
		return errors.New("display name too short")
	}

	user.DisplayName = displayName
	wires.UsersUpdate.Offer(user)
	return nil
}

func isAllowedCharacterDisplayName(s string) bool {
	for _, r := range s {
		if !(unicode.IsLetter(r) || unicode.IsNumber(r) || r == '_' || r == '-') {
			return false
		}
	}
	return true
}

func isAllowedCharacterUsername(s string) bool {
	for _, r := range s {
		if !(unicode.IsLetter(r) || unicode.IsNumber(r)) {
			return false
		}
	}
	return true
}

func SendNotifcation(uuid string, note interface{}) {
	UserList[uuid].Sender.Notifications.Offer(note)
}
