package types

import (
	"context"
	"time"
)

// DbUser golang representation of the pg users table
type DbUser struct {
	tableName struct{} `pg:"users,alias:users"`
	ID        int64
	Name      string
	Email     string `pg:",unique"`
	Password  string
	CreatedAt time.Time `pg:",notnull,use_zero"`
	UpdatedAt time.Time `pg:",notnull,use_zero"`
}

// BeforeInsert Before insert trigger
func (o *DbUser) BeforeInsert(c context.Context) (context.Context, error) {
	o.CreatedAt = time.Now()
	o.UpdatedAt = time.Now()
	return c, nil
}

// BeforeUpdate Before update trigger
func (o *DbUser) BeforeUpdate(c context.Context) (context.Context, error) {
	o.UpdatedAt = time.Now()
	return c, nil
}
