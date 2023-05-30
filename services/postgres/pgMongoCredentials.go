package postgres

import (
	"log"
	"time"

	"github.com/TFTPL/AWS-Cost-Calculator/services/types"
	"github.com/go-pg/pg/v10"
	"github.com/pkg/errors"
)

// CreateMongoCredentials for create a new credentials
func CreateMongoCredentials(projectId int64, OrganizationId, PublicId, PrivateId string) (*types.DbMongoCredentials, error) {
	credentials := &types.DbMongoCredentials{
		ProjectId:           projectId,
		MongoOrganizationId: OrganizationId,
		PublicKey:           PublicId,
		PrivateKey:          PrivateId,
	}
	err := Insert(credentials)

	if err != nil {
		return nil, errors.Wrap(err, "pgProject:CreateMongoCred failed - could not insert row")
	}
	return credentials, nil
}

// ReadMongoCredentials for read credentials  by project Id
func ReadMongoCredentials(projectId int64) (*types.DbMongoCredentials, error) {
	credentials := &types.DbMongoCredentials{
		ProjectId: projectId,
	}

	err := getDB().Model(credentials).Where("project_id = ?", projectId).Select()

	if err != nil {
		if errors.Is(err, pg.ErrNoRows) {
			return nil, nil
		} else {
			return nil, errors.Wrap(err, "pgcredentials:ReadMongocredentials failed - could not query row")
		}
	}

	return credentials, nil
}

// UpdateMongoCredentials updates a credentials
func UpdateMongoCredentials(id int64, OrganizationId, PublicId, PrivateId string) (*types.DbMongoCredentials, error) {

	credentials := &types.DbMongoCredentials{
		ProjectId:           id,
		MongoOrganizationId: OrganizationId,
		PublicKey:           PublicId,
		PrivateKey:          PrivateId,
		UpdatedAt:           time.Now(),
	}

	_, err := GetDB().Model(credentials).
		Set("mongo_organization_id =?mongo_organization_id").
		Set("public_key = ?public_key").
		Set("private_key = ?private_key").
		Set("updated_at=?updated_at").
		Where("project_id = ?project_id").
		Update()

	if err != nil {
		return nil, errors.Wrap(err, "pgProject:UpdateMongoCredentials failed - could not update row")
	}

	return credentials, nil
}

// DeleteMongoCredentials updates a credentials
func DeleteMongoCredentials(id int64) error {

	stmt := `DELETE FROM "mongo_credentials" WHERE project_id=?`
	_, err := GetDB().Exec(stmt, id)
	if err != nil {
		log.Println("unable to excute the query")
		return errors.Errorf("Unable to execute query: %v\n", err)
	}

	return nil
}

// ReadMongoCredentials for read Mongo credentials
func ReadAllMongoCredentials() ([]*types.DbMongoCredentials, error) {
	credentials := make([]*types.DbMongoCredentials, 0)

	err := getDB().Model(&types.DbMongoCredentials{}).Select(&credentials)

	// if err != nil {
	// 	return nil, errors.Wrap(err, "pgcredentials:ReadMongocredentials failed - could not query row")
	// }
	if err != nil {
		if errors.Is(err, pg.ErrNoRows) {
			return nil, nil
		} else {
			return nil, errors.Wrap(err, "pgcredentials:ReadMongocredentials failed - could not query row")
		}
	}

	return credentials, nil
}
