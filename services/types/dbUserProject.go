package types

// DbUserProject golang representation of the pg service_cost_aws table
type DbUserProject struct {
	tableName struct{} `pg:"user_projects,alias:user_projects"`
	ID        int64
	UserId    int64
	ProjectId int64
}
