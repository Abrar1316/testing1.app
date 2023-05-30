package types

import "time"

type CostSum struct {
	UsageDate              time.Time `pg:",notnull,use_zero"`
	AccruedCostMicrodollar int64
}

type GetServiceCostMinMAxTotalAvg struct {
	AccruedCostMicrodollar    int64 // `pg:"total_accrued_cost_microdollar"`
	MinAccruedCostMicrodollar int64 `pg:"min_accrued_cost_microdollar"`
	Count                     int64 `pg:"count"`
	MaxAccruedCostMicrodollar int64 `pg:"max_accrued_cost_microdollar"`
}

type DateServiceCost struct {
	UsageDate              time.Time `pg:",notnull,use_zero"`
	ServiceTitle           string
	AccruedCostMicrodollar int64
}

type DateCostByService struct {
	UsageDate              time.Time `pg:",notnull,use_zero"`
	AccruedCostMicrodollar int64
}

type GetProjectByProjectIdResponse struct {
	UserId       string
	Name         string
	Description  string
	StartedOn    string
	ActiveStatus string
	IsPinned     bool
}

type PinnedProjectResponse struct {
	ID           int64
	Name         string
	Description  string
	AwsSecretKey string
	AccessKey    string
	IsActive     bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
	IsPinned     bool
}

type TotalCostSumResonse struct {
	UsageDate              string
	AccruedCostMicrodollar float64
}

type TotalCostProject struct {
	ProjectName            string
	AccruedCostMicrodollar string
}

type TotalCost struct {
	AccruedCostMicrodollar int64
}

type Projectdetail struct {
	Id   int64
	Name string
}
