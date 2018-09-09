package session

import (
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/stock-simulator-server/src/lock"
	"time"
)

const expireTime = time.Hour * 24
const tokenLength = 32
var sessions = make(map[string]*sessionToken)
var sessionsLock = lock.NewLock("sessions")

type sessionToken struct {
	userId       string
	sessionToken string
	lastUse      time.Time
}

func NewSessionToken(userId string) string{
	sessionsLock.Acquire("New token")
	defer sessionsLock.Release()
	newToken, err := generateRandomString(tokenLength)
	if err != nil{
		print("new token err: " + err.Error())
		return ""
	}

	for {
		if _, exists := sessions[newToken]; !exists {
			break
		}
		fmt.Printf("yeah this is bad")
		newToken, err = generateRandomString(tokenLength)
		if err != nil {
			panic(err)
		}
	}
	sessions[newToken] = &sessionToken{
		userId:       userId,
		sessionToken: newToken,
		lastUse:      time.Now(),
	}
	return newToken
}

func GetUserId(sessionToken string) (string, error){
	sessionsLock.Acquire("get sessions")
	defer sessionsLock.Release()
	token, exists := sessions[sessionToken]
	if !exists{
		return "", errors.New("invalid token")
	}
	token.lastUse = time.Now()
	return token.userId, nil

}

func StartSessionCleaner(){
	go runCleanSessionTokens()
}

func runCleanSessionTokens(){
	for {
		sessionsLock.Acquire("clean sessions")
		for key, tokens := range sessions {
			if time.Since(tokens.lastUse) > expireTime{
				delete(sessions, key)
			}
		}
		sessionsLock.Release()
		<- time.After(time.Hour)
	}
}

// GenerateRandomBytes returns securely generated random bytes.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}



// GenerateRandomString returns a securely generated random string.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func generateRandomString(n int) (string, error) {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"
	bytes, err := generateRandomBytes(n)
	if err != nil {
		return "", err
	}
	for i, b := range bytes {
		bytes[i] = letters[b%byte(len(letters))]
	}
	return string(bytes), nil
}

