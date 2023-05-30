package types

// DbProject golang representation of the pg analytics table
type DbAnalytics struct {
	tableName            struct{} `pg:"analytics,alias:analytics"`
	ID                   int64    // ProjectId
	Name                 string
	Email                string
	DailyCostMicrodollar int64
}
