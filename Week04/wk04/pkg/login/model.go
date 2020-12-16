package login

import (
	"fmt"
)

const (
	NoQueryStaus    = 1
	ErrSystemStatus = 99
	Success         = 200
)

type ErrNoQuery struct{}
type ErrSystemErr struct{}

type loginUser struct {
	Id       int64
	Name     string
	Password string
}

func (e ErrNoQuery) Error() string {
	return fmt.Sprintf("%w%s", NoQueryStaus, "no query")
}

func (e ErrSystemErr) Error() string {
	return fmt.Sprintf("%w%s", ErrSystemStatus, "system panic")
}
