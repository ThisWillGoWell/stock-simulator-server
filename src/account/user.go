package account

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/stock-simulator-server/src/duplicator"
	"github.com/stock-simulator-server/src/lock"
	"github.com/stock-simulator-server/src/portfolio"
	"github.com/stock-simulator-server/src/session"
	"github.com/stock-simulator-server/src/utils"
)

// keep the uuid to user
var userList = make(map[string]*User)

// keep the username to uuid list
var uuidList = make(map[string]string)
var userListLock = lock.NewLock("user-list")

var NewObjectChannel = duplicator.MakeDuplicator("New User")
var UpdateChannel = duplicator.MakeDuplicator("User Update")

const minPasswordLength = 4
const minDisplayNameLength = 4
const maxDisplayNameLength = 20

/*
User Object
Represents a unique individual of the system
*/
type User struct {
	UserName      string     `json:"-"`
	Password      string     `json:"-"`
	DisplayName   string     `json:"display_name" change:"-"`
	Uuid          string     `json:"-"`
	Active        bool       `json:"active" change:"-"`
	ActiveClients int64      `json:"-"`
	Lock          *lock.Lock `json:"-"`
	PortfolioId   string     `json:"portfolio_uuid"`
	Config 		  map[string]interface{}	 `json:"-"`
}

/**
Return a user provided the username and Password
If the Password is correct return user, else return err
*/
func GetUser(username, password string) (*User, error) {
	userListLock.Acquire("get-user")
	defer userListLock.Release()
	userUuid, exists := uuidList[username]
	if !exists {
		return nil, errors.New("user does not exist")
	}
	user := userList[userUuid]

	if !comparePasswords(user.Password, password) {
		return nil, errors.New("Password is incorrect")
	}
	user.Active = true
	UpdateChannel.Offer(user)
	return user, nil
}

func RenewUser(sessionToken string)(*User, error) {
	userId, err := session.GetUserId(sessionToken)
	if err != nil{
		return nil, err
	}
	userListLock.Acquire("renew-user")
	defer userListLock.Release()
	user, exists := userList[userId]
	if !exists{
		return nil, errors.New("user found in session list but not in current users")
	}
	user.Active = true
	UpdateChannel.Offer(user)
	return user, nil
}

/**
Build a new user
set their Password to that provided
*/
func NewUser(username, displayName,  password string) (*User, error) {
	uuid := utils.PseudoUuid()
	if len(username) > 20{
		return nil, errors.New("username too long")
	}
	if len(username) < 4{
		return nil, errors.New("username too short")
	}
	if len(displayName) > maxDisplayNameLength {
		return nil, errors.New("display name too long")
	}
	if len(displayName) < minDisplayNameLength{
		return nil, errors.New("display name too short")
	}
	if len(password) < minPasswordLength{
		return nil, errors.New("password too short")
	}
	hashedPassword := hashAndSalt(password)
	user, err := MakeUser(uuid, username, displayName, hashedPassword, "", `{"config":"this is a config", "nested":{"one":1}}`)
	if err != nil {
		utils.RemoveUuid(uuid)
		return nil, err
	}
	port, _ := portfolio.NewPortfolio(uuid)
	user.PortfolioId = port.UUID
	return user, nil
}

func MakeUser(uuid, username, displayName, password, portfolioUUID, config string) (*User, error) {
	userListLock.Acquire("new-user")
	defer userListLock.Release()
	_, userNameExists := uuidList[username]
	if userNameExists {
		return nil, errors.New("username already exists")
	}
	var configMap map[string]interface{}
	err := json.Unmarshal([]byte(config), &configMap)
	if err != nil{
		fmt.Println("error making config json in MakeUser: ", err)
		configMap = make(map[string]interface{})
	}
	uuidList[username] = uuid
	userList[uuid] = &User {
		UserName:    username,
		DisplayName: displayName,
		Password:    password,
		Uuid:        uuid,
		PortfolioId: portfolioUUID,
		Lock:        lock.NewLock("user"),
		Active:      true,
		Config:      configMap,
	}
	NewObjectChannel.Offer(userList[uuid])
	utils.RegisterUuid(uuid, userList[uuid])
	return userList[uuid], nil

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
	UpdateChannel.Offer(user)
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
	userListLock.Acquire("get all users")
	defer userListLock.Release()
	lst := make([]*User, len(userList))
	i := 0
	for _, val := range userList {
		lst[i] = val
		i += 1
	}
	return lst
}

func  (user *User) SetConfig(config map[string]interface{}){
	user.Config = config
	UpdateChannel.Offer(user)
}


func  (user *User) SetPassword(pass string) error{
	if len(pass) > minPasswordLength{
		return errors.New("password too short")
	}
	hashedPassword := hashAndSalt(pass)
	user.Password = hashedPassword
	UpdateChannel.Offer(user)
	return nil
}


func  (user *User) SetDisplayName(displayName string)error{
	if len(displayName) > maxDisplayNameLength {
		return errors.New("display name too long")
	}
	if len(displayName) < minDisplayNameLength {
		return errors.New("display name too short")
	}
	user.DisplayName = displayName
	UpdateChannel.Offer(user)
	return nil
}