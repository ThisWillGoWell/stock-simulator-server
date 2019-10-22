package account

import (
	"errors"

	"github.com/ThisWillGoWell/stock-simulator-server/src/wires"

	"github.com/ThisWillGoWell/stock-simulator-server/src/portfolio"
	"github.com/ThisWillGoWell/stock-simulator-server/src/session"
	"github.com/ThisWillGoWell/stock-simulator-server/src/utils"
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
		return nil, errors.New("user found in session list but not in current users")
	}
	user.Active = true
	wires.UsersUpdate.Offer(user)
	return user, nil
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
	uuid := utils.SerialUuid()
	portUuid := utils.SerialUuid()

	hashedPassword := hashAndSalt(password)
	user, err := MakeUser(uuid, username, displayName, hashedPassword, portUuid, "{}")
	if err != nil {
		utils.RemoveUuid(uuid)
		utils.RemoveUuid(portUuid)
		return "", err
	}
	portfolio.NewPortfolio(portUuid, uuid)
	sessionToken := session.NewSessionToken(user.Uuid)
	return sessionToken, nil
}
