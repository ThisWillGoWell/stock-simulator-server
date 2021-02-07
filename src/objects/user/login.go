package user

import (
	"errors"
	"fmt"

	"github.com/ThisWillGoWell/stock-simulator-server/src/objects"
	"github.com/ThisWillGoWell/stock-simulator-server/src/objects/effect"

	"github.com/ThisWillGoWell/stock-simulator-server/src/id"

	"github.com/ThisWillGoWell/stock-simulator-server/src/app/log"

	"github.com/ThisWillGoWell/stock-simulator-server/src/wires"

	"github.com/ThisWillGoWell/stock-simulator-server/src/database"
	"github.com/ThisWillGoWell/stock-simulator-server/src/objects/portfolio"
	"github.com/ThisWillGoWell/stock-simulator-server/src/web/session"
)

const minPasswordLength = 4
const minDisplayNameLength = 4
const maxDisplayNameLength = 20

/**
Return a user provided the username and Password
If the Password is correct return user, else return err
*/
func ValidateUser(username, password string) (string, error) {
	UserListLock.Acquire("get-user")
	defer UserListLock.Release()
	userUuid, exists := uuidList[username]
	if !exists {
		return "", errors.New("user does not exist")
	}
	user := UserList[userUuid]

	if !comparePasswords(user.Password, password) {
		return "", errors.New("password is incorrect")
	}
	sessionToken := session.NewSessionToken(user.Uuid)
	return sessionToken, nil
}

/**
Renew a user user a session token
*/
func ConnectUser(sessionToken string) (*User, error) {
	userId, err := session.GetUserId(sessionToken)
	if err != nil {
		return nil, err
	}
	UserListLock.Acquire("renew-user")
	defer UserListLock.Release()
	user, exists := UserList[userId]
	if !exists {
		log.Log.Errorf("a user was found in session token but not in user list? uid=%s", userId)
		return nil, fmt.Errorf("oops! somehting unknown happened 0x58")
	}
	user.Active = true
	wires.UsersUpdate.Offer(user.User)
	return user, nil
}

func GetUserFromToken(sessionToken string) (objects.User, error) {
	u := objects.User{}
	userId, err := session.GetUserId(sessionToken)
	if err != nil {
		return u, err
	}

	user, exists := UserList[userId]
	if !exists {
		log.Log.Errorf("a user was found in session token but not in user list? uid=%s", userId)
		return u, fmt.Errorf("oops! something unknown happened 0x78")
	}
	u = user.User
	return u, nil
}

/**
Build a new user
set their Password to that provided
*/
func NewUser(username, displayName, password string) (string, error) {

	if len(username) > 20 {
		return "", errors.New("username too long")
	}
	if len(username) < 4 {
		return "", errors.New("username too short")
	}
	if !isAllowedCharacterUsername(username) {
		return "", errors.New("username is not allowed")
	}

	if len(password) < minPasswordLength {
		return "", errors.New("password too short")
	}

	if len(displayName) > maxDisplayNameLength {
		return "", errors.New("display name too long")
	}
	if len(displayName) < minDisplayNameLength {
		return "", errors.New("display name too short")
	}
	if !isAllowedCharacterDisplayName(displayName) {
		return "", errors.New("display name contains invalid characters")
	}
	UserListLock.Acquire("make-user")
	defer UserListLock.Release()

	uuid := id.SerialUuid()
	portUuid := id.SerialUuid()

	hashedPassword := hashAndSalt(password)
	u, err := MakeUser(objects.User{Uuid: uuid, PortfolioId: portUuid, UserName: username, DisplayName: displayName, Password: hashedPassword, Config: nil})
	if err != nil {
		id.RemoveUuid(uuid)
		id.RemoveUuid(portUuid)
		log.Log.Errorf("failed to make user err=[%v]", err)
		return "", fmt.Errorf("opps! Something went wrong 0x834")
	}
	port, err := portfolio.NewPortfolio(portUuid, uuid)
	if err != nil {
		id.RemoveUuid(portUuid)
		log.Log.Errorf("failed to make portfolio err=[%v]", err)
		deleteUser(u.Uuid, true)
		return "", fmt.Errorf("opps! Something went wrong 0x042")
	}

	baseEffect, err := effect.NewBaseTradeEffect(portUuid)
	if err != nil {
		portfolio.DeletePortfolio(portUuid)
		deleteUser(u.Uuid, true)
		effect.DeleteEffect(baseEffect)
		log.Log.Errorf("failed to make base trade effect err=[%v]", err)
		return "", fmt.Errorf("opps!, something went wrong 0x72")
	}

	if dbErr := database.Db.Execute([]interface{}{port, u, baseEffect}, nil); dbErr != nil {
		deleteUser(u.Uuid, true)
		portfolio.DeletePortfolio(portUuid)
		log.Log.Errorf("failed to make new user database err=[%v]", err)
		return "", fmt.Errorf("oops! something went wrong 0x48")
	}

	sessionToken := session.NewSessionToken(u.Uuid)

	wires.UsersNewObject.Offer(u.User)
	wires.PortfolioNewObject.Offer(port.Portfolio)
	wires.EffectsNewObject.Offer(baseEffect.Effect)
	return sessionToken, nil
}
