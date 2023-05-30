package postgres

import (
	"github.com/TFTPL/AWS-Cost-Calculator/services/types"
	"github.com/pkg/errors"
)

// CreateUserProject creates userProject mapping
func CreateUserProject(userId, projectId int64) (*types.DbUserProject, error) {
	userProject := &types.DbUserProject{
		UserId:    userId,
		ProjectId: projectId,
	}

	err := Insert(userProject)

	if err != nil {
		return nil, errors.Wrap(err, "pgUserProject:CreateUserProject failed - could not insert row")
	}

	return userProject, nil
}

// ReadUserProject reads userProject mapping
func ReadUserProject(id int64) (*types.DbUserProject, error) {
	userProject := &types.DbUserProject{
		ID: id,
	}

	err := getDB().Model(userProject).Where("id = ?", id).Select()

	if err != nil {
		return nil, errors.Wrap(err, "pgUserProject:ReadUserProject failed - could not query row")
	}

	return userProject, nil
}

// UpdateUserProject updates userProject mapping
func UpdateUserProject(id, userId, projectId int64) (*types.DbUserProject, error) {
	userProject := &types.DbUserProject{
		ID:        id,
		UserId:    userId,
		ProjectId: projectId,
	}
	err := Update(userProject)

	if err != nil {
		return nil, errors.Wrap(err, "pgUserProject:UpdateUserProject failed - could not update row")
	}

	return userProject, nil
}

// DeleteUserProject deletes userProject mapping
func DeleteUserProject(id int64) (*types.DbUserProject, error) {
	userProject := &types.DbUserProject{
		ID: id,
	}

	err := Delete(userProject)

	if err != nil {
		return nil, errors.Wrap(err, "pgUserProject:DeleteUserProject failed - could not delete row")
	}

	return userProject, nil
}

func GetProjectsByUserId(userId uint64) ([]types.DbUserProject, error) {
	var userProject []types.DbUserProject

	err := GetDB().Model(&userProject).Where("user_id =?", userId).Select()
	if err != nil {
		return nil, err
	}

	return userProject, nil
}

func GetProjectsIdsByUserId(userid int64) ([]*types.DbUserProject, error) {
	userProject := make([]*types.DbUserProject, 0)

	err := GetDB().Model(&userProject).Where("user_id =?", userid).Select(&userProject)
	if err != nil {
		return nil, err
	}
	return userProject, nil

}
