package postgres

import (
	"log"
	"time"

	"github.com/TFTPL/AWS-Cost-Calculator/services/types"
	"github.com/pkg/errors"
)

// CreateAWSCredentials for create a new asw credentials
func CreateAwsCredentials(projectId int64, awsAccessKey, awsSecretKey string) (*types.DbAwsCredentials, error) {
	credentials := &types.DbAwsCredentials{
		ProjectId:    projectId,
		AwsSecretKey: awsSecretKey,
		AccessKey:    awsAccessKey,
	}
	err := Insert(credentials)

	if err != nil {
		return nil, errors.Wrap(err, "pgProject:CreateProject failed - could not insert row")
	}
	return credentials, nil
}

// ReadAwsCredentials for read aws credentials  by project Id
func ReadAwsCredentials(projectId int64) (*types.DbAwsCredentials, error) {
	credentials := &types.DbAwsCredentials{
		ProjectId: projectId,
	}

	err := getDB().Model(credentials).Where("project_id = ?", projectId).Select()

	if err != nil {
		return nil, errors.Wrap(err, "pgcredentials:ReadAwscredentials failed - could not query row")
	}

	return credentials, nil
}

// UpdateAwsCredentials updates a credentials
func UpdateAwsCredentials(id int64, secretKey, accessKey string) (*types.DbAwsCredentials, error) {

	credentials := &types.DbAwsCredentials{
		ProjectId:    id,
		AwsSecretKey: secretKey,
		AccessKey:    accessKey,
		UpdatedAt:    time.Now(),
	}

	_, err := GetDB().Model(credentials).
		Set("aws_secret_key = ?aws_secret_key").
		Set("access_key = ?access_key").
		Set("updated_at=?updated_at").
		Where("project_id = ?project_id").
		Update()

	if err != nil {
		return nil, errors.Wrap(err, "pgProject:UpdateAwsCredentials failed - could not update row")
	}

	return credentials, nil
}

// Deletecredentials deletes a project
// func DeleteAwsCredentials(id int64) (*types.DbAwsCredentials, error) {
// 	credentials := &types.DbAwsCredentials{
// 		ProjectId: id,
// 	}

// 	err := Delete(credentials)

// 	if err != nil {
// 		return nil, errors.Wrap(err, "pgProject:DeleteAwscredentials failed - could not delete row")
// 	}

// 	return credentials, nil
// }

func DeleteAwsCredentials(id int64) error {

	stmt := `DELETE FROM "aws_credentials" WHERE project_id=?`
	_, err := GetDB().Exec(stmt, id)
	if err != nil {
		log.Println("unable to excute the query")
		return errors.Errorf("Unable to execute query: %v\n", err)
	}

	return nil
}
