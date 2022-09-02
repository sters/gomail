package gomail

import (
	"bytes"
	"net/smtp"
)

const (
	loginMechanism = "LOGIN"
)

// loginAuth is an smtp.Auth that implements the LOGIN authentication mechanism.
type loginAuth struct {
	username string
	password string
	host     string
}

func (a *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	if !server.TLS {
		advertised := false
		for _, mechanism := range server.Auth {
			if mechanism == loginMechanism {
				advertised = true
				break
			}
		}
		if !advertised {
			return "", nil, ErrUnencryptedConnection
		}
	}
	if server.Name != a.host {
		return "", nil, ErrWrongHostName
	}
	return loginMechanism, nil, nil
}

func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if !more {
		return nil, nil
	}

	switch {
	case bytes.Equal(fromServer, []byte("Username:")):
		return []byte(a.username), nil
	case bytes.Equal(fromServer, []byte("Password:")):
		return []byte(a.password), nil
	default:
		return nil, &UnexpectedServerChallengeError{fromServer}
	}
}
