package types

import (
	"context"
	"time"
)

// DbMongoCredentials golang representation of the pg aws_credentials table
type DbMongoCredentials struct {
	tableName      struct{} `pg:"mongo_credentials,alias:mongo_credentials"`
	ID             int64
	ProjectId      int64
	MongoOrganizationId string
	PublicKey      string
	PrivateKey     string
	CreatedAt      time.Time `pg:",notnull,use_zero"`
	UpdatedAt      time.Time `pg:",notnull,use_zero"`
}

// BeforeInsert Before insert trigger
func (o *DbMongoCredentials) BeforeInsert(c context.Context) (context.Context, error) {
	o.CreatedAt = time.Now()
	o.UpdatedAt = time.Now()
	return c, nil
}

// BeforeUpdate Before update trigger
func (o *DbMongoCredentials) BeforeUpdate(c context.Context) (context.Context, error) {
	o.UpdatedAt = time.Now()
	return c, nil
}
