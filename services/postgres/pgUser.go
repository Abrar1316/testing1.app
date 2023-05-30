package postgres

import (
	"github.com/TFTPL/AWS-Cost-Calculator/services/types"
	"github.com/go-pg/pg/v10"
	"github.com/pkg/errors"
)

// CreateUser creates a new user
func CreateUser(name, email, password string) (*types.DbUser, error) {
	user := &types.DbUser{
		Name:     name,
		Email:    email,
		Password: password,
	}
	err := Insert(user)

	if err != nil {
		return nil, errors.Wrap(err, "pgUser:CreateUser failed - could not insert row")
	}

	return user, nil
}

// ReadUser reads a user
func ReadUser(id int64) (*types.DbUser, error) {
	user := &types.DbUser{
		ID: id,
	}

	err := getDB().Model(user).Where("id = ?", id).Select()

	if err != nil {
		return nil, errors.Wrap(err, "pgUser:ReadUser failed - could not query row")
	}

	return user, nil
}

// UpdateUser updates a user
func UpdateUser(id int64, name, email, password string) (*types.DbUser, error) {
	user := &types.DbUser{
		ID:       id,
		Name:     name,
		Email:    email,
		Password: password,
	}
	err := Update(user)

	if err != nil {
		return nil, errors.Wrap(err, "pgUser:UpdateUser failed - could not update row")
	}

	return user, nil
}

// DeleteUser deletes a user
func DeleteUser(id int64) (*types.DbUser, error) {
	user := &types.DbUser{
		ID: id,
	}

	err := Delete(user)

	if err != nil {
		return nil, errors.Wrap(err, "pgUser:DeleteUser failed - could not delete row")
	}

	return user, nil
}

// GetUserByEmail reads a user by email
func GetUserByEmail(email string) (*types.DbUser, error) {
	user := &types.DbUser{}

	err := GetDB().Model(user).
		Where("lower(email) = lower(?)", email).
		Select()
	if err != nil {
		if errors.Is(err, pg.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

func GetAllUsers() ([]types.DbUser, error) {

	var userdetails = []types.DbUser{}

	query := GetDB().Model(&types.DbUser{})

	err := query.Select(&userdetails)
	if err != nil {
		if errors.Is(err, pg.ErrNoRows) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return userdetails, nil

}
