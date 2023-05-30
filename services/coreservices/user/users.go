package user

import (
	"fmt"

	"regexp"
	"strings"

	"github.com/TFTPL/AWS-Cost-Calculator/services/postgres"
	"github.com/TFTPL/AWS-Cost-Calculator/services/types"
	"github.com/pkg/errors"
)

func (u *UserInfo) SignupService() error {

	// Validate the name
	if len(strings.TrimSpace(u.Name)) == 0 || len(strings.TrimSpace(u.Email)) == 0 || len(strings.TrimSpace(u.Password)) == 0 {
		return fmt.Errorf("fill all the required feild")
	}
	//check mail and password
	err := u.IsemailpasswordValidate()
	if err != nil {
		return err
	}

	// Validate the email
	// emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	// if !emailRegex.MatchString(email) {
	// 	return fmt.Errorf("invalid email or password")
	// }

	// // Validate the password
	// if len(strings.TrimSpace(email)) == 0 || strings.TrimSpace(pass) == "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"{
	// 	return fmt.Errorf("email or password is required")
	// }
	// Add code to check if user already exist or
	isUserExist, err := u.IsUserSignupFoundByEmail()
	if err != nil {
		return err
	}

	if isUserExist {
		return errors.New("user already exists")
	}
	_, err = postgres.CreateUser(u.Name, u.Email, u.Password)
	if err != nil {
		return errors.Wrap(err, "coreservices:user:SignupService failed -")
	}

	return nil
}

func (u *UserInfo) Login() (*types.User, error) {

	err := u.IsemailpasswordValidate()
	if err != nil {
		return nil, err
	}

	res, err := u.IsUserFoundByEmail()
	if err != nil {
		return nil, err
	}

	if !res {
		return nil, err
	}

	user, err := postgres.GetUserByEmail(u.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user details: login un-successful")
	}

	if u.Password != user.Password {
		return nil, fmt.Errorf("password mismatched")
	}

	token, err := GenerateJwt(int(user.ID))
	if err != nil {
		fmt.Println("Login : Error in token generation", err.Error())
		return nil, err
	}

	return &types.User{Key: user.ID, Token: token}, nil
}

func (u *UserInfo) IsUserFoundByEmail() (bool, error) {
	res, err := postgres.GetUserByEmail(u.Email)
	if err != nil {
		// some error is database
		return false, fmt.Errorf("some issue from our end, try to login again")
	}
	if res == nil {
		// user does not exist
		return false, fmt.Errorf("user doesn't exists signup first")
	}
	// user exist
	return true, nil
}

func (u *UserInfo) IsemailpasswordValidate() error {
	if len(strings.TrimSpace(u.Email)) == 0 || strings.TrimSpace(u.Password) == "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855" {
		return fmt.Errorf("email or password is required")
	}

	// Validate the email
	emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if !emailRegex.MatchString(u.Email) {
		return fmt.Errorf("invalid email or password")
	}
	return nil
}

func (u *UserInfo) IsUserSignupFoundByEmail() (bool, error) {
	res, err := postgres.GetUserByEmail(u.Email)
	if err != nil {
		// some error is database
		return false, nil
	}
	if res == nil {
		// user does not exist
		return false, nil
	}
	// user exist
	return true, nil
}

func UpdateUserProfile(id int64, name, email, password string) error {

	UserDetail, err := postgres.ReadUser(id)
	if err != nil {
		return err
	}
	if len(name) == 0 {
		name = UserDetail.Name
	}
	if len(email) == 0 {
		email = UserDetail.Email
	}
	if len(password) == 0 {
		password = UserDetail.Password
	}

	_, err = postgres.UpdateUser(id, name, email, password)
	if err != nil {
		return err
	}

	return nil
}
