package types

import (
	"context"
	"time"
)

// DbAwsCredentials golang representation of the pg aws_credentials table
type DbAwsCredentials struct {
	tableName    struct{} `pg:"aws_credentials,alias:aws_credentials"`
	ID           int64
	ProjectId    int64
	AwsSecretKey string
	AccessKey    string
	CreatedAt    time.Time `pg:",notnull,use_zero"`
	UpdatedAt    time.Time `pg:",notnull,use_zero"`
}

// BeforeInsert Before insert trigger
func (o *DbAwsCredentials) BeforeInsert(c context.Context) (context.Context, error) {
	o.CreatedAt = time.Now()
	o.UpdatedAt = time.Now()
	return c, nil
}

// BeforeUpdate Before update trigger
func (o *DbAwsCredentials) BeforeUpdate(c context.Context) (context.Context, error) {
	o.UpdatedAt = time.Now()
	return c, nil
}
