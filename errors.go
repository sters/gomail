package gomail

import (
	"errors"
	"fmt"
)

var (
	ErrUnencryptedConnection    = errors.New("gomail: unencrypted connection")
	ErrWrongHostName            = errors.New("gomail: wrong host name")
	ErrInvalidMessageFromAbsent = errors.New(`gomail: invalid message, "From" field is absent`)
	ErrCannotWriteAsWriter      = errors.New("gomail: cannot write as writer is in error")
)

// A SendError represents the failure to transmit a Message, detailing the cause
// of the failure and index of the Message within a batch.
type SendError struct {
	// Index specifies the index of the Message within a batch.
	Index uint
	Cause error
}

func (err *SendError) Error() string {
	return fmt.Sprintf("gomail: could not send email %d: %v",
		err.Index+1, err.Cause)
}

func (*SendError) Is(err error) bool {
	if _, ok := err.(*SendError); ok {
		return true
	}
	return false
}

type UnexpectedServerChallengeError struct {
	fromServer []byte
}

func (u *UnexpectedServerChallengeError) Error() string {
	return fmt.Sprintf("gomail: unexpected server challenge: %s", u.fromServer)
}

func (*UnexpectedServerChallengeError) Is(err error) bool {
	if _, ok := err.(*UnexpectedServerChallengeError); ok {
		return true
	}
	return false
}

type InvalidAddress struct {
	field string
	err   error
}

func (i *InvalidAddress) Error() string {
	return fmt.Sprintf("gomail: invalid address %q: %v", i.field, i.err)
}

func (i *InvalidAddress) Unwrap() error {
	return i.err
}

func (*InvalidAddress) Is(err error) bool {
	if _, ok := err.(*InvalidAddress); ok {
		return true
	}
	return false
}

var _ = []error{
	(*SendError)(nil),
	(*UnexpectedServerChallengeError)(nil),
	(*InvalidAddress)(nil),
}
