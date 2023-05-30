package types

import (
	"context"
	"time"
)

// DbProject golang representation of the pg projects table
type DbProject struct {
	tableName    struct{} `pg:"projects,alias:projects"`
	ID           int64
	Name         string
	Description  string
	AwsSecretKey string // depricated
	AccessKey    string // depricated
	IsActive     bool
	CreatedAt    time.Time `pg:",notnull,use_zero"`
	UpdatedAt    time.Time `pg:",notnull,use_zero"`
	IsPinned     bool
}

// BeforeInsert Before insert trigger
func (o *DbProject) BeforeInsert(c context.Context) (context.Context, error) {
	o.CreatedAt = time.Now()
	o.UpdatedAt = time.Now()
	o.IsActive = true
	return c, nil
}

// // BeforeUpdate Before update trigger
func (o *DbProject) BeforeUpdate(c context.Context) (context.Context, error) {
	o.UpdatedAt = time.Now()
	return c, nil
}
