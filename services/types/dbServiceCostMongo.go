package types

import (
	"time"
)

// DbServiceCostMongo golang representation of the pg service_cost_aws table
type DbServiceCostMongo struct {
	tableName              struct{} `pg:"service_cost_mongo,alias:service_cost_mongo"`
	ID                     int64
	ProjectId              int64
	MongoProjectName       string
	ServiceTitle           string
	AccruedCostMicrodollar int64
	UsageDate              time.Time `pg:",notnull,use_zero"`
	Unit                   string
	UnitPrice              int64
	Quantity               int64
}
