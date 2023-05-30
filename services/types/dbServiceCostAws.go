package types

import (
	"time"
)

// DbServiceCostAws golang representation of the pg service_cost_aws table
type DbServiceCostAws struct {
	tableName              struct{} `pg:"service_cost_aws,alias:service_cost_aws"`
	ID                     int64
	ProjectId              int64
	ServiceTitle           string
	UsageDate              time.Time `pg:",notnull,use_zero"`
	AccruedCostMicrodollar int64
}
