package user

import (
	"github.com/TFTPL/AWS-Cost-Calculator/services/types"
)

type UserInfo struct {
	Name     string
	Email    string
	Password string
}

type UserMethods interface {
	SignupService() error
	Login() (*types.User, error)
	IsUserFoundByEmail() (bool, error)
	IsemailpasswordValidate() error
	IsUserSignupFoundByEmail() (bool, error)
}
